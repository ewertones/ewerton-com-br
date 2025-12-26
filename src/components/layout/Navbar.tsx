"use client";

import { Link } from "@/i18n/routing";
import { ThemeToggle } from "@/components/ui/ThemeToggle";
import { LanguageSwitcher } from "@/components/ui/LanguageSwitcher";
import styles from "./Navbar.module.css";

export function Navbar() {
    return (
        <nav className={styles.navbar}>
            <div className={styles.navContent}>
                <Link href="/" className={styles.logo}>
                    ewerton
                </Link>
                <div className={styles.navActions}>
                    <LanguageSwitcher />
                    <ThemeToggle />
                </div>
            </div>
        </nav>
    );
}
