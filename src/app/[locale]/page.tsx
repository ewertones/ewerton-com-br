"use client";

import { useTranslations, useLocale } from "next-intl";
import { useState, useEffect } from "react";
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
    Search,
} from "lucide-react";
import styles from "./page.module.css";
import Image from "next/image";

import { LINKEDIN_URL, AIOMOVER_URL, RESUME_PATH } from "@/constants";
import { StarRating } from "@/components/ui/StarRating";
import { TECHNOLOGY_CATEGORIES } from "@/data/technologies";
import { EDUCATION_DATA } from "@/data/education";
import { Calendar } from "lucide-react";

const HELLO_WORDS = [
    "Hello",
    "Olá",
    "Hola",
    "Bonjour",
    "Ciao",
    "Hej",
    "Hallo",
    "Salut",
    "Merhaba",
    "こんにちは",
    "Привет",
    "مرحبا",
];

const PARTNER_GROUPS = [
    {
        slug: "financial",
        partners: [
            { name: "DashFi", logo: "/partners/dashfi.webp", url: "https://dash.fi", displayUrl: "dash.fi" },
            { name: "CERC", logo: "/partners/cerc.webp", url: "https://cerc.com", displayUrl: "cerc.com" },
        ]
    },
    {
        slug: "marketing",
        partners: [
            { name: "Stagwell", logo: "/partners/stagwell.webp", url: "https://stagwellglobal.com", displayUrl: "stagwellglobal.com" },
            { name: "BeGrowth", logo: "/partners/begrowth.webp", url: "https://begrowth.com.br", displayUrl: "begrowth.com.br" },
        ]
    },
    {
        slug: "tax",
        partners: [
            { name: "EY", logo: "/partners/ey.webp", url: "https://ey.com", displayUrl: "ey.com" },
        ]
    },
    {
        slug: "food",
        partners: [
            { name: "Morsum", logo: "/partners/morsum.webp", url: "https://www.dotfoods.com", displayUrl: "dotfoods.com" },
        ]
    },
    {
        slug: "outsourcing",
        partners: [
            { name: "Softensity", logo: "/partners/softensity.webp", url: "https://softensity.com", displayUrl: "softensity.com" },
            { name: "Truelogic", logo: "/partners/truelogic.webp", url: "https://truelogic.io", displayUrl: "truelogic.io" },
            { name: "BlueShift", logo: "/partners/blueshift.webp", url: "https://blueshift.com.br", displayUrl: "blueshift.com.br" },
        ]
    }
];

