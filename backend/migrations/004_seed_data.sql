-- ============================================================
-- FILE: migrations/004_seed_data.sql
-- ============================================================

-- CITIES
INSERT INTO cities (id, name, country, region, lat, lng, population, timezone) VALUES
  ('c0000001-0000-0000-0000-000000000001', 'San Francisco', 'USA',         'Americas',    37.7749, -122.4194, 883305,   'America/Los_Angeles'),
  ('c0000001-0000-0000-0000-000000000002', 'New York',      'USA',         'Americas',    40.7128,  -74.0060, 8336817,  'America/New_York'),
  ('c0000001-0000-0000-0000-000000000003', 'Austin',        'USA',         'Americas',    30.2672,  -97.7431, 961855,   'America/Chicago'),
  ('c0000001-0000-0000-0000-000000000004', 'Seattle',       'USA',         'Americas',    47.6062, -122.3321, 737255,   'America/Los_Angeles'),
  ('c0000001-0000-0000-0000-000000000005', 'Toronto',       'Canada',      'Americas',    43.6532,  -79.3832, 2731571,  'America/Toronto'),
  ('c0000001-0000-0000-0000-000000000006', 'London',        'UK',          'Europe',      51.5074,   -0.1278, 8982000,  'Europe/London'),
  ('c0000001-0000-0000-0000-000000000007', 'Berlin',        'Germany',     'Europe',      52.5200,   13.4050, 3769495,  'Europe/Berlin'),
  ('c0000001-0000-0000-0000-000000000008', 'Stockholm',     'Sweden',      'Europe',      59.3293,   18.0686, 975551,   'Europe/Stockholm'),
  ('c0000001-0000-0000-0000-000000000009', 'Amsterdam',     'Netherlands', 'Europe',      52.3676,    4.9041, 872680,   'Europe/Amsterdam'),
  ('c0000001-0000-0000-0000-000000000010', 'Zurich',        'Switzerland', 'Europe',      47.3769,    8.5417, 415367,   'Europe/Zurich'),
  ('c0000001-0000-0000-0000-000000000011', 'Tel Aviv',      'Israel',      'Middle East', 32.0853,   34.7818, 460613,   'Asia/Jerusalem'),
  ('c0000001-0000-0000-0000-000000000012', 'Dubai',         'UAE',         'Middle East', 25.2048,   55.2708, 3331420,  'Asia/Dubai'),
  ('c0000001-0000-0000-0000-000000000013', 'Singapore',     'Singapore',   'Asia',         1.3521,  103.8198, 5850342,  'Asia/Singapore'),
  ('c0000001-0000-0000-0000-000000000014', 'Tokyo',         'Japan',       'Asia',        35.6762,  139.6503, 13960000, 'Asia/Tokyo'),
  ('c0000001-0000-0000-0000-000000000015', 'Seoul',         'South Korea', 'Asia',        37.5665,  126.9780, 9776000,  'Asia/Seoul'),
  ('c0000001-0000-0000-0000-000000000016', 'Bangalore',     'India',       'Asia',        12.9716,   77.5946, 12765000, 'Asia/Kolkata'),
  ('c0000001-0000-0000-0000-000000000017', 'Sydney',        'Australia',   'Oceania',    -33.8688,  151.2093, 5312000,  'Australia/Sydney'),
  ('c0000001-0000-0000-0000-000000000018', 'Lagos',         'Nigeria',     'Africa',       6.5244,    3.3792, 14800000, 'Africa/Lagos'),
  ('c0000001-0000-0000-0000-000000000019', 'Nairobi',       'Kenya',       'Africa',      -1.2921,   36.8219, 4397073,  'Africa/Nairobi'),
  ('c0000001-0000-0000-0000-000000000020', 'Tashkent',      'Uzbekistan',  'Central Asia',41.2995,   69.2401, 2571668,  'Asia/Tashkent');

