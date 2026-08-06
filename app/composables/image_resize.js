import { usePush } from "notivue";

const MAX_EDGE = 800;
const QUALITY_LADDER = [0.8, 0.7, 0.6];
const TARGET_BYTES = 150 * 1024;
const UNRESIZABLE_TYPES = ["image/gif"];

export const readAsBase64 = (file) =>
  new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.onload = (e) => resolve(e.target.result);
    reader.onerror = () => reject(reader.error);
    reader.readAsDataURL(file);
  });

let webpSupport = null;

const supportsWebp = () => {
  if (webpSupport !== null) return webpSupport;
  const probe = document.createElement("canvas");
  probe.width = 1;
  probe.height = 1;
  webpSupport = probe.toDataURL("image/webp").startsWith("data:image/webp");
  return webpSupport;
};

const decodeImage = async (file) => {
  if (typeof createImageBitmap === "function") {
    try {
      return await createImageBitmap(file, { imageOrientation: "from-image" });
    } catch {
      // falls through to the <img> decoder
    }
  }

  return new Promise((resolve) => {
    const objectUrl = URL.createObjectURL(file);
    const image = new Image();
    const settle = (value) => {
      URL.revokeObjectURL(objectUrl);
      resolve(value);
    };
    image.onload = () => settle(image);
    image.onerror = () => settle(null);
    image.src = objectUrl;
  });
};

const toBlob = (canvas, type, quality) =>
  new Promise((resolve) => canvas.toBlob(resolve, type, quality));

export const resizeImageToDataUrl = async (file) => {
  if (typeof document === "undefined") return null;
  if (UNRESIZABLE_TYPES.includes(file.type)) return null;

  const source = await decodeImage(file);
  if (!source) return null;

  try {
    const width = source.naturalWidth || source.width;
    const height = source.naturalHeight || source.height;
    if (!width || !height) return null;

    const scale = Math.min(1, MAX_EDGE / Math.max(width, height));
    const targetWidth = Math.max(1, Math.round(width * scale));
    const targetHeight = Math.max(1, Math.round(height * scale));

    const canvas = document.createElement("canvas");
    canvas.width = targetWidth;
    canvas.height = targetHeight;

    const context = canvas.getContext("2d");
    if (!context) return null;
    context.imageSmoothingEnabled = true;
    context.imageSmoothingQuality = "high";

    const type = supportsWebp() ? "image/webp" : "image/jpeg";
    if (type === "image/jpeg") {
      context.fillStyle = "#ffffff";
      context.fillRect(0, 0, targetWidth, targetHeight);
    }
    context.drawImage(source, 0, 0, targetWidth, targetHeight);

    let blob = null;
    for (const quality of QUALITY_LADDER) {
      blob = await toBlob(canvas, type, quality);
      if (!blob) return null;
      if (blob.size <= TARGET_BYTES) break;
    }
    if (!blob) return null;

    if (blob.size >= file.size) {
      return { dataUrl: await readAsBase64(file), resized: false };
    }

    return { dataUrl: await readAsBase64(blob), resized: true };
  } finally {
    if (typeof source.close === "function") source.close();
  }
};

export const useImagePicker = () => {
  const toast = usePush();
  const app = useNuxtApp();
  const config = useRuntimeConfig().public;

  const pickImageFile = async (file) => {
    if (!file) return null;

    if (!app.$validImageTypes.includes(file.type)) {
      toast.error(
        "Please upload a valid image file (JPEG, PNG, GIF, WEBP, HEIC, HEIF)."
      );
      return null;
    }

    if (file.size > config.maxImageUploadSize) {
      const limitMb = Math.round(config.maxImageUploadSize / (1024 * 1024));
      toast.error(`Please upload an image less than ${limitMb} MB.`);
      return null;
    }

    const shrunk = await resizeImageToDataUrl(file).catch(() => null);
    if (shrunk) return { dataUrl: shrunk.dataUrl, name: file.name };

    if (file.size > config.maxImageFileSize) {
      const limitKb = Math.round(config.maxImageFileSize / 1024);
      toast.error(`Please upload an image less than ${limitKb} KB.`);
      return null;
    }

    return { dataUrl: await readAsBase64(file), name: file.name };
  };

  return { pickImageFile };
};
