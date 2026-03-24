import { 
    Code2, 
    Database, 
    Network, 
    LineChart, 
    Layout, 
    Server, 
    Cloud, 
    Box, 
    Settings, 
    GitBranch, 
    LayoutDashboard, 
    Route, 
    ListTree, 
    Activity, 
    Megaphone, 
    Monitor,
    type LucideIcon
} from "lucide-react";

export interface Technology {
    name: string;
    rating: number; // 0.0 to 5.0
    keywords?: string[];
}

export interface TechnologyCategory {
    category: string;
    icon: LucideIcon;
    skills: Technology[];
}

export const TECHNOLOGY_CATEGORIES: TechnologyCategory[] = [
    {
        category: "Programming Languages",
        icon: Code2,
        skills: [
            { name: "Python", rating: 4.5, keywords: ["py", "backend", "scripting", "ai", "ml", "data science", "automation", "pandas", "numpy", "django", "fastapi", "flask", "general purpose"] },
            { name: "SQL", rating: 4.5, keywords: ["query", "database", "rdbms", "relational", "data", "joins", "aggregation", "stored procedures", "transactions", "normalization", "ddl", "dml"] },
            { name: "JavaScript", rating: 2.5, keywords: ["js", "ecmascript", "frontend", "web", "browser", "typescript", "node", "dom", "async", "promises", "json", "dynamic"] },
            { name: "Bash", rating: 2.5, keywords: ["shell", "terminal", "linux", "automation", "script", "cli", "unix", "commands", "pipes", "cron", "ssh", "sh"] },
            { name: "Java", rating: 2.0, keywords: ["backend", "jvm", "enterprise", "spring", "maven", "oop", "multithreading", "microservices", "compiled"] },
            { name: "PHP", rating: 2.0, keywords: ["web", "backend", "wordpress", "laravel", "symfony", "cms", "server-side", "scripting"] },
            { name: "R", rating: 2.0, keywords: ["statistics", "data science", "analysis", "ggplot", "tidyverse", "statistical modeling", "visualization", "cran"] },
            { name: "MATLAB", rating: 2.0, keywords: ["math", "engineering", "simulation", "numerical analysis", "signal processing", "scientific computing", "toolbox"] },
        ]
    },
    {
        category: "Database Technologies",
        icon: Database,
        skills: [
            { name: "MySQL", rating: 4.0, keywords: ["sql", "rdbms", "relational", "database", "innodb", "transactions", "indexing", "web", "open source", "mariadb"] },
            { name: "PostgreSQL", rating: 4.0, keywords: ["postgres", "psql", "sql", "rdbms", "relational", "database", "jsonb", "extensions", "stored procedures", "open source", "transactions", "pg"] },
            { name: "BigQuery", rating: 4.0, keywords: ["gbq", "google cloud", "gcp", "data warehouse", "dw", "analytics", "sql", "serverless", "columnar", "bi", "petabyte"] },
            { name: "Spanner", rating: 3.5, keywords: ["google cloud", "gcp", "distributed", "sql", "rdbms", "database", "globally distributed", "transactions", "scalable", "strong consistency", "horizontally scalable"] },
            { name: "ElasticSearch", rating: 3.0, keywords: ["elastic", "search", "indexing", "nosql", "elk", "full-text search", "kibana", "logstash", "aggregations", "json", "opensearch"] },
            { name: "CDC", rating: 2.5, keywords: ["change data capture", "replication", "streaming", "debezium", "kafka", "real-time", "sync", "event-driven", "incremental"] },
            { name: "MongoDB", rating: 2.0, keywords: ["mongo", "nosql", "database", "document", "json", "atlas", "aggregation", "schemaless", "crud", "bson", "collections"] },
            { name: "Redis", rating: 2.0, keywords: ["cache", "in-memory", "vss", "nosql", "key-value", "pub/sub", "queue", "session", "real-time", "caching layer", "sorted sets"] },
        ]
    },
    {
        category: "Big Data & Data Processing",
        icon: Network,
        skills: [
            { name: "Apache Spark", rating: 2.5, keywords: ["spark", "big data", "processing", "scala", "pyspark", "analytics", "distributed computing", "rdd", "dataframe", "streaming", "ml"] },
            { name: "Hadoop", rating: 2.5, keywords: ["big data", "mapreduce", "distributed", "hdfs", "hive", "hbase", "yarn", "batch processing", "storage", "cluster"] },
            { name: "Apache Beam", rating: 2.5, keywords: ["beam", "dataflow", "streaming", "batch", "pipeline", "unified model", "gcp", "google cloud", "real-time", "etl"] },
        ]
    },
    {
        category: "Data Science & ML",
        icon: LineChart,
        skills: [
            { name: "Pandas", rating: 4.5, keywords: ["data analysis", "python", "dataframe", "cleaning", "wrangling", "manipulation", "csv", "excel", "etl", "tabular data"] },
            { name: "NumPy", rating: 4.5, keywords: ["math", "python", "arrays", "scientific computing", "linear algebra", "matrices", "vectorization", "numerical", "ndarray"] },
            { name: "Matplotlib", rating: 4.0, keywords: ["visualization", "python", "plotting", "charts", "graphs", "figures", "subplots", "pyplot", "data viz"] },
            { name: "Seaborn", rating: 3.5, keywords: ["visualization", "python", "statistical", "charts", "heatmaps", "distributions", "correlation", "matplotlib", "data viz"] },
            { name: "TensorFlow", rating: 3.0, keywords: ["tf", "machine learning", "deep learning", "ai", "neural networks", "keras", "gpu", "training", "inference", "model"] },
            { name: "Scikit-learn", rating: 2.5, keywords: ["sklearn", "machine learning", "ml", "python", "statistics", "classification", "regression", "clustering", "preprocessing", "pipeline"] },
            { name: "statsmodels", rating: 2.0, keywords: ["statistics", "python", "econometrics", "modeling", "regression", "time series", "hypothesis testing", "anova", "glm"] },
            { name: "PyTorch", rating: 2.0, keywords: ["torch", "machine learning", "deep learning", "ai", "neural networks", "autograd", "gpu", "tensors", "training", "research"] },
        ]
    },
    {
        category: "Web Development",
        icon: Layout,
        skills: [
            { name: "HTML", rating: 4.0, keywords: ["web", "frontend", "markup", "structure", "html5", "semantic", "dom", "templates", "accessibility"] },
            { name: "CSS", rating: 4.0, keywords: ["web", "frontend", "styling", "design", "layout", "css3", "flexbox", "grid", "animations", "responsive", "sass"] },
            { name: "FastAPI", rating: 4.0, keywords: ["python", "api", "backend", "web", "rest", "async", "openapi", "swagger", "pydantic", "asyncio", "microservices"] },
            { name: "Flask", rating: 3.5, keywords: ["python", "api", "backend", "web", "microframework", "jinja", "wsgi", "routes", "rest", "lightweight"] },
            { name: "Bootstrap", rating: 3.5, keywords: ["css", "ui", "frontend", "responsive", "styling", "grid system", "components", "mobile-first", "framework"] },
            { name: "Node.js", rating: 3.0, keywords: ["node", "nodejs", "js", "javascript", "backend", "server", "npm", "express", "api", "async", "event-driven"] },
            { name: "Django", rating: 2.5, keywords: ["python", "web", "backend", "framework", "orm", "admin", "rest framework", "drf", "authentication", "migrations", "mtv"] },
            { name: "React.js", rating: 2.0, keywords: ["react", "reactjs", "js", "frontend", "ui", "web", "hooks", "components", "spa", "jsx", "state management"] },
            { name: "NextJS", rating: 1.5, keywords: ["next.js", "next", "react", "frontend", "ssr", "fullstack", "web", "ssg", "api routes", "typescript", "vercel", "routing"] },
        ]
    },
    {
        category: "Data Engineering",
        icon: Server,
        skills: [
            { name: "ETL", rating: 4.0, keywords: ["data pipeline", "integration", "extraction", "transformation", "loading", "elt", "batch", "real-time", "data flow", "orchestration", "ingestion"] },
            { name: "Data Lakes", rating: 3.5, keywords: ["storage", "big data", "s3", "gcs", "lakehouse", "parquet", "delta lake", "iceberg", "raw data", "ingestion", "object storage"] },
            { name: "Data Warehouses", rating: 3.5, keywords: ["bigquery", "snowflake", "redshift", "analytics", "dw", "olap", "columnar", "reporting", "star schema", "dimensional modeling"] },
            { name: "Database Modelling", rating: 3.5, keywords: ["architecture", "schema", "design", "normalization", "star schema", "erd", "entity relationship", "oltp", "olap", "data modeling", "snowflake schema"] },
        ]
    },
    {
        category: "Cloud Providers",
        icon: Cloud,
        skills: [
            { name: "GCP", rating: 4.0, keywords: ["google cloud platform", "google cloud", "cloud", "serverless", "computing", "bigquery", "cloud run", "gke", "cloud functions", "pub/sub", "gcs"] },
            { name: "AWS", rating: 3.0, keywords: ["amazon web services", "amazon", "cloud", "s3", "ec2", "lambda", "rds", "dynamodb", "iam", "cloudwatch", "ecs", "eks"] },
            { name: "Snowflake", rating: 2.5, keywords: ["cloud", "data warehouse", "analytics", "saas", "sql", "data sharing", "time travel", "virtual warehouse", "marketplace"] },
            { name: "Azure", rating: 1.5, keywords: ["microsoft", "cloud", "computing", "enterprise", "azure devops", "active directory", "blob storage", "aks", "functions", "ad"] },
        ]
    },
    {
        category: "Containerization",
        icon: Box,
        skills: [
            { name: "Docker", rating: 4.0, keywords: ["containers", "deployment", "devops", "infrastructure", "dockerfile", "compose", "images", "registry", "microservices", "isolation"] },
            { name: "Kubernetes", rating: 2.5, keywords: ["k8s", "orchestration", "devops", "containers", "cloud", "pods", "helm", "scaling", "cluster", "deployments", "services"] },
        ]
    },
    {
        category: "DevOps & CI/CD",
        icon: Settings,
        skills: [
            { name: "GitHub Actions", rating: 3.5, keywords: ["automation", "ci", "cd", "workflows", "git", "yaml", "runners", "pipelines", "testing", "deployment", "actions"] },
            { name: "Jenkins", rating: 2.0, keywords: ["automation", "ci", "cd", "pipeline", "groovy", "jobs", "plugins", "build", "deploy", "open source"] },
        ]
    },
    {
        category: "Version Control",
        icon: GitBranch,
        skills: [
            { name: "Git", rating: 4.5, keywords: ["vcs", "version control", "terminal", "code", "commits", "branches", "merge", "rebase", "pull requests", "diff"] },
            { name: "GitHub", rating: 4.5, keywords: ["git", "hosting", "collaboration", "open source", "pull requests", "issues", "actions", "repositories", "code review", "gh"] },
            { name: "GitLab", rating: 2.5, keywords: ["git", "hosting", "ci", "cd", "devops", "merge requests", "pipelines", "runner", "self-hosted", "repositories"] },
        ]
    },
    {
        category: "Dashboarding",
        icon: LayoutDashboard,
        skills: [
            { name: "Looker Studio", rating: 4.0, keywords: ["datastudio", "data studio", "visualization", "bi", "google", "reporting", "charts", "kpi", "dashboard", "data connectors", "sharing"] },
            { name: "Power BI", rating: 2.5, keywords: ["powerbi", "bi", "microsoft", "visualization", "reporting", "dax", "datasets", "service", "desktop", "analytics"] },
            { name: "Redash", rating: 2.5, keywords: ["visualization", "sql", "dashboard", "bi", "open source", "queries", "charts", "sharing", "data sources"] },
        ]
    },
    {
        category: "Agile",
        icon: Route,
        skills: [
            { name: "Jira", rating: 4.0, keywords: ["agile", "project management", "kanban", "scrum", "tasks", "issues", "sprints", "backlog", "roadmap", "atlassian"] },
            { name: "Scrum", rating: 4.0, keywords: ["agile", "methodology", "sprints", "standups", "retrospectives", "sprint planning", "story points", "product owner", "ceremonies"] },
            { name: "Confluence", rating: 3.5, keywords: ["documentation", "wiki", "collaboration", "knowledge base", "atlassian", "pages", "spaces", "team", "notes"] },
        ]
    },
    {
        category: "Orchestration",
        icon: ListTree,
        skills: [
            { name: "Apache Airflow", rating: 4.0, keywords: ["airflow", "orchestration", "workflow", "dag", "python", "data engineering", "scheduler", "tasks", "operators", "xcom", "etl"] },
            { name: "Zapier", rating: 3.0, keywords: ["automation", "integration", "low code", "workflow", "zaps", "triggers", "actions", "no code", "apps", "saas"] },
        ]
    },
    {
        category: "Testing Frameworks",
        icon: Activity,
        skills: [
            { name: "Pytest", rating: 3.5, keywords: ["python", "testing", "unit test", "automation", "fixtures", "mocking", "parametrize", "plugins", "coverage"] },
            { name: "unittest", rating: 3.5, keywords: ["python", "testing", "unit test", "standard library", "mock", "assertions", "test cases", "test suite", "tdd"] },
            { name: "Selenium", rating: 3.0, keywords: ["automation", "web", "browser", "testing", "qa", "webdriver", "e2e", "integration", "ui testing", "scraping"] },
            { name: "TDD", rating: 2.5, keywords: ["test driven development", "agile", "quality", "clean code", "red-green-refactor", "unit tests", "bdd", "coverage", "design"] },
        ]
    },
    {
        category: "Logging",
        icon: Activity, // Reuse generic monitoring activity
        skills: [
            { name: "Cloud Logging", rating: 2.5, keywords: ["gcp", "google cloud", "monitoring", "observability", "logs", "log explorer", "alerts", "metrics", "audit", "stackdriver"] },
            { name: "New Relic", rating: 2.0, keywords: ["monitoring", "apm", "observability", "performance", "tracing", "dashboards", "alerts", "logs", "distributed tracing", "infrastructure"] },
        ]
    },
    {
        category: "Ads",
        icon: Megaphone,
        skills: [
            { name: "Facebook Ads", rating: 3.0, keywords: ["meta", "marketing", "advertising", "social media", "paid ads", "facebook", "instagram", "campaigns", "targeting", "pixel"] },
            { name: "Google Ads", rating: 3.0, keywords: ["sem", "marketing", "advertising", "search", "paid ads", "adwords", "ppc", "campaigns", "keywords", "display"] },
            { name: "Google Analytics", rating: 2.5, keywords: ["ga", "tracking", "web analytics", "data", "measurement", "ga4", "events", "conversions", "audience", "reports"] },
            { name: "Google Ad Manager", rating: 2.5, keywords: ["gam", "advertising", "monetization", "publisher", "dfp", "ad units", "line items", "trafficking", "yield management"] },
        ]
    },
    {
        category: "OS",
        icon: Monitor,
        skills: [
            { name: "Windows", rating: 4.0, keywords: ["os", "operating system", "microsoft", "desktop", "powershell", "registry", "active directory", "wsl", "cmd"] },
            { name: "Linux", rating: 4.0, keywords: ["os", "operating system", "bash", "terminal", "server", "unix", "ubuntu", "debian", "centos", "ssh", "permissions"] },
        ]
    }
];