-- CITY SCORES
INSERT INTO city_scores (city_id, score, job_growth_pct, remote_score, ai_investment, talent_demand, cost_of_living) VALUES
  ('c0000001-0000-0000-0000-000000000001', 95, 12.0, 85, 98, 96, 22),
  ('c0000001-0000-0000-0000-000000000002', 91,  9.0, 78, 90, 93, 28),
  ('c0000001-0000-0000-0000-000000000003', 89, 14.0, 83, 85, 88, 42),
  ('c0000001-0000-0000-0000-000000000004', 90, 11.0, 82, 92, 91, 32),
  ('c0000001-0000-0000-0000-000000000005', 86, 13.0, 82, 80, 85, 48),
  ('c0000001-0000-0000-0000-000000000006', 88, 10.0, 80, 84, 87, 35),
  ('c0000001-0000-0000-0000-000000000007', 85, 11.0, 78, 78, 83, 52),
  ('c0000001-0000-0000-0000-000000000008', 89, 14.0, 85, 82, 88, 45),
  ('c0000001-0000-0000-0000-000000000009', 87, 12.0, 83, 80, 86, 44),
  ('c0000001-0000-0000-0000-000000000010', 90, 10.0, 82, 86, 90, 30),
  ('c0000001-0000-0000-0000-000000000011', 91, 17.0, 80, 88, 90, 40),
  ('c0000001-0000-0000-0000-000000000012', 82, 18.0, 65, 75, 80, 38),
  ('c0000001-0000-0000-0000-000000000013', 92, 15.0, 75, 89, 91, 36),
  ('c0000001-0000-0000-0000-000000000014', 87,  9.0, 60, 83, 86, 48),
  ('c0000001-0000-0000-0000-000000000015', 84, 11.0, 65, 80, 83, 50),
  ('c0000001-0000-0000-0000-000000000016', 83, 20.0, 70, 76, 84, 72),
  ('c0000001-0000-0000-0000-000000000017', 80,  8.0, 79, 72, 78, 44),
  ('c0000001-0000-0000-0000-000000000018', 68, 25.0, 55, 58, 65, 80),
  ('c0000001-0000-0000-0000-000000000019', 70, 22.0, 58, 60, 67, 82),
  ('c0000001-0000-0000-0000-000000000020', 64, 28.0, 52, 52, 62, 88);

-- PROFESSIONS (current)
INSERT INTO professions (slug, title, category, ai_risk_score, avg_salary_usd, description, is_future_job, demand_index, growth_pct) VALUES
  ('software-engineer',     'Software Engineer',     'Tech',        45, 120000, 'Design and build software systems and applications',                 false, 88,  8.0),
  ('data-scientist',        'Data Scientist',        'Tech',        38, 130000, 'Analyze complex data to extract actionable business insights',       false, 90, 10.0),
  ('cybersecurity-analyst', 'Cybersecurity Analyst', 'Tech',        28, 115000, 'Protect organizations from digital threats and breaches',            false, 92, 13.0),
  ('devops-engineer',       'DevOps Engineer',       'Tech',        42, 125000, 'Bridge development and operations to accelerate delivery',           false, 87,  9.0),
  ('ux-designer',           'UX Designer',           'Creative',    40,  95000, 'Design intuitive and delightful user experiences',                   false, 82,  7.0),
  ('graphic-designer',      'Graphic Designer',      'Creative',    75,  60000, 'Create visual content for brands and communications',                false, 60, -2.0),
  ('content-writer',        'Content Writer',        'Creative',    80,  55000, 'Produce written content for digital and print media',                false, 55, -5.0),
  ('financial-analyst',     'Financial Analyst',     'Finance',     85,  90000, 'Analyze financial data and advise on investment decisions',           false, 65, -3.0),
  ('accountant',            'Accountant',            'Finance',     94,  70000, 'Manage financial records and ensure regulatory compliance',          false, 55, -8.0),
  ('marketing-manager',     'Marketing Manager',     'Business',    60,  85000, 'Develop and execute marketing strategies for growth',                false, 75,  4.0),
  ('project-manager',       'Project Manager',       'Business',    50,  95000, 'Lead cross-functional teams to deliver projects on time',            false, 80,  5.0),
  ('lawyer',                'Lawyer',                'Legal',       55, 140000, 'Provide legal counsel and represent clients in disputes',            false, 72,  2.0),
  ('doctor',                'Doctor',                'Healthcare',  25, 200000, 'Diagnose and treat illness using evidence-based medicine',           false, 95,  6.0),
  ('nurse',                 'Nurse',                 'Healthcare',  20,  75000, 'Provide direct patient care and clinical support',                   false, 96,  9.0),
  ('teacher',               'Teacher',               'Education',   30,  55000, 'Educate and inspire students across subjects and age groups',        false, 85,  4.0),
  ('truck-driver',          'Truck Driver',          'Transport',   90,  50000, 'Transport goods across long distances by road',                      false, 40,-12.0);

