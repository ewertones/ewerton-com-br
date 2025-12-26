import React from "react";
import { clsx } from "clsx";
import { Loader2 } from "lucide-react";
import styles from "./Button.module.css";

interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
    variant?: "primary" | "secondary" | "ghost" | "danger";
    size?: "sm" | "md" | "lg";
    isLoading?: boolean;
    leftIcon?: React.ReactNode;
    rightIcon?: React.ReactNode;
    fullWidth?: boolean;
}

export function Button({
    className,
    variant = "primary",
    size = "md",
    isLoading,
    leftIcon,
    rightIcon,
    fullWidth,
    children,
    disabled,
    ...props
}: ButtonProps) {
    return (
        <button
            className={clsx(styles.button, styles[variant], styles[size], fullWidth && styles.fullWidth, className)}
            disabled={disabled || isLoading}
            {...props}
        >
            {isLoading && <Loader2 className={styles.spinner} size={16} />}
            {!isLoading && leftIcon && <span className={styles.iconLeft}>{leftIcon}</span>}
            {children}
            {!isLoading && rightIcon && <span className={styles.iconRight}>{rightIcon}</span>}
        </button>
    );
}
