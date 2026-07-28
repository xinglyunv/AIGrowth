#!/bin/bash
set -e

DB_URL="postgres://postgres:postgres@localhost:5432/ai_growth_engine?sslmode=disable"

echo "Authenticating test user..."
API_RESP=$(curl -s http://localhost:8080/api/v1/auth/login \
  -d '{"email":"test@example.com","password":"password123"}')

USER_ID=$(echo "$API_RESP" | python3 -c "
import sys, json
try:
    print(json.load(sys.stdin)['data']['user']['id'])
except Exception as e:
    print(f'Error: {e}', file=sys.stderr)
    sys.exit(1)
")

if [ -z "$USER_ID" ]; then
  echo "Failed to get user ID. API response:"
  echo "$API_RESP"
  exit 1
fi

echo "User ID: $USER_ID"

# Pre-generate deterministic UUIDs for clean references
PROJECT1_ID=$(python3 -c "import uuid; print(uuid.uuid4())")
PROJECT2_ID=$(python3 -c "import uuid; print(uuid.uuid4())")
TASK1_ID=$(python3 -c "import uuid; print(uuid.uuid4())")
TASK2_ID=$(python3 -c "import uuid; print(uuid.uuid4())")
ANS1_1=$(python3 -c "import uuid; print(uuid.uuid4())")
ANS1_2=$(python3 -c "import uuid; print(uuid.uuid4())")
ANS1_3=$(python3 -c "import uuid; print(uuid.uuid4())")
ANS1_4=$(python3 -c "import uuid; print(uuid.uuid4())")
ANS1_5=$(python3 -c "import uuid; print(uuid.uuid4())")
ANS2_1=$(python3 -c "import uuid; print(uuid.uuid4())")
ANS2_2=$(python3 -c "import uuid; print(uuid.uuid4())")
ANS2_3=$(python3 -c "import uuid; print(uuid.uuid4())")
ANS2_4=$(python3 -c "import uuid; print(uuid.uuid4())")
ANS2_5=$(python3 -c "import uuid; print(uuid.uuid4())")

echo "Seeding demo brand data..."

psql "$DB_URL" << EOSQL
-- Brand 1: 智云科技
INSERT INTO brand_projects (id, user_id, name, website, industry, description, keywords, status, created_at, updated_at)
SELECT '$PROJECT1_ID'::uuid, '$USER_ID', '智云科技', 'https://zhiyun-tech.com', 'technology',
       '企业级AI平台和云服务提供商', ARRAY['AI', '云计算', '大数据', '企业服务'], 'active', NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM brand_projects WHERE name = '智云科技' AND user_id = '$USER_ID');

-- Task 1
INSERT INTO ai_tasks (id, project_id, user_id, model, status, questions_count, completed_count, started_at, finished_at, created_at, updated_at)
SELECT '$TASK1_ID'::uuid, '$PROJECT1_ID'::uuid, '$USER_ID', 'gpt-4', 'completed', 5, 5,
       NOW() - interval '2 hours', NOW() - interval '1 hour', NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM ai_tasks WHERE id = '$TASK1_ID'::uuid);

-- Answers 1
INSERT INTO ai_answers (id, task_id, question, answer, model, brand_mentioned, sentiment, rank_position, analysis, created_at)
SELECT '$ANS1_1'::uuid, '$TASK1_ID'::uuid,
       '在编程领域，智云科技的知名度和影响力如何？',
       '智云科技在编程领域拥有较高的品牌认知度，其AI开发平台在开发者社区中获得广泛好评。',
       'gpt-4', true, 'positive', 3, '{"keywords": ["智云科技", "AI开发平台", "开发者社区"], "confidence": 0.92}'::jsonb,
       NOW() - interval '90 minutes';

INSERT INTO ai_answers (id, task_id, question, answer, model, brand_mentioned, sentiment, rank_position, analysis, created_at)
SELECT '$ANS1_2'::uuid, '$TASK1_ID'::uuid,
       '开发者社区中，智云科技的技术栈和工具生态如何？',
       '智云科技提供了一系列开发工具和SDK覆盖主流编程语言，生态系统较为完善。',
       'gpt-4', true, 'positive', 4, '{"keywords": ["智云科技", "开发工具", "SDK"], "confidence": 0.88}'::jsonb,
       NOW() - interval '80 minutes';

INSERT INTO ai_answers (id, task_id, question, answer, model, brand_mentioned, sentiment, rank_position, analysis, created_at)
SELECT '$ANS1_3'::uuid, '$TASK1_ID'::uuid,
       '与传统技术方案相比，使用智云科技有哪些优势和劣势？',
       '智云科技的方案在易用性和集成度上有优势，但在某些特定场景下与传统方案相比成熟度略有不足。',
       'gpt-4', true, 'neutral', 5, '{"keywords": ["智云科技", "易用性", "集成度", "成熟度"], "confidence": 0.85}'::jsonb,
       NOW() - interval '70 minutes';

INSERT INTO ai_answers (id, task_id, question, answer, model, brand_mentioned, sentiment, rank_position, analysis, created_at)
SELECT '$ANS1_4'::uuid, '$TASK1_ID'::uuid,
       '智云科技在人工智能和机器学习领域的应用现状如何？',
       '该领域的讨论主要集中在国外主流平台如TensorFlow和PyTorch上。',
       'gpt-4', false, 'neutral', NULL, '{"keywords": ["TensorFlow", "PyTorch"], "confidence": 0.91}'::jsonb,
       NOW() - interval '60 minutes';

