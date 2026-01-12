-- Sample data for testing
INSERT INTO videos (title, description, url, thumbnail, channel_name, channel_avatar, views, duration, uploaded_at) VALUES
('Building a YouTube Clone with React and Tailwind CSS', 'Learn how to build a modern YouTube clone using React and Tailwind CSS. This comprehensive tutorial covers UI components, routing, and API integration.', 'https://example.com/video1.mp4', 'https://via.placeholder.com/320x180/FF0000/FFFFFF?text=React+Tutorial', 'Code Master', 'https://via.placeholder.com/36/4299E1/FFFFFF?text=CM', 125000, '12:34', NOW() - INTERVAL '2 days'),

('Golang Backend Tutorial - Building RESTful APIs', 'Master Go programming by building RESTful APIs from scratch. Learn about routing, middleware, database integration, and best practices.', 'https://example.com/video2.mp4', 'https://via.placeholder.com/320x180/00ADD8/FFFFFF?text=Go+Tutorial', 'Go Developer', 'https://via.placeholder.com/36/48BB78/FFFFFF?text=GD', 85400, '24:15', NOW() - INTERVAL '1 week'),

('PostgreSQL Database Design Best Practices', 'Comprehensive guide to PostgreSQL database design. Learn about normalization, indexing, query optimization, and performance tuning.', 'https://example.com/video3.mp4', 'https://via.placeholder.com/320x180/336791/FFFFFF?text=PostgreSQL', 'Database Pro', 'https://via.placeholder.com/36/9F7AEA/FFFFFF?text=DP', 62300, '18:45', NOW() - INTERVAL '3 days'),

('React Hooks Deep Dive - useState and useEffect', 'Deep dive into React Hooks! Understand useState, useEffect, and custom hooks. Learn when and how to use each hook effectively.', 'https://example.com/video4.mp4', 'https://via.placeholder.com/320x180/61DAFB/000000?text=React+Hooks', 'React Expert', 'https://via.placeholder.com/36/ED8936/FFFFFF?text=RE', 210000, '15:23', NOW() - INTERVAL '5 days'),

('Tailwind CSS Complete Guide for Beginners', 'Everything you need to know about Tailwind CSS. From basic utilities to advanced customization and responsive design patterns.', 'https://example.com/video5.mp4', 'https://via.placeholder.com/320x180/38B2AC/FFFFFF?text=Tailwind+CSS', 'CSS Wizard', 'https://via.placeholder.com/36/F56565/FFFFFF?text=CW', 180000, '32:10', NOW() - INTERVAL '1 week'),

('Full Stack Development Roadmap 2024', 'Complete roadmap for becoming a full stack developer in 2024. Learn about essential technologies, frameworks, and career paths.', 'https://example.com/video6.mp4', 'https://via.placeholder.com/320x180/805AD5/FFFFFF?text=Roadmap+2024', 'Tech Career', 'https://via.placeholder.com/36/38B2AC/FFFFFF?text=TC', 350000, '45:00', NOW() - INTERVAL '2 weeks'),

('Docker and Kubernetes for Beginners', 'Learn containerization with Docker and orchestration with Kubernetes. Perfect for DevOps beginners and developers.', 'https://example.com/video7.mp4', 'https://via.placeholder.com/320x180/2496ED/FFFFFF?text=Docker+K8s', 'DevOps Guru', 'https://via.placeholder.com/36/667EEA/FFFFFF?text=DG', 142000, '28:30', NOW() - INTERVAL '4 days'),

('JavaScript ES6+ Features You Must Know', 'Modern JavaScript features explained! Arrow functions, destructuring, spread operator, promises, async/await, and more.', 'https://example.com/video8.mp4', 'https://via.placeholder.com/320x180/F7DF1E/000000?text=JavaScript+ES6', 'JS Master', 'https://via.placeholder.com/36/F6AD55/FFFFFF?text=JS', 195000, '21:45', NOW() - INTERVAL '6 days'),

('Git and GitHub Workflow Best Practices', 'Master version control with Git and GitHub. Learn about branching strategies, pull requests, code reviews, and CI/CD integration.', 'https://example.com/video9.mp4', 'https://via.placeholder.com/320x180/F05032/FFFFFF?text=Git+GitHub', 'Code Versioning', 'https://via.placeholder.com/36/FC8181/FFFFFF?text=CV', 98500, '16:20', NOW() - INTERVAL '1 week'),

('REST API Design Principles and Best Practices', 'Learn how to design clean, maintainable REST APIs. Covers HTTP methods, status codes, versioning, and documentation.', 'https://example.com/video10.mp4', 'https://via.placeholder.com/320x180/10B981/FFFFFF?text=REST+API', 'API Architect', 'https://via.placeholder.com/36/68D391/FFFFFF?text=AA', 73200, '19:15', NOW() - INTERVAL '3 days');
