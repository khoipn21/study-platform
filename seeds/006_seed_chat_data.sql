-- Development seed data for chat tables

-- Seed chat sessions
INSERT INTO chat_sessions (id, user_id, title, created_at, updated_at, message_count, is_active) VALUES

-- Alice's chat sessions
('chat001', '55555555-5555-5555-5555-555555555555', 'Help with JavaScript Closures', '2024-01-20 10:30:00', '2024-01-20 11:15:00', 6, false),
('chat002', '55555555-5555-5555-5555-555555555555', 'React Component Questions', '2024-02-05 14:20:00', '2024-02-05 14:45:00', 4, false),
('chat003', '55555555-5555-5555-5555-555555555555', 'Current Session', '2024-02-25 09:00:00', '2024-02-25 09:00:00', 1, true),

-- Bob's chat sessions
('chat004', '66666666-6666-6666-6666-666666666666', 'Database Design Help', '2024-02-10 16:40:00', '2024-02-10 17:25:00', 8, false),
('chat005', '66666666-6666-6666-6666-666666666666', 'SQL Query Optimization', '2024-02-18 11:15:00', '2024-02-18 11:50:00', 5, false),

-- Carol's chat sessions
('chat006', '77777777-7777-7777-7777-777777777777', 'CSS Grid Layout Issues', '2024-02-12 13:30:00', '2024-02-12 14:10:00', 7, false),

-- David's chat sessions
('chat007', '88888888-8888-8888-8888-888888888888', 'Career Advice for Developers', '2024-02-15 10:45:00', '2024-02-15 11:20:00', 6, false),
('chat008', '88888888-8888-8888-8888-888888888888', 'JavaScript Debugging Tips', '2024-02-22 15:30:00', '2024-02-22 16:00:00', 4, false);

-- Seed chat history
INSERT INTO chat_history (id, session_id, user_id, message, sender, response, created_at) VALUES

-- Alice's JavaScript Closures conversation
('msg001', 'chat001', '55555555-5555-5555-5555-555555555555', 'I''m having trouble understanding JavaScript closures. Can you explain them simply?', 'user', 'A closure in JavaScript is when a function has access to variables from its outer (enclosing) scope even after the outer function has finished executing. Think of it as the function "remembering" variables from where it was created.', '2024-01-20 10:30:00'),
('msg002', 'chat001', '55555555-5555-5555-5555-555555555555', 'Can you give me a simple example?', 'user', 'Sure! Here''s a basic example:\n\n```javascript\nfunction outerFunction(x) {\n  return function innerFunction(y) {\n    return x + y; // x is remembered from outerFunction\n  };\n}\n\nconst addFive = outerFunction(5);\nconsole.log(addFive(3)); // 8\n```\n\nThe inner function "closes over" the variable x.', '2024-01-20 10:35:00'),
('msg003', 'chat001', '55555555-5555-5555-5555-555555555555', 'Why are closures useful?', 'user', 'Closures are useful for:\n1. Creating private variables\n2. Function factories\n3. Event handlers that need to remember state\n4. Module patterns\n5. Callbacks that need access to their creation context\n\nThey''re fundamental to how JavaScript works!', '2024-01-20 10:40:00'),

-- Bob's Database Design conversation
('msg004', 'chat004', '66666666-6666-6666-6666-666666666666', 'I''m designing a database for an e-commerce site. What tables should I start with?', 'user', 'For an e-commerce database, you''ll typically need these core tables:\n\n1. **Users** - customer information\n2. **Products** - product details, prices, inventory\n3. **Categories** - product categorization\n4. **Orders** - order information\n5. **Order_Items** - products within each order\n6. **Shopping_Cart** - temporary cart items\n7. **Addresses** - shipping and billing addresses\n\nStart with these basics and add more as needed.', '2024-02-10 16:40:00'),
('msg005', 'chat004', '66666666-6666-6666-6666-666666666666', 'How should I handle the relationship between products and categories?', 'user', 'You have a few options depending on your needs:\n\n**Option 1: Simple many-to-one**\n- Add category_id to products table\n- Each product belongs to one category\n\n**Option 2: Many-to-many**\n- Create product_categories junction table\n- Products can belong to multiple categories\n\n**Option 3: Hierarchical categories**\n- Add parent_id to categories table\n- Support nested categories\n\nFor most e-commerce sites, Option 2 or 3 works best.', '2024-02-10 16:50:00'),

-- Carol's CSS Grid conversation
('msg006', 'chat006', '77777777-7777-7777-7777-777777777777', 'I''m trying to create a responsive grid layout but it''s not working on mobile. Any tips?', 'user', 'Common issues with responsive CSS Grid:\n\n1. **Use fr units instead of fixed widths**\n```css\ngrid-template-columns: repeat(auto-fit, minmax(250px, 1fr));\n```\n\n2. **Add proper breakpoints**\n```css\n@media (max-width: 768px) {\n  grid-template-columns: 1fr;\n}\n```\n\n3. **Don''t forget the viewport meta tag**\n```html\n<meta name="viewport" content="width=device-width, initial-scale=1">\n```\n\nWhat specific issues are you seeing?', '2024-02-12 13:30:00'),
('msg007', 'chat006', '77777777-7777-7777-7777-777777777777', 'The grid items are overflowing on small screens', 'user', 'Grid overflow usually happens because of:\n\n1. **Fixed widths** - use minmax() instead\n2. **Too many columns** - use auto-fit or media queries\n3. **Content doesn''t wrap** - add word-break: break-word\n\n```css\n.grid-container {\n  display: grid;\n  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));\n  gap: 1rem;\n}\n\n.grid-item {\n  word-break: break-word;\n  overflow-wrap: break-word;\n}\n```\n\nThis should fix the overflow issue!', '2024-02-12 13:40:00'),

-- David's Career Advice conversation  
('msg008', 'chat007', '88888888-8888-8888-8888-888888888888', 'I''ve been learning for 8 months. How do I know when I''m ready to apply for jobs?', 'user', 'Great question! You''re likely ready to start applying when you can:\n\n✅ **Technical Skills**\n- Build a complete web application from scratch\n- Use Git/GitHub confidently\n- Debug your own code effectively\n- Understand core programming concepts\n\n✅ **Portfolio**\n- 3-4 solid projects showcasing different skills\n- Clean, documented code\n- Deployed applications (not just localhost)\n\n✅ **Soft Skills**\n- Explain your code to others\n- Learn new technologies independently\n- Handle feedback and iteration\n\nDon''t wait until you feel "perfect" - apply when you''re 70-80% confident!', '2024-02-15 10:45:00'),
('msg009', 'chat007', '88888888-8888-8888-8888-888888888888', 'What should my portfolio projects include?', 'user', 'Your portfolio should demonstrate variety and depth:\n\n**Project Ideas:**\n1. **CRUD Application** - Todo app, blog, inventory system\n2. **API Integration** - Weather app, movie database, social media clone\n3. **Full-Stack Project** - E-commerce, social platform, project management\n4. **Your Choice** - Something you''re passionate about!\n\n**Each project should show:**\n- Clean, readable code\n- Responsive design\n- Error handling\n- Some complexity (not just tutorials)\n- README with setup instructions\n\n**Bonus points:**\n- Testing\n- CI/CD\n- Performance optimization\n- Accessibility features', '2024-02-15 11:00:00');