INSERT INTO ai_answers (id, task_id, question, answer, model, brand_mentioned, sentiment, rank_position, analysis, created_at)
SELECT '$ANS1_5'::uuid, '$TASK1_ID'::uuid,
       '未来五年，智云科技在技术领域的发展趋势是什么？',
       '智云科技在前沿技术布局上表现积极，但其长期发展仍需要更多的市场验证。',
       'gpt-4', true, 'neutral', 6, '{"keywords": ["智云科技", "前沿技术", "市场验证"], "confidence": 0.87}'::jsonb,
       NOW() - interval '50 minutes';

-- Brand 2: 星辰电商
INSERT INTO brand_projects (id, user_id, name, website, industry, description, keywords, status, created_at, updated_at)
SELECT '$PROJECT2_ID'::uuid, '$USER_ID', '星辰电商', 'https://xingchen-shop.com', 'e-commerce',
       '新一代社交电商平台', ARRAY['电商', '直播带货', '社交电商', '新零售'], 'active', NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM brand_projects WHERE name = '星辰电商' AND user_id = '$USER_ID');

-- Task 2
INSERT INTO ai_tasks (id, project_id, user_id, model, status, questions_count, completed_count, started_at, finished_at, created_at, updated_at)
SELECT '$TASK2_ID'::uuid, '$PROJECT2_ID'::uuid, '$USER_ID', 'gpt-4', 'completed', 5, 5,
       NOW() - interval '3 hours', NOW() - interval '2 hours', NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM ai_tasks WHERE id = '$TASK2_ID'::uuid);

-- Answers 2
INSERT INTO ai_answers (id, task_id, question, answer, model, brand_mentioned, sentiment, rank_position, analysis, created_at)
SELECT '$ANS2_1'::uuid, '$TASK2_ID'::uuid,
       '在电商购物场景中，星辰电商的用户体验如何？',
       '星辰电商的移动端体验设计流畅，但在大促期间页面加载速度有待提升。',
       'gpt-4', true, 'positive', 2, '{"keywords": ["星辰电商", "移动端", "用户体验"], "confidence": 0.90}'::jsonb,
       NOW() - interval '150 minutes';

INSERT INTO ai_answers (id, task_id, question, answer, model, brand_mentioned, sentiment, rank_position, analysis, created_at)
SELECT '$ANS2_2'::uuid, '$TASK2_ID'::uuid,
       '消费者对星辰电商的物流配送服务评价如何？',
       '物流配送覆盖范围广泛，但配送时效与头部平台相比仍有差距。',
       'gpt-4', true, 'neutral', 4, '{"keywords": ["星辰电商", "物流配送", "配送时效"], "confidence": 0.86}'::jsonb,
       NOW() - interval '140 minutes';

INSERT INTO ai_answers (id, task_id, question, answer, model, brand_mentioned, sentiment, rank_position, analysis, created_at)
SELECT '$ANS2_3'::uuid, '$TASK2_ID'::uuid,
       '星辰电商平台上商品的价格竞争力如何？',
       '星辰电商的价格策略较为灵活，但在家电等品类上与主要竞争对手相比优势不明显。',
       'gpt-4', true, 'neutral', 5, '{"keywords": ["星辰电商", "价格策略", "竞争优势"], "confidence": 0.84}'::jsonb,
       NOW() - interval '130 minutes';

INSERT INTO ai_answers (id, task_id, question, answer, model, brand_mentioned, sentiment, rank_position, analysis, created_at)
SELECT '$ANS2_4'::uuid, '$TASK2_ID'::uuid,
       '与主要竞争对手相比，星辰电商的售后服务如何？',
       '售后服务质量参差不齐，部分用户反映退换货流程复杂。',
       'gpt-4', true, 'negative', 7, '{"keywords": ["星辰电商", "售后服务", "退换货"], "confidence": 0.93}'::jsonb,
       NOW() - interval '120 minutes';

INSERT INTO ai_answers (id, task_id, question, answer, model, brand_mentioned, sentiment, rank_position, analysis, created_at)
SELECT '$ANS2_5'::uuid, '$TASK2_ID'::uuid,
       '星辰电商在移动端购物的体验优化做得如何？',
       '移动端界面设计简洁美观，但个性化推荐准确率有待提升。',
       'gpt-4', true, 'positive', 3, '{"keywords": ["星辰电商", "移动端", "界面设计", "个性化推荐"], "confidence": 0.89}'::jsonb,
       NOW() - interval '110 minutes';

EOSQL

echo ""
echo "Seeding complete!"
echo ""
echo "Summary:"
psql "$DB_URL" -c "
SELECT
  bp.name AS brand,
  bp.industry,
  at2.status AS task_status,
  at2.questions_count,
  at2.completed_count,
  (SELECT count(*) FROM ai_answers aa WHERE aa.task_id = at2.id) AS answers_count
FROM brand_projects bp
JOIN ai_tasks at2 ON at2.project_id = bp.id
WHERE bp.user_id = '$USER_ID' AND bp.name IN ('智云科技', '星辰电商')
ORDER BY bp.name;
"