-- FUTURE PROFESSIONS
INSERT INTO professions (slug, title, category, ai_risk_score, avg_salary_usd, description, is_future_job, demand_index, growth_pct) VALUES
  ('ai-ethics-officer',         'AI Ethics Officer',                'Tech',        5,  170000, 'Govern AI systems for fairness, transparency and legal compliance',           true, 95, 340.0),
  ('prompt-engineer',           'Prompt Engineer',                  'Tech',       10,  130000, 'Design and optimize AI systems and workflows across industries',               true, 92, 290.0),
  ('quantum-developer',         'Quantum Developer',                'Tech',        8,  190000, 'Build next-gen applications on quantum computing platforms',                  true, 85, 310.0),
  ('climate-tech-engineer',     'Climate Tech Engineer',            'Green',       5,  140000, 'Engineer large-scale climate change mitigation solutions',                    true, 90, 280.0),
  ('human-ai-specialist',       'Human-AI Collaboration Spec.',     'Business',   10,  130000, 'Bridge human workers and intelligent AI systems in the workplace',            true, 89, 260.0),
  ('longevity-scientist',       'Longevity Scientist',              'Life Science',5,  160000, 'Research anti-aging therapies and life extension biotechnologies',            true, 80, 220.0),
  ('biosecurity-analyst',       'Biosecurity Analyst',              'Life Science',8,  140000, 'Protect nations and organizations from emerging biological threats',          true, 80, 200.0),
  ('green-energy-consultant',   'Green Energy Consultant',          'Green',       6,  120000, 'Guide organizations through full renewable energy transition strategies',     true, 88, 240.0),
  ('neural-interface-developer','Neural Interface Developer',       'Tech',        7,  200000, 'Build software and applications for brain-computer interface platforms',      true, 72, 380.0),
  ('digital-wellness-coach',    'Digital Wellness Coach',           'Wellness',   12,   90000, 'Help individuals maintain wellbeing and balance in a hyper-digital world',   true, 76, 175.0),
  ('metaverse-architect',       'Metaverse Architect',              'Creative',   12,  150000, 'Design immersive virtual worlds and spatial computing experiences',           true, 70, 190.0),
  ('space-tourism-specialist',  'Space Tourism Specialist',         'Space',       6,  115000, 'Design and manage commercial space travel experiences for civilians',         true, 65, 450.0);

-- CAREER TRANSITIONS (Software Engineer → future roles)
INSERT INTO career_transitions (from_profession, to_profession, match_score, transition_reason, avg_reskill_months)
SELECT f.id, t.id, match, reason, months FROM (VALUES
  ('software-engineer', 'prompt-engineer',           90, 'Programming mindset is the #1 foundation for prompt engineering', 3),
  ('software-engineer', 'ai-ethics-officer',         85, 'System design knowledge is essential for AI governance work',     6),
  ('software-engineer', 'quantum-developer',         78, 'Algorithms and math background transfers directly',               12),
  ('software-engineer', 'human-ai-specialist',       82, 'Deep tech understanding helps bridge humans and AI tools',        5),
  ('software-engineer', 'neural-interface-developer',74, 'Software architecture is core to BCI platform development',      18),
  ('data-scientist',    'longevity-scientist',        88, 'Data analysis is central to aging and genomics research',        8),
  ('data-scientist',    'climate-tech-engineer',      84, 'Climate modeling demands your exact analytical skill set',       6),
  ('data-scientist',    'biosecurity-analyst',        80, 'Epidemiological modeling is a core biosecurity skill',           10),
  ('data-scientist',    'ai-ethics-officer',          76, 'Model bias and fairness work is highly valued in AI governance', 6),
  ('data-scientist',    'quantum-developer',          70, 'Quantum ML is the fastest-growing frontier in your field',       14),
  ('lawyer',            'ai-ethics-officer',          92, 'Legal expertise is the #1 skill needed in AI governance',        4),
  ('lawyer',            'biosecurity-analyst',        78, 'Policy and regulatory knowledge is crucial for biosecurity',     10),
  ('lawyer',            'human-ai-specialist',        72, 'Negotiation and contract logic transfers to AI mediation',        6),
  ('teacher',           'human-ai-specialist',        86, 'Ability to translate complex concepts is invaluable here',       5),
  ('teacher',           'digital-wellness-coach',     82, 'Coaching and mentoring skills transfer directly',                4),
  ('financial-analyst', 'ai-ethics-officer',          80, 'Risk assessment and compliance skills are directly applicable',  6),
  ('financial-analyst', 'green-energy-consultant',    75, 'ESG financial modeling is in extremely high demand',             8)
) AS v(from_slug, to_slug, match, reason, months)
JOIN professions f ON f.slug = v.from_slug
JOIN professions t ON t.slug = v.to_slug;
