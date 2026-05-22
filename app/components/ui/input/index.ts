import type { VariantProps } from "class-variance-authority";
import { cva } from "class-variance-authority";

export { default as Input } from "./Input.vue";
export { default as Textarea } from "./Textarea.vue";
export { default as Label } from "./Label.vue";

export const inputVariants = cva(
  "w-full font-body text-jv-ink placeholder:text-jv-muted bg-jv-white border-[2px] border-jv-ink outline-none transition-shadow focus:shadow-brutal-sm disabled:cursor-not-allowed disabled:opacity-60",
  {
    variants: {
      size: {
        sm: "h-9 px-3 text-sm rounded-[6px]",
        md: "h-11 px-4 text-base rounded-[8px]",
        lg: "h-12 px-5 text-base rounded-[10px]",
      },
      invalid: {
        true: "border-jv-coral",
        false: "",
      },
    },
    defaultVariants: {
      size: "md",
      invalid: false,
    },
  }
);

export type InputVariants = VariantProps<typeof inputVariants>;