export default function HomePage() {
    const t = useTranslations();
    const locale = useLocale();
    const [showPdfViewer, setShowPdfViewer] = useState(false);

    // Typewriter state
    const [helloIndex, setHelloIndex] = useState(0);
    const [displayedHello, setDisplayedHello] = useState(HELLO_WORDS[0]);
    const [isDeleting, setIsDeleting] = useState(false);

    // Technologies search state
    const [techSearch, setTechSearch] = useState("");

    const filteredTechnologies = TECHNOLOGY_CATEGORIES.map((cat) => {
        const query = techSearch.toLowerCase().trim();
        if (!query) return cat;

        const filteredSkills = cat.skills.filter(
            (s) =>
                s.name.toLowerCase().includes(query) ||
                (s.keywords && s.keywords.some((k) => k.toLowerCase().includes(query))),
        );
        return { ...cat, skills: filteredSkills };
    }).filter((cat) => cat.skills.length > 0 || cat.category.toLowerCase().includes(techSearch.toLowerCase().trim()));

    useEffect(() => {
        const currentWord = HELLO_WORDS[helloIndex];
        let delay: number;

        if (!isDeleting && displayedHello === currentWord) {
            // Fully typed — long pause before deleting
            delay = 2000;
            const timer = setTimeout(() => setIsDeleting(true), delay);
            return () => clearTimeout(timer);
        }

        if (isDeleting && displayedHello === "") {
            // Fully deleted — move to next word
            delay = 200;
            const timer = setTimeout(() => {
                setIsDeleting(false);
                setHelloIndex((i) => (i + 1) % HELLO_WORDS.length);
            }, delay);
            return () => clearTimeout(timer);
        }

        delay = isDeleting ? 30 : 60;
        const timer = setTimeout(() => {
            setDisplayedHello(
                isDeleting ? displayedHello.slice(0, -1) : currentWord.slice(0, displayedHello.length + 1),
            );
        }, delay);
        return () => clearTimeout(timer);
    }, [displayedHello, isDeleting, helloIndex]);

    return (
        <div className={styles.page}>
            <Navbar />

            <main className={styles.main}>
                {/* Hero + About merged */}
                <section className={styles.hero}>
                    {/* Centered animated greeting */}
                    <h1 className={styles.heroGreeting}>
                        <span className={styles.helloWord}>{displayedHello}</span>
                        <span className={styles.helloCursor} />
                        {", "}
                        {t("landing.hero.greeting")}
                    </h1>

                    {/* Two-column content */}
                    <div className={styles.heroContent}>
                        <div className={styles.heroText}>
                            <p className={styles.heroSubtitle}>{t("landing.hero.subtitle")}</p>

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
                        </div>

                        <div className={styles.heroPhoto}>
                            <Image
                                src="/photos/ewerton.webp"
                                alt="Ewerton"
                                width={280}
                                height={280}
                                className={styles.photo}
                            />
                        </div>
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

                {/* Technologies Section */}
                <section className={styles.section}>
                    <h2 className={styles.sectionTitle}>{t("landing.technologies.title")}</h2>
                    <p className={styles.sectionSubtitle}>{t("landing.technologies.subtitle")}</p>

                    <div className={styles.techSearchContainer}>
                        <Search className={styles.techSearchIcon} size={20} />
                        <input
                            type="text"
                            placeholder="Search skills (e.g. Postgres, React, AWS...)"
                            className={styles.techSearchInput}
                            value={techSearch}
                            onChange={(e) => setTechSearch(e.target.value)}
                        />
                    </div>

                    <div className={styles.techCategoriesGrid}>
                        {filteredTechnologies.map((cat) => (
                            <div key={cat.category} className={styles.techCategory}>
                                <div className={styles.techCategoryHeader}>
                                    <cat.icon size={20} className={styles.techCategoryIcon} />
                                    <h3 className={styles.techCategoryTitle}>{cat.category}</h3>
                                </div>
                                <div className={styles.techSkillsGrid}>
                                    {cat.skills.map((skill) => (
                                        <div key={skill.name} className={styles.techSkillCard}>
                                            <span className={styles.techSkillName}>{skill.name}</span>
                                            <StarRating rating={skill.rating} size={14} />
                                        </div>
                                    ))}
                                </div>
                            </div>
                        ))}
                    </div>
                </section>

                {/* About & Education Section */}
                <section className={styles.section}>
                    <h2 className={styles.sectionTitle}>{t("landing.about.title")}</h2>
                    <p className={styles.sectionSubtitle}>{t("landing.about.subtitle")}</p>

                    <div className={styles.eduGrid}>
                        {EDUCATION_DATA[locale.startsWith("pt") ? "pt" : "en"].map((edu) => (
                            <div key={edu.id} className={styles.eduCard}>
                                <div className={styles.eduHeader}>
                                    <a
                                        href={edu.websiteUrl}
                                        target="_blank"
                                        rel="noopener noreferrer"
                                        className={styles.eduLogo}
                                    >
                                        <Image src={edu.logo} alt={edu.institution} width={40} height={40} />
                                    </a>
                                    <div className={styles.eduDetails}>
                                        <a
                                            href={edu.websiteUrl}
                                            target="_blank"
                                            rel="noopener noreferrer"
                                            className={styles.eduInstLink}
                                        >
                                            <h3 className={styles.eduInst}>{edu.institution}</h3>
                                        </a>
                                        <div className={styles.eduDegreeWrapper}>
                                            {edu.degree.includes(",") ?
                                                <>
                                                    <span className={styles.eduDegreeType}>
                                                        {edu.degree.split(",")[0].trim()}
                                                    </span>
                                                    <span className={styles.eduDegreeName}>
                                                        {edu.degree.split(",").slice(1).join(",").trim()}
                                                    </span>
                                                </>
                                            :   <span className={styles.eduDegreeName}>{edu.degree}</span>}
                                        </div>
                                        <p className={styles.eduMeta}>
                                            {edu.period} • {edu.grade}
                                        </p>
                                    </div>
                                </div>
                                <div className={styles.eduContent}>
                                    <ul className={styles.eduPoints}>
                                        {edu.points.map((pt, i) => (
                                            <li key={i}>{pt}</li>
                                        ))}
                                    </ul>
                                </div>
                            </div>
                        ))}
                    </div>
                </section>

                {/* Partners Section */}
                <section className={styles.section}>
                    <h2 className={styles.sectionTitle}>{t("landing.partners.title")}</h2>
                    <p className={styles.sectionSubtitle}>{t("landing.partners.subtitle")}</p>

                    <div className={styles.partnerGroups}>
                        {PARTNER_GROUPS.map((group) => (
                            <div key={group.slug} className={styles.partnerGroup}>
                                <h3 className={styles.partnerGroupTitle}>
                                    {t(`landing.partners.categories.${group.slug}`)}
                                </h3>
                                <div className={styles.partnersGrid}>
                                    {group.partners.map((partner) => (
                                        <a
                                            key={partner.name}
                                            href={partner.url}
                                            target="_blank"
                                            rel="noopener noreferrer"
                                            className={styles.partnerCard}
                                        >
                                            <div className={styles.partnerLogo}>
                                                <Image src={partner.logo} alt={partner.name} width={72} height={72} />
                                            </div>
                                            <p className={styles.partnerName}>{partner.name}</p>
                                            <p className={styles.partnerUrl}>{partner.displayUrl}</p>
                                        </a>
                                    ))}
                                </div>
                            </div>
                        ))}
                    </div>
                </section>

                {/* Closing CTA */}
                <div className={styles.closingCtaCard}>
                    <div className={styles.closingCtaContent}>
                        <h2 className={styles.closingCtaTitle}>{t("landing.closingCta.title")}</h2>
                        <p className={styles.closingCtaSubtitle}>{t("landing.closingCta.subtitle")}</p>
                    </div>
                    <div className={styles.closingCtaButton}>
                        <a href="https://calendly.com/ewertones/30min" target="_blank" rel="noopener noreferrer">
                            <Button variant="primary" size="lg" leftIcon={<Calendar size={20} />}>
                                {t("landing.closingCta.button")}
                            </Button>
                        </a>
                    </div>
                </div>
            </main>

            <Footer />

            {/* PDF Viewer Modal */}
            {showPdfViewer && (
                <div className={styles.pdfViewerOverlay} onClick={() => setShowPdfViewer(false)}>
                    <div className={styles.pdfViewerContainer} onClick={(e) => e.stopPropagation()}>
                        <iframe src={RESUME_PATH} className={styles.pdfViewerFrame} title="Resume" />
                    </div>
                </div>
            )}
        </div>
    );
}
