package data

type EducationEntry struct {
	ID          string
	Logo        string
	Institution string
	Degree      string
	Period      string
	Grade       string
	Points      []string
	Skills      string
	WebsiteURL  string
}

var EducationData = map[string][]EducationEntry{
	"en": {
		{
			ID:          "usp",
			Logo:        "/static/images/partners/usp.webp",
			Institution: "USP - Universidade de São Paulo",
			WebsiteURL:  "https://www5.usp.br/",
			Degree:      "Master in Business Administration (MBA), Data Science and Analytics",
			Period:      "2021 – 2022",
			Grade:       "Grade: 99/100",
			Points:      []string{"Gained a strong foundation in analytical and data skills, with a focus on applying them within a business context."},
			Skills:      "Skills: Google Cloud Platform (GCP), Data Engineering, and more.",
		},
		{
			ID:          "puc-grad",
			Logo:        "/static/images/partners/pucminas.webp",
			Institution: "PUC Minas",
			WebsiteURL:  "https://www.pucminas.br/",
			Degree:      "Graduate Degree, Computer Software Engineering",
			Period:      "2021 – 2022",
			Grade:       "Grade: 93/100",
			Points:      []string{"Agile methodologies and processes, such as Scrum and Kanban, which are commonly used for software development and project management."},
			Skills:      "Skills: Software Architecture, DevOps, Agile.",
		},
		{
			ID:          "puc-undergrad",
			Logo:        "/static/images/partners/pucminas.webp",
			Institution: "PUC Minas",
			WebsiteURL:  "https://www.pucminas.br/",
			Degree:      "Undergraduate degree, Chemical Engineering",
			Period:      "2014 – 2019",
			Grade:       "Grade: 85/100",
			Points:      []string{"Gained a strong understanding of the principles of chemical engineering, including thermodynamics, kinetics, transport phenomena, and reactor design."},
			Skills:      "Skills: Chemical Processes, Thermodynamics, Optimization.",
		},
	},
	"pt": {
		{
			ID:          "usp",
			Logo:        "/static/images/partners/usp.webp",
			Institution: "USP - Universidade de São Paulo",
			WebsiteURL:  "https://www5.usp.br/",
			Degree:      "Master in Business Administration (MBA), Ciência de Dados e Analytics",
			Period:      "2021 – 2022",
			Grade:       "Nota: 99/100",
			Points:      []string{"Obtive uma forte base analítica, focada na aplicação de dados no contexto de negócios."},
			Skills:      "Habilidades: Google Cloud Platform (GCP), Engenharia de Dados, e mais.",
		},
		{
			ID:          "puc-grad",
			Logo:        "/static/images/partners/pucminas.webp",
			Institution: "PUC Minas",
			WebsiteURL:  "https://www.pucminas.br/",
			Degree:      "Pós-Graduação, Engenharia de Software",
			Period:      "2021 – 2022",
			Grade:       "Nota: 93/100",
			Points:      []string{"Metodologias ágeis e processos, como Scrum e Kanban, comumente usados para desenvolvimento de software e gestão de projetos."},
			Skills:      "Habilidades: Arquitetura de Software, DevOps, Agile.",
		},
		{
			ID:          "puc-undergrad",
			Logo:        "/static/images/partners/pucminas.webp",
			Institution: "PUC Minas",
			WebsiteURL:  "https://www.pucminas.br/",
			Degree:      "Graduação, Engenharia Química",
			Period:      "2014 – 2019",
			Grade:       "Nota: 85/100",
			Points:      []string{"Ganhei profundo entendimento em princípios de engenharia química (termodinâmica, cinética, fenômenos de transporte, design de reatores)."},
			Skills:      "Habilidades: Processos Químicos, Termodinâmica, Otimização.",
		},
	},
}
