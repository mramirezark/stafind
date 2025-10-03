-- Insert comprehensive modern skills (2025) with normalized categories

-- First, insert all categories
INSERT INTO categories (name) VALUES
('Programming Language'),
('Frontend Framework'),
('Backend Framework'),
('Database'),
('Cloud Platform'),
('DevOps'),
('API'),
('Version Control'),
('Architecture'),
('Frontend Technology'),
('Build Tool'),
('Testing Framework'),
('Mobile Development'),
('Desktop Development'),
('Data Science & AI'),
('Security'),
('Soft Skills');

-- Insert skills
INSERT INTO skills (name) VALUES
-- Programming Languages
('JavaScript'),
('Python'),
('Go'),
('TypeScript'),
('Java'),
('C++'),
('C#'),
('Rust'),
('Swift'),
('Kotlin'),
('PHP'),
('Ruby'),
('Scala'),
('R'),
('MATLAB'),
('Dart'),
('Lua'),
('Haskell'),
('Clojure'),
('Erlang'),
('Elixir'),
('F#'),
('Assembly'),
('C'),
('Shell'),
('PowerShell'),
('Bash'),

-- Frontend Frameworks
('React'),
('Vue.js'),
('Angular'),
('Next.js'),
('Nuxt.js'),
('Svelte'),
('SvelteKit'),
('Gatsby'),
('Astro'),
('Remix'),
('SolidJS'),
('Alpine.js'),
('jQuery'),
('Lit'),
('Preact'),

-- Backend Frameworks
('Node.js'),
('Express.js'),
('FastAPI'),
('Django'),
('Flask'),
('Spring Boot'),
('Laravel'),
('Symfony'),
('Rails'),
('ASP.NET'),
('Deno'),
('Bun'),
('NestJS'),
('Koa'),
('Hapi'),
('Sails'),
('Meteor'),
('Feathers'),
('AdonisJS'),
('Fastify'),
('Gin'),
('Echo'),
('Fiber'),

-- Databases
('PostgreSQL'),
('MongoDB'),
('Redis'),
('MySQL'),
('SQLite'),
('Oracle'),
('SQL Server'),
('MariaDB'),
('Neo4j'),
('CouchDB'),
('RethinkDB'),
('InfluxDB'),
('TimescaleDB'),
('ClickHouse'),
('Snowflake'),
('BigQuery'),
('Redshift'),
('Firebase'),
('Supabase'),
('PlanetScale'),
('CockroachDB'),
('FaunaDB'),
('ArangoDB'),
('OrientDB'),
('Amazon RDS'),
('Azure SQL'),
('DynamoDB'),
('Cassandra'),
('Elasticsearch'),

-- Cloud Platforms
('AWS'),
('Azure'),
('Google Cloud'),
('Vercel'),
('Netlify'),
('DigitalOcean'),
('Linode'),
('Heroku'),
('Railway'),
('Render'),
('Fly.io'),
('Cloudflare'),
('Alibaba Cloud'),
('IBM Cloud'),
('Oracle Cloud'),
('Tencent Cloud'),

-- DevOps & Infrastructure
('Docker'),
('Kubernetes'),
('Terraform'),
('Ansible'),
('Jenkins'),
('GitLab CI'),
('GitHub Actions'),
('CircleCI'),
('Travis CI'),
('Bamboo'),
('TeamCity'),
('Helm'),
('Prometheus'),
('Grafana'),
('ELK Stack'),
('Splunk'),
('Datadog'),
('New Relic'),
('PagerDuty'),
('Vault'),
('Consul'),
('Istio'),
('Linkerd'),
('Nginx'),
('Apache'),
('HAProxy'),
('Traefik'),
('AWS Lambda'),
('Azure Functions'),
('Google Cloud Functions'),
('CI/CD'),

-- API Technologies
('GraphQL'),
('REST'),
('gRPC'),
('WebSocket'),
('WebRTC'),
('OpenAPI'),
('Swagger'),
('Postman'),
('Insomnia'),
('Thunder Client'),
('Apollo'),
('Hasura'),
('Prisma'),
('Supabase API'),
('Firebase API'),

