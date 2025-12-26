import { useTranslations } from "next-intl";
import { Link } from "@/i18n/routing";
import Image from "next/image";
import styles from "./Footer.module.css";

export function Footer() {
    const t = useTranslations("footer");

    return (
        <footer className={styles.footer}>
            <div className={styles.content}>
                <Link href="/" className={styles.logoLink}>
                    <Image src="/logo.svg" alt="ewerton" width={48} height={48} className={styles.logo} />
                </Link>

                <div className={styles.bottom}>
                    <p className={styles.copyright}>{t("copyright")}</p>
                    <p className={styles.madeWith}>
                        {t("madeWith")} <span className={styles.brazil}>{t("brazil")}</span>
                    </p>
                </div>
            </div>
        </footer>
    );
}
