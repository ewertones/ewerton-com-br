import { ReactNode } from "react";
import { Metadata, Viewport } from "next";
import "@/app/globals.css";

export const viewport: Viewport = {
    width: "device-width",
    initialScale: 1,
};

export const metadata: Metadata = {
    metadataBase: new URL(process.env.NEXT_PUBLIC_APP_URL || "https://ewerton.com.br"),
    title: "ewerton",
    description: "ewerton is a software house dedicated to building innovative solutions that make a difference.",
    openGraph: {
        title: "ewerton",
        description: "Improving the world through technology.",
        url: "https://ewerton.com.br",
        siteName: "ewerton",
        type: "website",
    },
    twitter: {
        card: "summary_large_image",
        title: "ewerton",
        description: "Improving the world through technology.",
    },
    icons: {
        icon: [{ url: "/logo.svg", type: "image/svg+xml" }],
    },
};

type Props = {
    children: ReactNode;
};

export default function RootLayout({ children }: Props) {
    return children;
}