-- Version Control
('Git'),
('GitHub'),
('GitLab'),
('Bitbucket'),
('Azure DevOps'),
('Mercurial'),
('SVN'),
('Perforce'),
('Plastic SCM'),

-- Architecture Patterns
('Microservices'),
('Serverless'),
('Event-Driven Architecture'),
('CQRS'),
('Event Sourcing'),
('Domain-Driven Design'),
('Clean Architecture'),
('Hexagonal Architecture'),
('Layered Architecture'),
('Monolithic'),
('Service Mesh'),
('API Gateway'),
('Load Balancing'),
('Circuit Breaker'),
('Bulkhead'),
('Saga Pattern'),

-- Frontend Technologies
('HTML5'),
('CSS3'),
('Sass'),
('Less'),
('Stylus'),
('PostCSS'),
('Tailwind CSS'),
('Material-UI'),
('Ant Design'),
('Chakra UI'),
('Bulma'),
('Foundation'),
('Semantic UI'),
('Vuetify'),
('Quasar'),

-- Build Tools & Bundlers
('Webpack'),
('Vite'),
('Parcel'),
('Rollup'),
('esbuild'),
('SWC'),
('Turbo'),
('Nx'),
('Lerna'),
('Rush'),
('Yarn'),
('pnpm'),
('npm'),

-- Testing Frameworks
('Jest'),
('Cypress'),
('Playwright'),
('Selenium'),
('Puppeteer'),
('Vitest'),
('Testing Library'),
('Mocha'),
('Chai'),
('Sinon'),
('Enzyme'),
('Karma'),
('Jasmine'),
('Ava'),
('Tap'),

-- Mobile Development
('React Native'),
('Flutter'),
('Ionic'),
('Xamarin'),
('Cordova'),
('PhoneGap'),
('Expo'),
('NativeScript'),
('Quasar Mobile'),
('Framework7'),

-- Desktop Development
('Electron'),
('Tauri'),
('Flutter Desktop'),
('Proton Native'),
('Neutralino'),
('Wails'),
('Qt'),
('GTK'),
('WxWidgets'),

-- Data Science & AI
('TensorFlow'),
('PyTorch'),
('Scikit-learn'),
('Pandas'),
('NumPy'),
('Matplotlib'),
('Seaborn'),
('Plotly'),
('D3.js'),
('Apache Spark'),
('Hadoop'),
('Kafka'),
('Airflow'),
('Jupyter'),
('MLflow'),
('Weights & Biases'),
('Comet'),
('Neptune'),

-- Security
('OWASP'),
('JWT'),
('OAuth'),
('SAML'),
('LDAP'),
('SSL/TLS'),
('HTTPS'),
('CORS'),
('CSRF'),
('XSS'),
('SQL Injection'),
('Penetration Testing'),
('Vulnerability Assessment'),
('Security Auditing'),
('Cryptography'),

-- Soft Skills
('Team Leadership'),
('Project Management'),
('Agile'),
('Scrum'),
('Kanban'),
('DevOps Culture'),
('Code Review'),
('Technical Writing'),
('Public Speaking'),
('Mentoring'),
('Cross-functional Collaboration'),
('Problem Solving'),
('Critical Thinking'),
('Communication'),
('Time Management');

-- Insert skills_categories relationships
-- Programming Languages (Category ID: 1)
INSERT INTO skills_categories (skill_id, category_id) 
SELECT s.id, 1 FROM skills s WHERE s.name IN (
    'JavaScript', 'Python', 'Go', 'TypeScript', 'Java', 'C++', 'C#', 'Rust', 
    'Swift', 'Kotlin', 'PHP', 'Ruby', 'Scala', 'R', 'MATLAB', 'Dart', 'Lua', 
    'Haskell', 'Clojure', 'Erlang', 'Elixir', 'F#', 'Assembly', 'C', 'Shell', 
    'PowerShell', 'Bash'
);

