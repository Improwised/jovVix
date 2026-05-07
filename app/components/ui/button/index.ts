import type { VariantProps } from "class-variance-authority";
import { cva } from "class-variance-authority";

export { default as Button } from "./Button.vue";

export const buttonVariants = cva(
  "inline-flex items-center justify-center font-headings text-center whitespace-nowrap jv-border-uneven shadow-brutal transition-transform transform transform-gpu backface-hidden active:shadow-none",
  {
    variants: {
      variant: {
        default: "bg-jv-yellow-2 text-jv-ink",
        outline:
          "border-border bg-background hover:bg-muted hover:text-foreground dark:bg-input/30 dark:border-input dark:hover:bg-input/50 aria-expanded:bg-muted aria-expanded:text-foreground",
        secondary: "bg-jv-canvas text-jv-ink",
        accent: "bg-jv-coral text-white",
        ghost:
          "hover:bg-muted hover:text-foreground dark:hover:bg-muted/50 aria-expanded:bg-muted aria-expanded:text-foreground",
        destructive:
          "bg-destructive/10 hover:bg-destructive/20 focus-visible:ring-destructive/20 dark:focus-visible:ring-destructive/40 dark:bg-destructive/20 text-destructive focus-visible:border-destructive/40 dark:hover:bg-destructive/30",
        link: "text-primary underline-offset-4 hover:underline",
      },
      size: {
        default:
          "gap-1.5 px-2.5 has-data-[icon=inline-end]:pr-2 has-data-[icon=inline-start]:pl-2",
        xs: "gap-1 px-5 py-2 text-xs in-data-[slot=button-group]:rounded-lg has-data-[icon=inline-end]:pr-1.5 has-data-[icon=inline-start]:pl-1.5 [&_svg:not([class*='size-'])]:size-3",
        sm: "gap-1 px-6 py-2.5 text-xs md:text-base",
        lg: "gap-1.5 px-7 py-4 has-data-[icon=inline-end]:pr-2 has-data-[icon=inline-start]:pl-2 text-sm md:text-lg",
        icon: "size-8",
        "icon-xs":
          "size-6 in-data-[slot=button-group]:rounded-lg [&_svg:not([class*='size-'])]:size-3",
        "icon-sm": "size-7 in-data-[slot=button-group]:rounded-lg",
        "icon-lg": "size-9",
      },
      tilt: {
        "-2": "hover:rotate-[-2deg]",
        "-1": "hover:rotate-[-1deg]",
        0: "hover:rotate-0",
        1: "hover:rotate-[1deg]",
        2: "hover:rotate-[2deg]",
      },
    },
    defaultVariants: {
      variant: "default",
      size: "default",
      tilt: 0,
    },
  }
);
export type ButtonVariants = VariantProps<typeof buttonVariants>;
