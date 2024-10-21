import { useToast } from "vue-toastification";
const toast = useToast();

export default function usecopyToClipboard(text) {
  navigator.clipboard
    .writeText(text)
    .then(() => {
      toast.success("Code copied to clipboard");
    })
    .catch((error) => {
      toast.warning("Error copying to clipboard", error);
    });
}