-- Frontend Frameworks (Category ID: 2)
INSERT INTO skills_categories (skill_id, category_id) 
SELECT s.id, 2 FROM skills s WHERE s.name IN (
    'React', 'Vue.js', 'Angular', 'Next.js', 'Nuxt.js', 'Svelte', 'SvelteKit', 
    'Gatsby', 'Astro', 'Remix', 'SolidJS', 'Alpine.js', 'jQuery', 'Lit', 'Preact'
);

-- Backend Frameworks (Category ID: 3)
INSERT INTO skills_categories (skill_id, category_id) 
SELECT s.id, 3 FROM skills s WHERE s.name IN (
    'Node.js', 'Express.js', 'FastAPI', 'Django', 'Flask', 'Spring Boot', 
    'Laravel', 'Symfony', 'Rails', 'ASP.NET', 'Deno', 'Bun', 'NestJS', 'Koa', 
    'Hapi', 'Sails', 'Meteor', 'Feathers', 'AdonisJS', 'Fastify', 'Gin', 'Echo', 'Fiber'
);

-- Databases (Category ID: 4)
INSERT INTO skills_categories (skill_id, category_id) 
SELECT s.id, 4 FROM skills s WHERE s.name IN (
    'PostgreSQL', 'MongoDB', 'Redis', 'MySQL', 'SQLite', 'Oracle', 'SQL Server', 
    'MariaDB', 'Neo4j', 'CouchDB', 'RethinkDB', 'InfluxDB', 'TimescaleDB', 
    'ClickHouse', 'Snowflake', 'BigQuery', 'Redshift', 'Firebase', 'Supabase', 
    'PlanetScale', 'CockroachDB', 'FaunaDB', 'ArangoDB', 'OrientDB', 'Amazon RDS', 
    'Azure SQL', 'DynamoDB', 'Cassandra', 'Elasticsearch'
);

-- Cloud Platforms (Category ID: 5)
INSERT INTO skills_categories (skill_id, category_id) 
SELECT s.id, 5 FROM skills s WHERE s.name IN (
    'AWS', 'Azure', 'Google Cloud', 'Vercel', 'Netlify', 'DigitalOcean', 'Linode', 
    'Heroku', 'Railway', 'Render', 'Fly.io', 'Cloudflare', 'Alibaba Cloud', 
    'IBM Cloud', 'Oracle Cloud', 'Tencent Cloud'
);

-- DevOps (Category ID: 6)
INSERT INTO skills_categories (skill_id, category_id) 
SELECT s.id, 6 FROM skills s WHERE s.name IN (
    'Docker', 'Kubernetes', 'Terraform', 'Ansible', 'Jenkins', 'GitLab CI', 
    'GitHub Actions', 'CircleCI', 'Travis CI', 'Bamboo', 'TeamCity', 'Helm', 
    'Prometheus', 'Grafana', 'ELK Stack', 'Splunk', 'Datadog', 'New Relic', 
    'PagerDuty', 'Vault', 'Consul', 'Istio', 'Linkerd', 'Nginx', 'Apache', 
    'HAProxy', 'Traefik', 'AWS Lambda', 'Azure Functions', 'Google Cloud Functions', 'CI/CD'
);

-- API (Category ID: 7)
INSERT INTO skills_categories (skill_id, category_id) 
SELECT s.id, 7 FROM skills s WHERE s.name IN (
    'GraphQL', 'REST', 'gRPC', 'WebSocket', 'WebRTC', 'OpenAPI', 'Swagger', 
    'Postman', 'Insomnia', 'Thunder Client', 'Apollo', 'Hasura', 'Prisma', 
    'Supabase API', 'Firebase API'
);

-- Version Control (Category ID: 8)
INSERT INTO skills_categories (skill_id, category_id) 
SELECT s.id, 8 FROM skills s WHERE s.name IN (
    'Git', 'GitHub', 'GitLab', 'Bitbucket', 'Azure DevOps', 'Mercurial', 'SVN', 
    'Perforce', 'Plastic SCM'
);

