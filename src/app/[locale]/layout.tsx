import { NextIntlClientProvider } from "next-intl";
import { GoogleAnalytics } from "@next/third-parties/google";
import { getMessages, getTranslations } from "next-intl/server";
import { notFound } from "next/navigation";
import { routing } from "@/i18n/routing";
import { ThemeProvider } from "@/lib/ThemeProvider";

import { ReactNode } from "react";
import { Inter } from "next/font/google";
import { Metadata } from "next";

const inter = Inter({
    subsets: ["latin"],
    variable: "--font-inter",
    display: "swap",
});

export async function generateMetadata({ params }: Omit<Props, "children">): Promise<Metadata> {
    const { locale } = await params;
    const t = await getTranslations({ locale, namespace: "metadata" });

    return {
        title: t("title"),
        description: "I'm a software engineer dedicated to building innovative solutions that make a difference.",
        openGraph: {
            locale: locale === "pt-br" ? "pt_BR" : "en_US",
            type: "website",
        },
        alternates: {
            canonical: "https://ewerton.com.br",
            languages: {
                en: "https://ewerton.com.br/en",
                "pt-br": "https://ewerton.com.br/pt-br",
            },
        },
    };
}

const jsonLd = {
    "@context": "https://schema.org",
    "@type": "Person",
    name: "Ewerton",
    url: "https://ewerton.com.br",
    image: "https://ewerton.com.br/logo.svg",
    description: "Software engineer improving the world through technology.",
};

export function generateStaticParams() {
    return routing.locales.map((locale) => ({ locale }));
}

type Props = {
    children: ReactNode;
    params: Promise<{ locale: string }>;
};

export default async function LocaleLayout({ children, params }: Props) {
    const { locale } = await params;

    if (!routing.locales.includes(locale as "en" | "pt-br")) {
        notFound();
    }

    const messages = await getMessages();

    return (
        <html lang={locale} suppressHydrationWarning>
            <head>
                <link rel="alternate" hrefLang="en" href="https://ewerton.com.br/en" />
                <link rel="alternate" hrefLang="pt-BR" href="https://ewerton.com.br/pt-br" />
                <link rel="alternate" hrefLang="x-default" href="https://ewerton.com.br/en" />
                <script
                    dangerouslySetInnerHTML={{
                        __html: `
              (function() {
                try {
                  var theme = localStorage.getItem('theme');
                  if (!theme) {
                    theme = 'light';
                  }
                  document.documentElement.setAttribute('data-theme', theme);
                } catch (e) {}
              })();
            `,
                    }}
                />
            </head>
            <body className={inter.className}>
                <ThemeProvider>
                    <NextIntlClientProvider messages={messages}>{children}</NextIntlClientProvider>
                </ThemeProvider>
                <GoogleAnalytics gaId="G-MFH72L1LCE" />
                <script type="application/ld+json" dangerouslySetInnerHTML={{ __html: JSON.stringify(jsonLd) }} />
            </body>
        </html>
    );
}
