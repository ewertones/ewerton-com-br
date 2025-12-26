"use client";

import { createContext, useContext, useEffect, useState, ReactNode } from "react";

type Theme = "light" | "dark";

type ThemeContextType = {
    theme: Theme;
    toggleTheme: () => void;
    setTheme: (theme: Theme) => void;
};

const ThemeContext = createContext<ThemeContextType | undefined>(undefined);

export function ThemeProvider({ children }: { children: ReactNode }) {
    const [theme, setThemeState] = useState<Theme>("light");
    const [mounted, setMounted] = useState(false);

    useEffect(() => {
        setMounted(true);
        const saved = localStorage.getItem("theme") as Theme | null;
        if (saved) {
            setThemeState(saved);
            document.documentElement.setAttribute("data-theme", saved);
        } else {
            setThemeState("light");
            document.documentElement.setAttribute("data-theme", "light");
        }
    }, []);

    const setTheme = (newTheme: Theme) => {
        setThemeState(newTheme);
        localStorage.setItem("theme", newTheme);
        document.documentElement.setAttribute("data-theme", newTheme);
    };

    const toggleTheme = () => {
        setTheme(theme === "light" ? "dark" : "light");
    };

    // Always provide the context, but use safe defaults before mounting
    const value: ThemeContextType = {
        theme: mounted ? theme : "light",
        toggleTheme: mounted ? toggleTheme : () => {},
        setTheme: mounted ? setTheme : () => {},
    };

    return <ThemeContext.Provider value={value}>{children}</ThemeContext.Provider>;
}

export function useTheme() {
    const context = useContext(ThemeContext);
    if (!context) {
        throw new Error("useTheme must be used within a ThemeProvider");
    }
    return context;
}