-- Architecture (Category ID: 9)
INSERT INTO skills_categories (skill_id, category_id) 
SELECT s.id, 9 FROM skills s WHERE s.name IN (
    'Microservices', 'Serverless', 'Event-Driven Architecture', 'CQRS', 
    'Event Sourcing', 'Domain-Driven Design', 'Clean Architecture', 
    'Hexagonal Architecture', 'Layered Architecture', 'Monolithic', 'Service Mesh', 
    'API Gateway', 'Load Balancing', 'Circuit Breaker', 'Bulkhead', 'Saga Pattern'
);

-- Frontend Technology (Category ID: 10)
INSERT INTO skills_categories (skill_id, category_id) 
SELECT s.id, 10 FROM skills s WHERE s.name IN (
    'HTML5', 'CSS3', 'Sass', 'Less', 'Stylus', 'PostCSS', 'Tailwind CSS', 
    'Material-UI', 'Ant Design', 'Chakra UI', 'Bulma', 'Foundation', 'Semantic UI', 
    'Vuetify', 'Quasar'
);

-- Build Tool (Category ID: 11)
INSERT INTO skills_categories (skill_id, category_id) 
SELECT s.id, 11 FROM skills s WHERE s.name IN (
    'Webpack', 'Vite', 'Parcel', 'Rollup', 'esbuild', 'SWC', 'Turbo', 'Nx', 
    'Lerna', 'Rush', 'Yarn', 'pnpm', 'npm'
);

-- Testing Framework (Category ID: 12)
INSERT INTO skills_categories (skill_id, category_id) 
SELECT s.id, 12 FROM skills s WHERE s.name IN (
    'Jest', 'Cypress', 'Playwright', 'Selenium', 'Puppeteer', 'Vitest', 
    'Testing Library', 'Mocha', 'Chai', 'Sinon', 'Enzyme', 'Karma', 'Jasmine', 
    'Ava', 'Tap'
);

-- Mobile Development (Category ID: 13)
INSERT INTO skills_categories (skill_id, category_id) 
SELECT s.id, 13 FROM skills s WHERE s.name IN (
    'React Native', 'Flutter', 'Ionic', 'Xamarin', 'Cordova', 'PhoneGap', 'Expo', 
    'NativeScript', 'Quasar Mobile', 'Framework7'
);

-- Desktop Development (Category ID: 14)
INSERT INTO skills_categories (skill_id, category_id) 
SELECT s.id, 14 FROM skills s WHERE s.name IN (
    'Electron', 'Tauri', 'Flutter Desktop', 'Proton Native', 'Neutralino', 'Wails', 
    'Qt', 'GTK', 'WxWidgets'
);

-- Data Science & AI (Category ID: 15)
INSERT INTO skills_categories (skill_id, category_id) 
SELECT s.id, 15 FROM skills s WHERE s.name IN (
    'TensorFlow', 'PyTorch', 'Scikit-learn', 'Pandas', 'NumPy', 'Matplotlib', 
    'Seaborn', 'Plotly', 'D3.js', 'Apache Spark', 'Hadoop', 'Kafka', 'Airflow', 
    'Jupyter', 'MLflow', 'Weights & Biases', 'Comet', 'Neptune'
);

-- Security (Category ID: 16)
INSERT INTO skills_categories (skill_id, category_id) 
SELECT s.id, 16 FROM skills s WHERE s.name IN (
    'OWASP', 'JWT', 'OAuth', 'SAML', 'LDAP', 'SSL/TLS', 'HTTPS', 'CORS', 'CSRF', 
    'XSS', 'SQL Injection', 'Penetration Testing', 'Vulnerability Assessment', 
    'Security Auditing', 'Cryptography'
);

-- Soft Skills (Category ID: 17)
INSERT INTO skills_categories (skill_id, category_id) 
SELECT s.id, 17 FROM skills s WHERE s.name IN (
    'Team Leadership', 'Project Management', 'Agile', 'Scrum', 'Kanban', 
    'DevOps Culture', 'Code Review', 'Technical Writing', 'Public Speaking', 
    'Mentoring', 'Cross-functional Collaboration', 'Problem Solving', 
    'Critical Thinking', 'Communication', 'Time Management'
);
