"use client";

import { Button as ButtonPrimitive } from "@base-ui/react/button";
import { cva, type VariantProps } from "class-variance-authority";

import { cn } from "@/lib/utils";

const buttonVariants = cva(
  "group/button inline-flex shrink-0 items-center justify-center rounded-md border border-transparent bg-clip-padding font-sans font-bold text-sm uppercase tracking-widest whitespace-nowrap transition-all outline-none select-none focus-visible:border-ring focus-visible:ring-2 focus-visible:ring-ring/30 active:not-aria-[haspopup]:translate-y-px aria-disabled:pointer-events-none aria-disabled:opacity-30 aria-invalid:border-destructive aria-invalid:ring-2 aria-invalid:ring-destructive/20 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4",
  {
    variants: {
      variant: {
        default:
          "relative bg-transparent border border-primary text-primary transition-all duration-350 ease-in overflow-hidden cursor-pointer rounded-none h-auto before:content-[''] before:absolute before:inset-0 before:bg-primary before:scale-x-0 before:origin-left before:transition-transform before:duration-350 before:ease-in before:z-0 hover:before:scale-x-100 hover:text-background hover:shadow-[0_0_16px_oklch(0.75_0.12_77/0.3)]",
        outline:
          "relative bg-transparent border border-border text-muted-foreground transition-all duration-300 ease-in cursor-pointer rounded-none h-auto hover:border-foreground hover:text-foreground",
        secondary:
          "relative bg-transparent border border-secondary text-secondary transition-all duration-350 ease-in overflow-hidden cursor-pointer rounded-none h-auto before:content-[''] before:absolute before:inset-0 before:bg-secondary before:scale-x-0 before:origin-left before:transition-transform before:duration-350 before:ease-in before:z-0 hover:before:scale-x-100 hover:text-secondary-foreground hover:shadow-[0_0_16px_oklch(0.72_0.09_190/0.3)]",
        ghost:
          "hover:bg-muted hover:text-foreground aria-expanded:bg-muted aria-expanded:text-foreground",
        destructive:
          "bg-destructive/10 text-destructive hover:bg-destructive/20 focus-visible:border-destructive/40 focus-visible:ring-destructive/20",
        link: "text-primary underline-offset-4 hover:underline",
      },
      size: {
        default:
          "gap-1 px-8 py-3.5 has-data-[icon=inline-end]:pr-6 has-data-[icon=inline-start]:pl-6 [&_svg:not([class*='size-'])]:size-3.5",
        xs: "h-5 gap-1 rounded-sm px-2 text-[0.625rem] has-data-[icon=inline-end]:pr-1.5 has-data-[icon=inline-start]:pl-1.5 [&_svg:not([class*='size-'])]:size-2.5",
        sm: "gap-1 px-5 py-2 text-xs has-data-[icon=inline-end]:pr-4 has-data-[icon=inline-start]:pl-4 [&_svg:not([class*='size-'])]:size-3",
        lg: "gap-1 px-10 py-4 has-data-[icon=inline-end]:pr-8 has-data-[icon=inline-start]:pl-8 [&_svg:not([class*='size-'])]:size-4 font-black",
        icon: "size-7 [&_svg:not([class*='size-'])]:size-3.5",
        "icon-xs": "size-5 rounded-sm [&_svg:not([class*='size-'])]:size-2.5",
        "icon-sm": "size-6 [&_svg:not([class*='size-'])]:size-3",
        "icon-lg": "size-8 [&_svg:not([class*='size-'])]:size-4",
      },
    },
    defaultVariants: {
      variant: "default",
      size: "default",
    },
  },
);

function Button({
  className,
  variant = "default",
  size = "default",
  children,
  disabled,
  ...props
}: ButtonPrimitive.Props & VariantProps<typeof buttonVariants>) {
  const wrapChildren = variant === "default" || variant === "secondary";
  return (
    <ButtonPrimitive
      data-slot="button"
      className={cn(buttonVariants({ variant, size, className }))}
      aria-disabled={disabled}
      {...props}
    >
      {wrapChildren ? (
        <span className="relative z-1">{children}</span>
      ) : (
        children
      )}
    </ButtonPrimitive>
  );
}

export { Button, buttonVariants };
