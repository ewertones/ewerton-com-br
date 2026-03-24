export interface EducationEntry {
    id: string;
    logo: string;
    institution: string;
    degree: string;
    period: string;
    grade: string;
    points: string[];
    skills: string;
    websiteUrl: string;
}

export const EDUCATION_DATA: { en: EducationEntry[], pt: EducationEntry[] } = {
    en: [
        {
            id: "usp",
            logo: "/partners/usp.webp", // Fallback generated or added later if missing
            institution: "USP - Universidade de São Paulo",
            websiteUrl: "https://www5.usp.br/",
            degree: "Master in Business Administration (MBA), Data Science and Analytics",
            period: "2021 – 2022",
            grade: "Grade: 99/100",
            points: [
                "Gained a strong foundation in analytical and data skills, with a focus on applying them within a business context."
            ],
            skills: "Skills: Google Cloud Platform (GCP), Data Engineering, and more."
        },
        {
            id: "puc-grad",
            logo: "/partners/pucminas.webp",
            institution: "PUC Minas",
            websiteUrl: "https://www.pucminas.br/",
            degree: "Graduate Degree, Computer Software Engineering",
            period: "2021 – 2022",
            grade: "Grade: 93/100",
            points: [
                "Agile methodologies and processes, such as Scrum and Kanban, which are commonly used for software development and project management."
            ],
            skills: "Skills: Software Architecture, DevOps, Agile."
        },
        {
            id: "puc-undergrad",
            logo: "/partners/pucminas.webp",
            institution: "PUC Minas",
            websiteUrl: "https://www.pucminas.br/",
            degree: "Undergraduate degree, Chemical Engineering",
            period: "2014 – 2019",
            grade: "Grade: 85/100",
            points: [
                "Gained a strong understanding of the principles of chemical engineering, including thermodynamics, kinetics, transport phenomena, and reactor design."
            ],
            skills: "Skills: Chemical Processes, Thermodynamics, Optimization."
        }
    ],
    pt: [
        {
            id: "usp",
            logo: "/partners/usp.webp",
            institution: "USP - Universidade de São Paulo",
            websiteUrl: "https://www5.usp.br/",
            degree: "Master in Business Administration (MBA), Ciência de Dados e Analytics",
            period: "2021 – 2022",
            grade: "Nota: 99/100",
            points: [
                "Obtive uma forte base analítica, focada na aplicação de dados no contexto de negócios."
            ],
            skills: "Habilidades: Google Cloud Platform (GCP), Engenharia de Dados, e mais."
        },
        {
            id: "puc-grad",
            logo: "/partners/pucminas.webp",
            institution: "PUC Minas",
            websiteUrl: "https://www.pucminas.br/",
            degree: "Pós-Graduação, Engenharia de Software",
            period: "2021 – 2022",
            grade: "Nota: 93/100",
            points: [
                "Metodologias ágeis e processos, como Scrum e Kanban, comumente usados para desenvolvimento de software e gestão de projetos."
            ],
            skills: "Habilidades: Arquitetura de Software, DevOps, Agile."
        },
        {
            id: "puc-undergrad",
            logo: "/partners/pucminas.webp",
            institution: "PUC Minas",
            websiteUrl: "https://www.pucminas.br/",
            degree: "Graduação, Engenharia Química",
            period: "2014 – 2019",
            grade: "Nota: 85/100",
            points: [
                "Ganhei profundo entendimento em princípios de engenharia química (termodinâmica, cinética, fenômenos de transporte, design de reatores)."
            ],
            skills: "Habilidades: Processos Químicos, Termodinâmica, Otimização."
        }
    ]
};
