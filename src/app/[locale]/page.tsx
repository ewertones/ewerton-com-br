"use client";

import { useTranslations } from "next-intl";
import { useState } from "react";
import { Button } from "@/components/ui/Button";
import { Navbar } from "@/components/layout/Navbar";
import { Footer } from "@/components/layout/Footer";
import {
    Linkedin,
    FileText,
    MapPin,
    Code2,
    Zap,
    Shield,
    LayoutGrid,
    ArrowRight,
    Brain,
    Gamepad2,
    Users,
} from "lucide-react";
import styles from "./page.module.css";
import Image from "next/image";

const LINKEDIN_URL = "https://www.linkedin.com/in/ewertones";
const AIOMOVER_URL = "https://aiomover.com";

export default function HomePage() {
    const t = useTranslations();
    const [showPdfViewer, setShowPdfViewer] = useState(false);

    return (
        <div className={styles.page}>
            <Navbar />

            <main className={styles.main}>
                {/* Hero Section */}
                <section className={styles.hero}>
                    <h1 className={styles.heroTitle}>
                        {t("landing.hero.title1")}
                        <br />
                        {t("landing.hero.title2")}
                    </h1>
                    <p className={styles.heroSubtitle}>{t("landing.hero.subtitle")}</p>

                    <div className={styles.heroCta}>
                        <a href={LINKEDIN_URL} target="_blank" rel="noopener noreferrer">
                            <Button variant="primary" size="lg" leftIcon={<Linkedin size={20} />}>
                                {t("common.connectLinkedIn")}
                            </Button>
                        </a>
                        <Button
                            variant="secondary"
                            size="lg"
                            leftIcon={<FileText size={20} />}
                            onClick={() => setShowPdfViewer(true)}
                        >
                            {t("common.viewResume")}
                        </Button>
                    </div>
                </section>

                {/* Projects Section */}
                <section className={styles.section}>
                    <h2 className={styles.sectionTitle}>{t("landing.projects.title")}</h2>
                    <p className={styles.sectionSubtitle}>{t("landing.projects.subtitle")}</p>

                    <a href={AIOMOVER_URL} target="_blank" rel="noopener noreferrer" className={styles.projectCardLink}>
                        <div className={styles.projectCard}>
                            <div className={styles.projectHeader}>
                                <div className={styles.projectLogo}>
                                    <Image src="/projects/aiomover.svg" alt="aiomover" width={48} height={48} />
                                </div>
                                <div className={styles.projectTitleGroup}>
                                    <h3 className={styles.projectTitle}>{t("landing.projects.aiomover.title")}</h3>
                                    <p className={styles.projectTagline}>{t("landing.projects.aiomover.tagline")}</p>
                                </div>
                                <ArrowRight size={24} className={styles.projectArrow} />
                            </div>

                            <p className={styles.projectDescription}>{t("landing.projects.aiomover.description")}</p>

                            <div className={styles.projectFeatures}>
                                <div className={styles.featureItem}>
                                    <span className={styles.featureIcon}>
                                        <Zap size={20} />
                                    </span>
                                    {t("landing.projects.aiomover.features.fast")}
                                </div>
                                <div className={styles.featureItem}>
                                    <span className={styles.featureIcon}>
                                        <Shield size={20} />
                                    </span>
                                    {t("landing.projects.aiomover.features.secure")}
                                </div>
                                <div className={styles.featureItem}>
                                    <span className={styles.featureIcon}>
                                        <LayoutGrid size={20} />
                                    </span>
                                    {t("landing.projects.aiomover.features.simple")}
                                </div>
                            </div>
                        </div>
                    </a>

                    {/* questao - Disabled */}
                    <div className={`${styles.projectCard} ${styles.projectCardDisabled}`}>
                        <div className={styles.projectHeader}>
                            <div className={styles.projectLogo}>
                                <Image src="/projects/questao.svg" alt="questao" width={48} height={48} />
                            </div>
                            <div className={styles.projectTitleGroup}>
                                <h3 className={styles.projectTitle}>{t("landing.projects.questao.title")}</h3>
                                <p className={styles.projectTagline}>{t("landing.projects.questao.tagline")}</p>
                            </div>
                            <span className={styles.comingSoonBadge}>{t("common.comingSoon")}</span>
                        </div>

                        <p className={styles.projectDescription}>{t("landing.projects.questao.description")}</p>

                        <div className={styles.projectFeatures}>
                            <div className={styles.featureItem}>
                                <span className={styles.featureIcon}>
                                    <Brain size={20} />
                                </span>
                                {t("landing.projects.questao.features.ai")}
                            </div>
                            <div className={styles.featureItem}>
                                <span className={styles.featureIcon}>
                                    <Gamepad2 size={20} />
                                </span>
                                {t("landing.projects.questao.features.gamified")}
                            </div>
                            <div className={styles.featureItem}>
                                <span className={styles.featureIcon}>
                                    <Users size={20} />
                                </span>
                                {t("landing.projects.questao.features.community")}
                            </div>
                        </div>
                    </div>
                </section>

                {/* About Section */}
                <section className={styles.section}>
                    <h2 className={styles.sectionTitle}>{t("landing.about.title")}</h2>

                    <div className={styles.aboutContent}>
                        <Image
                            src="/photos/ewerton.webp"
                            alt="Ewerton"
                            width={300}
                            height={300}
                            className={styles.photo}
                        />

                        <div className={styles.aboutText}>
                            <p className={styles.aboutDescription}>{t("landing.about.description")}</p>

                            <div className={styles.aboutMeta}>
                                <div className={styles.aboutMetaItem}>
                                    <MapPin size={18} />
                                    {t("landing.about.location")}
                                </div>
                                <div className={styles.aboutMetaItem}>
                                    <Code2 size={18} />
                                    {t("landing.about.focus")}
                                </div>
                            </div>
                        </div>
                    </div>
                </section>
            </main>

            <Footer />

            {/* PDF Viewer Modal */}
            {showPdfViewer && (
                <div className={styles.pdfViewerOverlay} onClick={() => setShowPdfViewer(false)}>
                    <div className={styles.pdfViewerContainer} onClick={(e) => e.stopPropagation()}>
                        <iframe src="/documents/resume-ewerton.pdf" className={styles.pdfViewerFrame} title="Resume" />
                    </div>
                </div>
            )}
        </div>
    );
}
