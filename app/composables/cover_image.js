import { usePush } from "notivue";

// The API serializes sql.NullString as { String, Valid } when populated; null otherwise.
export const nullableString = (value) => {
  if (!value) return "";
  if (typeof value === "string") return value;
  return value.String || "";
};

export const readAsBase64 = (file) =>
  new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.onload = (e) => resolve(e.target.result);
    reader.onerror = () => reject(reader.error);
    reader.readAsDataURL(file);
  });

// Cover images are stored as base64 data URIs in the DB rather than in object
// storage, so they are read client-side and sent inline as a JSON string.
export const useCoverImage = () => {
  const toast = usePush();
  const app = useNuxtApp();
  const url = useRuntimeConfig().public;

  // Returns { dataUrl, name }, or null when the file was rejected.
  const pickCoverImage = async (event) => {
    const file = event.target.files?.[0];
    if (!file) return null;

    if (!app.$validImageTypes.includes(file.type)) {
      toast.error(
        "Please upload a valid image file (JPEG, PNG, GIF, WEBP, HEIC, HEIF)."
      );
      event.target.value = "";
      return null;
    }

    if (file.size > url.maxImageFileSize) {
      const limitKb = Math.round(url.maxImageFileSize / 1024);
      toast.error(`Please upload an image less than ${limitKb} KB.`);
      event.target.value = "";
      return null;
    }

    return { dataUrl: await readAsBase64(file), name: file.name };
  };

  return { readAsBase64, pickCoverImage };
};
