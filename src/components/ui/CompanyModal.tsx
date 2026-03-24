"use client";

import React, { useEffect } from "react";
import Image from "next/image";
import { X } from "lucide-react";
import { useTranslations } from "next-intl";
import styles from "./CompanyModal.module.css";

interface Role {
    title: string;
    period: string;
    description: string;
}

interface CompanyModalProps {
    isOpen: boolean;
    onClose: () => void;
    company: {
        slug: string;
        name: string;
        logo: string;
        url: string;
        displayUrl: string;
    } | null;
}

const CompanyModal: React.FC<CompanyModalProps> = ({ isOpen, onClose, company }) => {
    const t = useTranslations("landing.partners.descriptions");

    useEffect(() => {
        if (isOpen) {
            document.body.style.overflow = "hidden";
        } else {
            document.body.style.overflow = "unset";
        }
        return () => {
            document.body.style.overflow = "unset";
        };
    }, [isOpen]);

    if (!isOpen || !company) return null;

    const getRoles = () => {
        const roles: Role[] = [];
        const slug = company.slug;

        if (slug === "stagwell") {
            roles.push({
                title: t("stagwell.role1"),
                period: t("stagwell.period1"),
                description: t("stagwell.desc1"),
            });
            roles.push({
                title: t("stagwell.role2"),
                period: t("stagwell.period2"),
                description: t("stagwell.desc2"),
            });
        } else {
            roles.push({
                title: t(`${slug}.role`),
                period: t(`${slug}.period`),
                description: t(`${slug}.desc`),
            });
        }
        return roles;
    };

    const roles = getRoles();

    return (
        <div className={styles.overlay} onClick={onClose}>
            <div className={styles.modal} onClick={(e) => e.stopPropagation()}>
                <button className={styles.closeButton} onClick={onClose}>
                    <X size={24} />
                </button>

                <div className={styles.content}>
                    <div className={styles.header}>
                        <div className={styles.logo}>
                            <Image src={company.logo} alt={company.name} width={64} height={64} />
                        </div>
                        <div className={styles.companyInfo}>
                            <h2>{company.name}</h2>
                            <a href={company.url} target="_blank" rel="noopener noreferrer" className={styles.companyUrl}>
                                {company.displayUrl}
                            </a>
                        </div>
                    </div>

                    <div className={styles.roles}>
                        {roles.map((role, index) => (
                            <div key={index} className={styles.roleItem}>
                                <div className={styles.roleHeader}>
                                    <h3 className={styles.roleTitle}>{role.title}</h3>
                                    <span className={styles.rolePeriod}>{role.period}</span>
                                </div>
                                <p className={styles.roleDescription}>{role.description}</p>
                            </div>
                        ))}
                    </div>
                </div>
            </div>
        </div>
    );
};

export default CompanyModal;
