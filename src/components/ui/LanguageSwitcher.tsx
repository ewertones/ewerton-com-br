"use client";

import { useLocale } from "next-intl";
import { useRouter, usePathname } from "@/i18n/routing";
import { routing } from "@/i18n/routing";
import { useState, useRef, useEffect } from "react";
import { ChevronDown } from "lucide-react";
import styles from "./LanguageSwitcher.module.css";
import clsx from "clsx";

const localeData: Record<string, { flag: string; name: string }> = {
    en: { flag: "🇺🇸", name: "English" },
    "pt-br": { flag: "🇧🇷", name: "Português" },
};

export function LanguageSwitcher() {
    const locale = useLocale();
    const router = useRouter();
    const pathname = usePathname();
    const [isOpen, setIsOpen] = useState(false);
    const ref = useRef<HTMLDivElement>(null);

    useEffect(() => {
        const handleClickOutside = (event: MouseEvent) => {
            if (ref.current && !ref.current.contains(event.target as Node)) {
                setIsOpen(false);
            }
        };

        document.addEventListener("mousedown", handleClickOutside);
        return () => document.removeEventListener("mousedown", handleClickOutside);
    }, []);

    const handleSwitch = (newLocale: string) => {
        router.replace(pathname, { locale: newLocale as "en" | "pt-br" });
        setIsOpen(false);
    };

    return (
        <div className={styles.switcher} ref={ref}>
            <button className={styles.button} onClick={() => setIsOpen(!isOpen)} aria-label="Switch language">
                <span className={styles.flag}>{localeData[locale].flag}</span>
                <ChevronDown size={14} className={clsx(styles.chevron, isOpen && styles.chevronOpen)} />
            </button>

            {isOpen && (
                <div className={styles.dropdown}>
                    {routing.locales.map((loc) => (
                        <button
                            key={loc}
                            className={clsx(styles.option, loc === locale && styles.optionActive)}
                            onClick={() => handleSwitch(loc)}
                        >
                            <span className={styles.flag}>{localeData[loc].flag}</span>
                            {localeData[loc].name}
                        </button>
                    ))}
                </div>
            )}
        </div>
    );
}
