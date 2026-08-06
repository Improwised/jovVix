import { readAsBase64, useImagePicker } from "./image_resize";

// The API serializes sql.NullString as { String, Valid } when populated; null otherwise.
export const nullableString = (value) => {
  if (!value) return "";
  if (typeof value === "string") return value;
  return value.String || "";
};

// Cover images are stored as base64 data URIs in the DB rather than in object
// storage, so they are read client-side and sent inline as a JSON string.
export const useCoverImage = () => {
  const { pickImageFile } = useImagePicker();

  const pickCoverImage = async (event) => {
    const file = event.target.files?.[0];
    if (!file) return null;

    const picked = await pickImageFile(file);
    if (!picked) {
      event.target.value = "";
      return null;
    }

    return picked;
  };

  return { readAsBase64, pickCoverImage };
};
