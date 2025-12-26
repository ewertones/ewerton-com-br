"use client";

import { Sun, Moon } from "lucide-react";
import { useTheme } from "@/lib/ThemeProvider";
import styles from "./ThemeToggle.module.css";
import clsx from "clsx";

export function ThemeToggle() {
    const { theme, toggleTheme } = useTheme();

    return (
        <button
            onClick={toggleTheme}
            className={clsx(styles.toggle, theme === "dark" && styles.toggleDark)}
            aria-label={`Switch to ${theme === "light" ? "dark" : "light"} mode`}
            title={theme === "light" ? "Switch to dark mode" : "Switch to light mode"}
        >
            {theme === "light" ? <Moon size={18} className={styles.moon} /> : <Sun size={18} className={styles.sun} />}
        </button>
    );
}
