package data

type Technology struct {
	Name     string
	Rating   float64
	Keywords []string
}

type TechnologyCategory struct {
	Slug     string
	Category string
	Icon     string // icon name for template lookup
	Skills   []Technology
}

var TechnologyCategories = []TechnologyCategory{
	{
		Slug:     "programming_languages",
		Category: "Programming Languages",
		Icon:     "code2",
		Skills: []Technology{
			{Name: "Python", Rating: 4.5, Keywords: []string{"py", "backend", "scripting", "ai", "ml", "data science", "automation", "pandas", "numpy", "django", "fastapi", "flask", "general purpose"}},
			{Name: "SQL", Rating: 4.5, Keywords: []string{"query", "database", "rdbms", "relational", "data", "joins", "aggregation", "stored procedures", "transactions", "normalization", "ddl", "dml"}},
			{Name: "JavaScript", Rating: 2.5, Keywords: []string{"js", "ecmascript", "frontend", "web", "browser", "typescript", "node", "dom", "async", "promises", "json", "dynamic"}},
			{Name: "Bash", Rating: 2.5, Keywords: []string{"shell", "terminal", "linux", "automation", "script", "cli", "unix", "commands", "pipes", "cron", "ssh", "sh"}},
			{Name: "Java", Rating: 2.0, Keywords: []string{"backend", "jvm", "enterprise", "spring", "maven", "oop", "multithreading", "microservices", "compiled"}},
			{Name: "PHP", Rating: 2.0, Keywords: []string{"web", "backend", "wordpress", "laravel", "symfony", "cms", "server-side", "scripting"}},
			{Name: "R", Rating: 2.0, Keywords: []string{"statistics", "data science", "analysis", "ggplot", "tidyverse", "statistical modeling", "visualization", "cran"}},
			{Name: "MATLAB", Rating: 2.0, Keywords: []string{"math", "engineering", "simulation", "numerical analysis", "signal processing", "scientific computing", "toolbox"}},
		},
	},
	{
		Slug:     "databases",
		Category: "Databases",
		Icon:     "database",
		Skills: []Technology{
			{Name: "MySQL", Rating: 4.0, Keywords: []string{"sql", "rdbms", "relational", "database", "innodb", "transactions", "indexing", "web", "open source", "mariadb"}},
			{Name: "PostgreSQL", Rating: 4.0, Keywords: []string{"postgres", "psql", "sql", "rdbms", "relational", "database", "jsonb", "extensions", "stored procedures", "open source", "transactions", "pg", "vector", "pgvector", "embeddings", "chromadb", "pinecone"}},
			{Name: "BigQuery", Rating: 4.0, Keywords: []string{"gbq", "google cloud", "gcp", "data warehouse", "dw", "analytics", "sql", "serverless", "columnar", "bi", "petabyte"}},
			{Name: "Spanner", Rating: 3.5, Keywords: []string{"google cloud", "gcp", "distributed", "sql", "rdbms", "database", "globally distributed", "transactions", "scalable", "strong consistency", "horizontally scalable"}},
			{Name: "ElasticSearch", Rating: 3.0, Keywords: []string{"elastic", "search", "indexing", "nosql", "elk", "full-text search", "kibana", "logstash", "aggregations", "json", "opensearch", "vector", "vector search", "knn", "embeddings", "chromadb", "pinecone"}},
			{Name: "CDC", Rating: 2.5, Keywords: []string{"change data capture", "replication", "streaming", "debezium", "kafka", "real-time", "sync", "event-driven", "incremental"}},
			{Name: "MongoDB", Rating: 2.0, Keywords: []string{"mongo", "nosql", "database", "document", "json", "atlas", "aggregation", "schemaless", "crud", "bson", "collections"}},
			{Name: "Redis", Rating: 2.0, Keywords: []string{"cache", "in-memory", "vss", "nosql", "key-value", "pub/sub", "queue", "session", "real-time", "caching layer", "sorted sets"}},
		},
	},
	{
		Slug:     "big_data_processing",
		Category: "Big Data & Data Processing",
		Icon:     "network",
		Skills: []Technology{
			{Name: "Apache Spark", Rating: 2.5, Keywords: []string{"spark", "big data", "processing", "scala", "pyspark", "analytics", "distributed computing", "rdd", "dataframe", "streaming", "ml"}},
			{Name: "Hadoop", Rating: 2.5, Keywords: []string{"big data", "mapreduce", "distributed", "hdfs", "hive", "hbase", "yarn", "batch processing", "storage", "cluster"}},
			{Name: "Apache Beam", Rating: 2.5, Keywords: []string{"beam", "dataflow", "streaming", "batch", "pipeline", "unified model", "gcp", "google cloud", "real-time", "etl"}},
		},
	},
	{
		Slug:     "data_science_ml",
		Category: "Data Science & ML",
		Icon:     "line-chart",
		Skills: []Technology{
			{Name: "Pandas", Rating: 4.5, Keywords: []string{"data analysis", "python", "dataframe", "cleaning", "wrangling", "manipulation", "csv", "excel", "etl", "tabular data"}},
			{Name: "NumPy", Rating: 4.5, Keywords: []string{"math", "python", "arrays", "scientific computing", "linear algebra", "matrices", "vectorization", "numerical", "ndarray"}},
			{Name: "Matplotlib", Rating: 4.0, Keywords: []string{"visualization", "python", "plotting", "charts", "graphs", "figures", "subplots", "pyplot", "data viz"}},
			{Name: "Seaborn", Rating: 3.5, Keywords: []string{"visualization", "python", "statistical", "charts", "heatmaps", "distributions", "correlation", "matplotlib", "data viz"}},
			{Name: "TensorFlow", Rating: 3.0, Keywords: []string{"tf", "machine learning", "deep learning", "ai", "neural networks", "keras", "gpu", "training", "inference", "model"}},
			{Name: "Scikit-learn", Rating: 2.5, Keywords: []string{"sklearn", "machine learning", "ml", "python", "statistics", "classification", "regression", "clustering", "preprocessing", "pipeline"}},
			{Name: "statsmodels", Rating: 2.0, Keywords: []string{"statistics", "python", "econometrics", "modeling", "regression", "time series", "hypothesis testing", "anova", "glm"}},
			{Name: "PyTorch", Rating: 2.0, Keywords: []string{"torch", "machine learning", "deep learning", "ai", "neural networks", "autograd", "gpu", "tensors", "training", "research"}},
		},
	},
	{
		Slug:     "web_development",
		Category: "Web Development",
		Icon:     "layout",
		Skills: []Technology{
			{Name: "HTML", Rating: 4.0, Keywords: []string{"web", "frontend", "markup", "structure", "html5", "semantic", "dom", "templates", "accessibility"}},
			{Name: "CSS", Rating: 4.0, Keywords: []string{"web", "frontend", "styling", "design", "layout", "css3", "flexbox", "grid", "animations", "responsive", "sass"}},
			{Name: "FastAPI", Rating: 4.0, Keywords: []string{"python", "api", "backend", "web", "rest", "async", "openapi", "swagger", "pydantic", "asyncio", "microservices"}},
			{Name: "Flask", Rating: 3.5, Keywords: []string{"python", "api", "backend", "web", "microframework", "jinja", "wsgi", "routes", "rest", "lightweight"}},
			{Name: "Bootstrap", Rating: 3.5, Keywords: []string{"css", "ui", "frontend", "responsive", "styling", "grid system", "components", "mobile-first", "framework"}},
			{Name: "Node.js", Rating: 3.0, Keywords: []string{"node", "nodejs", "js", "javascript", "backend", "server", "npm", "express", "api", "async", "event-driven"}},
			{Name: "Django", Rating: 2.5, Keywords: []string{"python", "web", "backend", "framework", "orm", "admin", "rest framework", "drf", "authentication", "migrations", "mtv"}},
			{Name: "React.js", Rating: 2.0, Keywords: []string{"react", "reactjs", "js", "frontend", "ui", "web", "hooks", "components", "spa", "jsx", "state management"}},
			{Name: "NextJS", Rating: 1.5, Keywords: []string{"next.js", "next", "react", "frontend", "ssr", "fullstack", "web", "ssg", "api routes", "typescript", "vercel", "routing"}},
		},
	},
	{
		Slug:     "data_engineering",
		Category: "Data Engineering",
		Icon:     "server",
		Skills: []Technology{
			{Name: "ETL", Rating: 4.0, Keywords: []string{"data pipeline", "integration", "extraction", "transformation", "loading", "elt", "batch", "real-time", "data flow", "orchestration", "ingestion"}},
			{Name: "Data Lakes", Rating: 3.5, Keywords: []string{"storage", "big data", "s3", "gcs", "lakehouse", "parquet", "delta lake", "iceberg", "raw data", "ingestion", "object storage"}},
			{Name: "Data Warehouses", Rating: 3.5, Keywords: []string{"bigquery", "snowflake", "redshift", "analytics", "dw", "olap", "columnar", "reporting", "star schema", "dimensional modeling"}},
			{Name: "Database Modelling", Rating: 3.5, Keywords: []string{"architecture", "schema", "design", "normalization", "star schema", "erd", "entity relationship", "oltp", "olap", "data modeling", "snowflake schema"}},
		},
	},
	{
		Slug:     "cloud_providers",
		Category: "Cloud Providers",
		Icon:     "cloud",
		Skills: []Technology{
			{Name: "Google Cloud", Rating: 4.0, Keywords: []string{"gcp", "google cloud platform", "cloud", "serverless", "computing", "bigquery", "cloud run", "gke", "cloud functions", "pub/sub", "gcs"}},
			{Name: "AWS", Rating: 3.0, Keywords: []string{"amazon web services", "amazon", "cloud", "s3", "ec2", "lambda", "rds", "dynamodb", "iam", "cloudwatch", "ecs", "eks"}},
			{Name: "Snowflake", Rating: 2.5, Keywords: []string{"cloud", "data warehouse", "analytics", "saas", "sql", "data sharing", "time travel", "virtual warehouse", "marketplace"}},
			{Name: "Microsoft Azure", Rating: 1.5, Keywords: []string{"azure", "microsoft", "cloud", "computing", "enterprise", "azure devops", "active directory", "blob storage", "aks", "functions", "ad"}},
		},
	},
	{
		Slug:     "containerization",
		Category: "Containerization",
		Icon:     "box",
		Skills: []Technology{
			{Name: "Docker", Rating: 4.0, Keywords: []string{"containers", "deployment", "devops", "infrastructure", "dockerfile", "compose", "images", "registry", "microservices", "isolation"}},
			{Name: "Kubernetes", Rating: 2.5, Keywords: []string{"k8s", "orchestration", "devops", "containers", "cloud", "pods", "helm", "scaling", "cluster", "deployments", "services"}},
		},
	},
	{
		Slug:     "devops_ci_cd",
		Category: "DevOps & CI/CD",
		Icon:     "settings",
		Skills: []Technology{
			{Name: "GitHub Actions", Rating: 3.5, Keywords: []string{"automation", "ci", "cd", "workflows", "git", "yaml", "runners", "pipelines", "testing", "deployment", "actions"}},
			{Name: "Jenkins", Rating: 2.0, Keywords: []string{"automation", "ci", "cd", "pipeline", "groovy", "jobs", "plugins", "build", "deploy", "open source"}},
		},
	},
	{
		Slug:     "version_control",
		Category: "Version Control",
		Icon:     "git-branch",
		Skills: []Technology{
			{Name: "Git", Rating: 4.5, Keywords: []string{"vcs", "version control", "terminal", "code", "commits", "branches", "merge", "rebase", "pull requests", "diff"}},
			{Name: "GitHub", Rating: 4.5, Keywords: []string{"git", "hosting", "collaboration", "open source", "pull requests", "issues", "actions", "repositories", "code review", "gh"}},
			{Name: "GitLab", Rating: 2.5, Keywords: []string{"git", "hosting", "ci", "cd", "devops", "merge requests", "pipelines", "runner", "self-hosted", "repositories"}},
		},
	},
	{
		Slug:     "dashboarding",
		Category: "Dashboarding",
		Icon:     "layout-dashboard",
		Skills: []Technology{
			{Name: "Looker Studio", Rating: 4.0, Keywords: []string{"datastudio", "data studio", "visualization", "bi", "google", "reporting", "charts", "kpi", "dashboard", "data connectors", "sharing"}},
			{Name: "Power BI", Rating: 2.5, Keywords: []string{"powerbi", "bi", "microsoft", "visualization", "reporting", "dax", "datasets", "service", "desktop", "analytics"}},
			{Name: "Redash", Rating: 2.5, Keywords: []string{"visualization", "sql", "dashboard", "bi", "open source", "queries", "charts", "sharing", "data sources"}},
		},
	},
	{
		Slug:     "agile",
		Category: "Agile",
		Icon:     "route",
		Skills: []Technology{
			{Name: "Jira", Rating: 4.0, Keywords: []string{"agile", "project management", "kanban", "scrum", "tasks", "issues", "sprints", "backlog", "roadmap", "atlassian"}},
			{Name: "Scrum", Rating: 4.0, Keywords: []string{"agile", "methodology", "sprints", "standups", "retrospectives", "sprint planning", "story points", "product owner", "ceremonies"}},
			{Name: "Confluence", Rating: 3.5, Keywords: []string{"documentation", "wiki", "collaboration", "knowledge base", "atlassian", "pages", "spaces", "team", "notes"}},
		},
	},
	{
		Slug:     "orchestration",
		Category: "Orchestration",
		Icon:     "list-tree",
		Skills: []Technology{
			{Name: "Apache Airflow", Rating: 4.0, Keywords: []string{"airflow", "orchestration", "workflow", "dag", "python", "data engineering", "scheduler", "tasks", "operators", "xcom", "etl", "cloud composer", "composer"}},
			{Name: "Zapier", Rating: 3.0, Keywords: []string{"automation", "integration", "low code", "workflow", "zaps", "triggers", "actions", "no code", "apps", "saas"}},
		},
	},
	{
		Slug:     "testing_frameworks",
		Category: "Testing Frameworks",
		Icon:     "activity",
		Skills: []Technology{
			{Name: "Pytest", Rating: 3.5, Keywords: []string{"python", "testing", "unit test", "automation", "fixtures", "mocking", "parametrize", "plugins", "coverage"}},
			{Name: "unittest", Rating: 3.5, Keywords: []string{"python", "testing", "unit test", "standard library", "mock", "assertions", "test cases", "test suite", "tdd"}},
			{Name: "Selenium", Rating: 3.0, Keywords: []string{"automation", "web", "browser", "testing", "qa", "webdriver", "e2e", "integration", "ui testing", "scraping"}},
			{Name: "TDD", Rating: 3.0, Keywords: []string{"test driven development", "agile", "quality", "clean code", "red-green-refactor", "unit tests", "bdd", "coverage", "design"}},
		},
	},
	{
		Slug:     "logging",
		Category: "Logging",
		Icon:     "activity",
		Skills: []Technology{
			{Name: "Cloud Logging", Rating: 2.5, Keywords: []string{"gcp", "google cloud", "monitoring", "observability", "logs", "log explorer", "alerts", "metrics", "audit", "stackdriver"}},
			{Name: "New Relic", Rating: 2.0, Keywords: []string{"monitoring", "apm", "observability", "performance", "tracing", "dashboards", "alerts", "logs", "distributed tracing", "infrastructure"}},
		},
	},
	{
		Slug:     "ads",
		Category: "Ads",
		Icon:     "megaphone",
		Skills: []Technology{
			{Name: "Facebook Ads", Rating: 3.0, Keywords: []string{"meta", "marketing", "advertising", "social media", "paid ads", "facebook", "instagram", "campaigns", "targeting", "pixel"}},
			{Name: "Google Ads", Rating: 3.0, Keywords: []string{"sem", "marketing", "advertising", "search", "paid ads", "adwords", "ppc", "campaigns", "keywords", "display"}},
			{Name: "Google Analytics", Rating: 2.5, Keywords: []string{"ga", "tracking", "web analytics", "data", "measurement", "ga4", "events", "conversions", "audience", "reports"}},
			{Name: "Google Ad Manager", Rating: 2.5, Keywords: []string{"gam", "advertising", "monetization", "publisher", "dfp", "ad units", "line items", "trafficking", "yield management"}},
		},
	},
	{
		Slug:     "os",
		Category: "OS",
		Icon:     "monitor",
		Skills: []Technology{
			{Name: "Windows", Rating: 4.0, Keywords: []string{"os", "operating system", "microsoft", "desktop", "powershell", "registry", "active directory", "wsl", "cmd"}},
			{Name: "Linux", Rating: 4.0, Keywords: []string{"os", "operating system", "bash", "terminal", "server", "unix", "ubuntu", "debian", "centos", "ssh", "permissions"}},
		},
	},
	{
		Slug:     "ai_tooling",
		Category: "AI Tooling",
		Icon:     "brain",
		Skills: []Technology{
			{Name: "Vertex AI", Rating: 4.0, Keywords: []string{"google cloud", "gcp", "ai", "ml", "machine learning", "llm", "generative ai", "model deployment", "vertex", "gemini", "mlops", "pipelines"}},
			{Name: "Claude Code", Rating: 4.0, Keywords: []string{"anthropic", "claude", "ai", "llm", "coding assistant", "agentic", "cli", "developer tools", "code generation", "ai assistant"}},
			{Name: "GitHub Copilot", Rating: 3.0, Keywords: []string{"github", "ai", "coding assistant", "llm", "code completion", "developer tools", "openai", "copilot", "autocomplete", "pair programming"}},
			{Name: "Factory AI", Rating: 3.0, Keywords: []string{"ai", "developer tools", "coding assistant", "automation", "software factory", "agents", "code generation", "engineering"}},
		},
	},
}
