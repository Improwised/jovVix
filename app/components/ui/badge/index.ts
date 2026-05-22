import type { VariantProps } from "class-variance-authority";
import { cva } from "class-variance-authority";

export { default as Badge } from "./Badge.vue";

export const badgeVariants = cva(
  "inline-flex items-center gap-1 font-body font-bold border-[2px] border-jv-ink",
  {
    variants: {
      variant: {
        default: "bg-jv-yellow text-jv-ink",
        coral: "bg-jv-coral text-jv-white",
        mint: "bg-jv-mint text-jv-ink",
        salmon: "bg-jv-salmon text-jv-ink",
        lavender: "bg-jv-lavender text-jv-ink",
        white: "bg-jv-white text-jv-ink",
        ink: "bg-jv-ink text-jv-white",
        success: "bg-jv-mint text-jv-ink",
        warning: "bg-jv-yellow-soft text-jv-ink",
        danger: "bg-jv-coral text-jv-white",
        info: "bg-jv-lavender text-jv-ink",
      },
      size: {
        xs: "h-5 px-2 text-[10px] rounded-[6px]",
        sm: "h-6 px-2.5 text-xs rounded-[8px]",
        md: "h-7 px-3 text-sm rounded-[8px]",
        lg: "h-8 px-3.5 text-sm rounded-[10px]",
      },
      shadow: {
        none: "",
        sm: "shadow-brutal-sm",
        md: "shadow-brutal",
      },
    },
    defaultVariants: {
      variant: "default",
      size: "sm",
      shadow: "none",
    },
  }
);

export type BadgeVariants = VariantProps<typeof badgeVariants>;
