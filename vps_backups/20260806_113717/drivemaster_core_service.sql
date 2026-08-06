--
-- PostgreSQL database dump
--

\restrict EHYbGxzaUf53OF9JeJdInm5Prm97WPQG1fZ6qFH2A4pje7ok5XrVqqxEyyFzHc5

-- Dumped from database version 16.14
-- Dumped by pg_dump version 16.14

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

ALTER TABLE ONLY public.sale_items DROP CONSTRAINT fk_sales_items;
ALTER TABLE ONLY public.categories DROP CONSTRAINT fk_categories_parent;
ALTER TABLE ONLY public.articles DROP CONSTRAINT fk_articles_category;
ALTER TABLE ONLY public.article_tags DROP CONSTRAINT fk_article_tags_tag;
DROP INDEX public.idx_tags_slug;
DROP INDEX public.idx_tags_name;
DROP INDEX public.idx_sales_user_id;
DROP INDEX public.idx_sales_status;
DROP INDEX public.idx_sales_payment_id;
DROP INDEX public.idx_sales_package_id;
DROP INDEX public.idx_sales_order_number;
DROP INDEX public.idx_sale_items_sale_id;
DROP INDEX public.idx_related;
DROP INDEX public.idx_pages_slug;
DROP INDEX public.idx_pages_deleted_at;
DROP INDEX public.idx_faqs_deleted_at;
DROP INDEX public.idx_categories_slug;
DROP INDEX public.idx_cars_license_plate;
DROP INDEX public.idx_articles_slug;
DROP INDEX public.idx_articles_deleted_at;
DROP INDEX public.idx_article_tag;
ALTER TABLE ONLY public.tags DROP CONSTRAINT tags_pkey;
ALTER TABLE ONLY public.sales DROP CONSTRAINT sales_pkey;
ALTER TABLE ONLY public.sale_items DROP CONSTRAINT sale_items_pkey;
ALTER TABLE ONLY public.related_articles DROP CONSTRAINT related_articles_pkey;
ALTER TABLE ONLY public.regencies DROP CONSTRAINT regencies_pkey;
ALTER TABLE ONLY public.provinces DROP CONSTRAINT provinces_pkey;
ALTER TABLE ONLY public.pages DROP CONSTRAINT pages_pkey;
ALTER TABLE ONLY public.packages DROP CONSTRAINT packages_pkey;
ALTER TABLE ONLY public.package_benefits DROP CONSTRAINT package_benefits_pkey;
ALTER TABLE ONLY public.monthly_sales DROP CONSTRAINT monthly_sales_pkey;
ALTER TABLE ONLY public.general_settings DROP CONSTRAINT general_settings_pkey;
ALTER TABLE ONLY public.faqs DROP CONSTRAINT faqs_pkey;
ALTER TABLE ONLY public.districts DROP CONSTRAINT districts_pkey;
ALTER TABLE ONLY public.contact_inquiries DROP CONSTRAINT contact_inquiries_pkey;
ALTER TABLE ONLY public.categories DROP CONSTRAINT categories_pkey;
ALTER TABLE ONLY public.cars DROP CONSTRAINT cars_pkey;
ALTER TABLE ONLY public.articles DROP CONSTRAINT articles_pkey;
ALTER TABLE ONLY public.article_tags DROP CONSTRAINT article_tags_pkey;
ALTER TABLE ONLY public.add_ons DROP CONSTRAINT add_ons_pkey;
ALTER TABLE public.regencies ALTER COLUMN id DROP DEFAULT;
ALTER TABLE public.provinces ALTER COLUMN id DROP DEFAULT;
ALTER TABLE public.districts ALTER COLUMN id DROP DEFAULT;
DROP TABLE public.tags;
DROP TABLE public.sales;
DROP TABLE public.sale_items;
DROP TABLE public.related_articles;
DROP SEQUENCE public.regencies_id_seq;
DROP TABLE public.regencies;
DROP SEQUENCE public.provinces_id_seq;
DROP TABLE public.provinces;
DROP TABLE public.pages;
DROP TABLE public.packages;
DROP TABLE public.package_benefits;
DROP TABLE public.monthly_sales;
DROP TABLE public.general_settings;
DROP TABLE public.faqs;
DROP SEQUENCE public.districts_id_seq;
DROP TABLE public.districts;
DROP TABLE public.contact_inquiries;
DROP TABLE public.categories;
DROP TABLE public.cars;
DROP TABLE public.articles;
DROP TABLE public.article_tags;
DROP TABLE public.add_ons;
SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: add_ons; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.add_ons (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    title character varying(255) NOT NULL,
    description text,
    price numeric(10,2) NOT NULL,
    sessions bigint DEFAULT 1,
    status character varying(20) DEFAULT 'active'::character varying NOT NULL,
    image_url character varying(500),
    sort_order bigint DEFAULT 0,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.add_ons OWNER TO postgres;

--
-- Name: article_tags; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.article_tags (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    article_id uuid,
    tag_id uuid
);


ALTER TABLE public.article_tags OWNER TO postgres;

--
-- Name: articles; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.articles (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    title character varying(255) NOT NULL,
    slug character varying(255) NOT NULL,
    author character varying(100),
    content text,
    media jsonb DEFAULT '[]'::jsonb,
    lead_paragraph text,
    excerpt text,
    footer text,
    meta_title character varying(70),
    meta_description character varying(160),
    meta_keywords character varying(500),
    canonical_url character varying(500),
    og_title character varying(95),
    og_description character varying(200),
    og_image character varying(500),
    featured_image character varying(500),
    thumbnail_url character varying(500),
    category_id uuid,
    tags text[],
    author_id uuid NOT NULL,
    status character varying(20) DEFAULT 'draft'::character varying,
    published_at timestamp with time zone,
    scheduled_at timestamp with time zone,
    view_count bigint DEFAULT 0,
    like_count bigint DEFAULT 0,
    share_count bigint DEFAULT 0,
    reading_time bigint DEFAULT 0,
    is_featured boolean DEFAULT false,
    is_spotlight boolean DEFAULT false,
    priority bigint DEFAULT 0,
    highlight boolean DEFAULT false,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    language character varying(10) DEFAULT 'en'::character varying,
    version bigint DEFAULT 1,
    estimated_time bigint,
    video_url character varying(500),
    cta_text character varying(100),
    cta_link character varying(500),
    avg_read_time numeric,
    bounce_rate numeric,
    completion_rate numeric,
    body_blocks jsonb
);


ALTER TABLE public.articles OWNER TO postgres;

--
-- Name: cars; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.cars (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    brand character varying(100) NOT NULL,
    model character varying(100) NOT NULL,
    year bigint NOT NULL,
    license_plate character varying(20) NOT NULL,
    color character varying(50),
    transmission character varying(20) DEFAULT 'manual'::character varying NOT NULL,
    status character varying(20) DEFAULT 'available'::character varying NOT NULL,
    mileage bigint DEFAULT 0,
    image_url character varying(500),
    notes text,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.cars OWNER TO postgres;

--
-- Name: categories; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.categories (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    name character varying(100) NOT NULL,
    slug character varying(100) NOT NULL,
    description character varying(255),
    parent_id uuid,
    "order" bigint DEFAULT 0,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.categories OWNER TO postgres;

--
-- Name: contact_inquiries; Type: TABLE; Schema: public; Owner: admin_drive
--

CREATE TABLE public.contact_inquiries (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    name character varying(255) NOT NULL,
    email character varying(255) NOT NULL,
    subject character varying(255) NOT NULL,
    message text NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    status character varying(50) DEFAULT 'unread'::character varying NOT NULL
);


ALTER TABLE public.contact_inquiries OWNER TO admin_drive;

--
-- Name: districts; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.districts (
    id bigint NOT NULL,
    name character varying(100) NOT NULL,
    regency_id bigint NOT NULL
);


ALTER TABLE public.districts OWNER TO postgres;

--
-- Name: districts_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.districts_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.districts_id_seq OWNER TO postgres;

--
-- Name: districts_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.districts_id_seq OWNED BY public.districts.id;


--
-- Name: faqs; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.faqs (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    question character varying(500) NOT NULL,
    answer text NOT NULL,
    "order" bigint DEFAULT 0,
    category character varying(100) DEFAULT 'general'::character varying,
    is_active boolean DEFAULT true,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);


ALTER TABLE public.faqs OWNER TO postgres;

--
-- Name: general_settings; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.general_settings (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    business_name character varying(255),
    email character varying(255),
    phone character varying(50),
    fax character varying(50),
    whats_app character varying(50),
    address text,
    hours_mon_fri character varying(100),
    hours_sat_sun character varying(100),
    hours_night_shift character varying(100),
    promo_end_date date,
    notify_email boolean DEFAULT true,
    notify_sms boolean DEFAULT false,
    notify_whats_app boolean DEFAULT false,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    instagram character varying(255),
    youtube character varying(255),
    map_direction character varying(512)
);


ALTER TABLE public.general_settings OWNER TO postgres;

--
-- Name: monthly_sales; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.monthly_sales (
    year bigint NOT NULL,
    month bigint NOT NULL,
    total_sales bigint,
    total_revenue numeric,
    total_discount numeric,
    total_refunds numeric,
    net_revenue numeric,
    avg_order_value numeric,
    canceled_sales bigint,
    pending_sales bigint,
    completed_sales bigint,
    source_breakdown text,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.monthly_sales OWNER TO postgres;

--
-- Name: package_benefits; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.package_benefits (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    package_id uuid NOT NULL,
    title character varying(255) NOT NULL,
    description character varying(500),
    icon character varying(100),
    sort_order bigint DEFAULT 0,
    created_at timestamp with time zone
);


ALTER TABLE public.package_benefits OWNER TO postgres;

--
-- Name: packages; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.packages (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    name character varying(255) NOT NULL,
    description text,
    package_type character varying(50) NOT NULL,
    price numeric(10,2) NOT NULL,
    discount_price numeric(10,2) DEFAULT 0,
    duration_minutes bigint DEFAULT 60,
    total_sessions bigint DEFAULT 1,
    status character varying(20) DEFAULT 'active'::character varying NOT NULL,
    image_url character varying(500),
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    features text[],
    highlight boolean DEFAULT false,
    student_count bigint DEFAULT 0
);


ALTER TABLE public.packages OWNER TO postgres;

--
-- Name: pages; Type: TABLE; Schema: public; Owner: admin_drive
--

CREATE TABLE public.pages (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    title character varying(255) NOT NULL,
    slug character varying(255) NOT NULL,
    status character varying(20) DEFAULT 'draft'::character varying,
    sections text,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);


ALTER TABLE public.pages OWNER TO admin_drive;

--
-- Name: provinces; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.provinces (
    id bigint NOT NULL,
    name character varying(100) NOT NULL
);


ALTER TABLE public.provinces OWNER TO postgres;

--
-- Name: provinces_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.provinces_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.provinces_id_seq OWNER TO postgres;

--
-- Name: provinces_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.provinces_id_seq OWNED BY public.provinces.id;


--
-- Name: regencies; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.regencies (
    id bigint NOT NULL,
    province_id bigint NOT NULL,
    name character varying(100) NOT NULL,
    type character varying(100) NOT NULL
);


ALTER TABLE public.regencies OWNER TO postgres;

--
-- Name: regencies_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.regencies_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.regencies_id_seq OWNER TO postgres;

--
-- Name: regencies_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.regencies_id_seq OWNED BY public.regencies.id;


--
-- Name: related_articles; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.related_articles (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    article_id uuid,
    related_article_id uuid,
    relationship_type character varying(20)
);


ALTER TABLE public.related_articles OWNER TO postgres;

--
-- Name: sale_items; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.sale_items (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    sale_id uuid NOT NULL,
    package_id uuid NOT NULL,
    package_name character varying(255),
    quantity bigint DEFAULT 1,
    unit_price numeric(12,2) NOT NULL,
    discount numeric(12,2) DEFAULT 0,
    subtotal numeric(12,2) NOT NULL,
    created_at timestamp with time zone
);


ALTER TABLE public.sale_items OWNER TO postgres;

--
-- Name: sales; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.sales (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    order_number character varying(50) NOT NULL,
    payment_id uuid,
    user_id uuid NOT NULL,
    package_id uuid,
    package_name character varying(255),
    package_type character varying(50),
    total_amount numeric(12,2) NOT NULL,
    discount_amount numeric(12,2) DEFAULT 0,
    final_amount numeric(12,2) NOT NULL,
    status character varying(20) DEFAULT 'pending'::character varying,
    source character varying(50),
    payment_method character varying(50),
    currency character varying(3) DEFAULT 'IDR'::character varying,
    paid_at timestamp with time zone,
    refunded_at timestamp with time zone,
    notes text,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.sales OWNER TO postgres;

--
-- Name: tags; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.tags (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    name character varying(50) NOT NULL,
    slug character varying(50) NOT NULL,
    description character varying(255),
    usage_count bigint DEFAULT 0,
    created_at timestamp with time zone
);


ALTER TABLE public.tags OWNER TO postgres;

--
-- Name: districts id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.districts ALTER COLUMN id SET DEFAULT nextval('public.districts_id_seq'::regclass);


--
-- Name: provinces id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.provinces ALTER COLUMN id SET DEFAULT nextval('public.provinces_id_seq'::regclass);


--
-- Name: regencies id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.regencies ALTER COLUMN id SET DEFAULT nextval('public.regencies_id_seq'::regclass);


--
-- Data for Name: add_ons; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.add_ons (id, title, description, price, sessions, status, image_url, sort_order, created_at, updated_at) FROM stdin;
22222222-2222-2222-2222-222222222202	Night Session Add-on	Add night driving session to your package. Learn to drive safely in low-light conditions.	250000.00	0	active		2	2026-07-07 15:24:52.048007+07	2026-07-29 01:38:24.806571+07
22222222-2222-2222-2222-222222222203	Weekend Session Add-on	Add weekend driving session to your package. Perfect for those with busy weekday schedules.	200000.00	0	active		3	2026-07-07 15:24:52.048007+07	2026-07-29 01:38:24.84225+07
22222222-2222-2222-2222-222222222201	Extra Session	Additional driving session to enhance your skills. Each purchase adds 1 extra session to your package.	350000.00	1	active		1	2026-07-29 09:40:40.162023+07	2026-07-29 09:40:40.162023+07
\.


--
-- Data for Name: article_tags; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.article_tags (id, article_id, tag_id) FROM stdin;
\.


--
-- Data for Name: articles; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.articles (id, title, slug, author, content, media, lead_paragraph, excerpt, footer, meta_title, meta_description, meta_keywords, canonical_url, og_title, og_description, og_image, featured_image, thumbnail_url, category_id, tags, author_id, status, published_at, scheduled_at, view_count, like_count, share_count, reading_time, is_featured, is_spotlight, priority, highlight, created_at, updated_at, deleted_at, language, version, estimated_time, video_url, cta_text, cta_link, avg_read_time, bounce_rate, completion_rate, body_blocks) FROM stdin;
d7fdfb8a-d20d-4b0b-9bb1-2737b3d8e7e7	Tutorial: Drive Master Technique		Admin	<p>Key benefits of the <strong>Drive Master Technique</strong> include:</p><p></p><ul><li><p><strong>- Improved vehicle control</strong> during turns and lane changes.</p></li><li><p><strong>- Enhanced road safety</strong> through better awareness and reaction time.</p></li><li><p><strong>- Reduced fuel consumption</strong> by using smooth acceleration and braking techniques.</p></li><li><p><strong>- Greater driving confidence</strong> in both urban and highway environments.</p></li></ul><p></p><p>To practice the <strong>Drive Master Technique</strong>, drivers should focus on maintaining proper posture, keeping both hands on the steering wheel, and observing traffic conditions ahead. Regular practice helps develop muscle memory and quick decision-making skills. Over time, these techniques contribute to safer, more comfortable, and more efficient driving experiences for all road users.</p>	[]													\N	\N	00000000-0000-0000-0000-000000000000	published	\N	\N	0	0	0	1	t	f	0	f	2026-06-21 01:47:17.302037+07	2026-06-21 01:47:17.302037+07	2026-06-25 00:30:22.859856+07	en	1	0				0	0	0	\N
a0000001-0000-0000-0000-000000000001	5 Tips for First-Time Drivers to Pass Their SIM A Test	5-tips-first-time-drivers-pass-sim-a-test	Drive Master Team	Passing your SIM A driving test on the first try requires preparation and confidence. Here are five essential tips that have helped thousands of new drivers succeed.	[{"url": "https://example.com/images/driving-tips-hero.jpg", "name": "hero.jpg", "type": "image/jpeg"}]	Getting your SIM A license is a milestone. Here are expert tips to help you pass on your first attempt.										https://example.com/images/driving-tips-hero.jpg	https://example.com/images/driving-tips-thumb.jpg	c0000001-0000-0000-0000-000000000001	{tips,sim-a,beginners,driving-test}	00000000-0000-0000-0000-000000000000	published	2026-05-21 01:31:15.192634+07	\N	1525	89	45	5	t	f	10	t	2026-06-21 01:31:15.202321+07	2026-06-29 13:24:28.844278+07	2026-07-10 16:43:19.440414+07	en	1	0				0	0	0	\N
d6c1f196-34bf-408b-b36d-18c0dee11d7b	How to Pass Your Driving Test on First Try	how-to-pass-driving-test-first-try	Drive Master Team	Passing your driving test requires preparation and confidence. Here are the essential tips...	[{"id": "00000000-0000-0000-0000-000000000000", "url": "https://example.com/images/hero.jpg", "name": "hero-image.jpg", "size": "2MB", "type": "image/jpeg", "order": 1, "fileType": "image", "createdAt": "0001-01-01T00:00:00Z"}]	Learn the essential tips for passing your driving test on the first attempt.	Learn the essential tips for passing your driving test on the first attempt.									https://example.com/images/driving-tips.jpg		\N	\N	00000000-0000-0000-0000-000000000000	published	2024-01-15 17:00:00+07	\N	0	0	0	1	t	f	5	t	2026-06-21 01:32:30.94155+07	2026-06-21 01:32:30.94155+07	2026-07-10 16:43:28.723519+07	en	1	0				0	0	0	\N
a0000001-0000-0000-0000-000000000007	Parking Techniques for Beginners	parking-techniques-for-beginners	Drive Master Team	Good parking skills are essential for everyday driving. Practice these techniques regularly.	[{"url": "https://example.com/images/parking.jpg", "name": "parking.jpg", "type": "image/jpeg"}]	Parking is often the biggest challenge for new drivers. Master these techniques to park with confidence.										https://example.com/images/parking-hero.jpg		c0000001-0000-0000-0000-000000000001	{parking,skills,beginners,techniques}	00000000-0000-0000-0000-000000000000	draft	\N	\N	0	0	0	5	f	f	2	f	2026-06-21 01:31:15.202321+07	2026-06-21 01:31:15.202321+07	2026-07-10 16:43:09.358161+07	en	1	0				0	0	0	\N
a0000001-0000-0000-0000-000000000006	Understanding Car Maintenance Basics	understanding-car-maintenance-basics	Drive Master Team	Basic car maintenance knowledge helps prevent breakdowns and extends your vehicle's lifespan.	[{"url": "https://example.com/images/car-maintenance.jpg", "name": "car-maintenance.jpg", "type": "image/jpeg"}]	Every driver should know basic car maintenance. Keep your vehicle running smoothly with these essential tips.										https://example.com/images/car-maintenance-hero.jpg		c0000001-0000-0000-0000-000000000001	{maintenance,car-care,tips,knowledge}	00000000-0000-0000-0000-000000000000	draft	2026-06-14 01:31:15.192649+07	\N	278	19	8	1	f	f	3	f	2026-06-21 01:31:15.202321+07	2026-06-26 19:36:16.998714+07	2026-07-10 16:43:12.850653+07	en	1	0				0	0	0	\N
a0000001-0000-0000-0000-000000000005	The Benefits of Enrolling in a Professional Driving School	benefits-professional-driving-school	Drive Master Team	Professional driving schools offer advantages that self-learning simply cannot match. Learn about the benefits of professional training.	[{"url": "https://example.com/images/driving-school.jpg", "name": "driving-school.jpg", "type": "image/jpeg"}]	Professional driving schools offer advantages that self-learning simply cannot match.										https://example.com/images/driving-school-hero.jpg		c0000001-0000-0000-0000-000000000003	{driving-school,tips,beginners,preparation}	00000000-0000-0000-0000-000000000000	published	2026-06-06 01:31:15.192645+07	\N	323	24	15	4	f	f	6	f	2026-06-21 01:31:15.202321+07	2026-06-29 13:24:24.511518+07	2026-07-10 16:43:16.112634+07	en	1	0				0	0	0	\N
a0000001-0000-0000-0000-000000000004	How to Handle Wet Weather Driving	how-to-handle-wet-weather-driving	Drive Master Team	Rainy conditions require special driving techniques. This guide helps you navigate safely on wet roads.	[{"url": "https://example.com/images/rain-driving.jpg", "name": "rain-driving.jpg", "type": "image/jpeg"}]	Rainy conditions require special driving techniques. Here's how to stay safe on wet roads.										https://example.com/images/rain-driving-hero.jpg		c0000001-0000-0000-0000-000000000002	{safety,weather,rain,advanced-tips}	00000000-0000-0000-0000-000000000000	published	2026-02-21 01:31:15.192643+07	\N	445	28	12	5	f	f	4	f	2026-06-21 01:31:15.202321+07	2026-06-21 01:31:15.202321+07	2026-07-10 16:43:30.529176+07	en	1	0				0	0	0	\N
a0000001-0000-0000-0000-000000000003	Night Driving: Essential Tips for Safety	night-driving-essential-tips-for-safety	Drive Master Team	Night driving requires extra attention and different techniques than daytime driving. This guide covers essential safety tips for driving after dark.	[{"url": "https://example.com/images/night-driving.jpg", "name": "night-driving.jpg", "type": "image/jpeg"}]	Night driving presents unique challenges. Learn how to stay safe when driving after dark.										https://example.com/images/night-driving-hero.jpg		c0000001-0000-0000-0000-000000000002	{safety,night-driving,advanced,driving-tips}	00000000-0000-0000-0000-000000000000	published	2026-03-21 01:31:15.19264+07	\N	654	42	18	4	f	f	5	f	2026-06-21 01:31:15.202321+07	2026-06-21 01:31:15.202321+07	2026-07-10 16:43:32.413084+07	en	1	0				0	0	0	\N
24b1fd20-b527-422b-8d50-f01c69ab117d	Formula 1: Puncak Balap Mobil Dunia	formula-1-puncak-balap-mobil-dunia	\N	<p>Persaingan antara tim-tim besar seperti Ferrari, Mercedes-AMG Petronas Formula One Team, dan Oracle Red Bull Racing sering kali menjadi sorotan utama setiap musim.</p><p>Popularitas <strong>Formula 1</strong> terus meningkat berkat kombinasi antara olahraga, teknologi, dan hiburan. Balapan diselenggarakan di berbagai sirkuit ikonik di seluruh dunia, mulai dari Circuit de Monaco hingga Silverstone Circuit. Dengan munculnya generasi pembalap muda berbakat serta perkembangan teknologi yang berkelanjutan, Formula 1 diperkirakan akan tetap menjadi simbol inovasi dan kompetisi tingkat tertinggi dalam dunia motorsport.</p><p></p>	[]	Formula 1 (F1) merupakan ajang balap mobil paling bergengsi di dunia yang mempertemukan pembalap dan tim terbaik dari berbagai negara. Kejuaraan ini diselenggarakan oleh Fédération Internationale de l'Automobile (FIA) dan berlangsung dalam serangkaian balapan yang disebut Grand Prix. Dengan teknologi canggih, kecepatan tinggi, dan strategi yang kompleks, Formula 1 menjadi salah satu olahraga yang paling banyak ditonton di dunia.	\N									https://ik.imagekit.io/oy4rsvid5/blogs/24b1fd20-b527-422b-8d50-f01c69ab117d/WhatsApp_Image_2026-06-26_at_16.09.45_E-QN27jBm.jpeg?tr=w-700,h-450,fo-auto		\N	\N	cf475ead-91bf-4a55-a47e-cf93279240b6	draft	\N	\N	0	0	0	1	f	f	0	f	2026-06-26 11:09:21.366407+07	2026-06-26 18:52:59.625813+07	\N	en	1	0				0	0	0	\N
8b458099-05a6-4b11-a758-02653f665a8b	Tutorial: Drive Master Technique	tutorial-drive-master-technique	\N	<p><strong>Key benefits of the Drive Master Technique include:</strong></p><p>Improved vehicle control during turns and lane changes.</p><p>Enhanced road safety through better awareness and reaction time.</p><p>Reduced fuel consumption by using smooth acceleration and braking techniques.</p><p>Greater driving confidence in both urban and highway environments.</p><p>To practice the Drive Master Technique, drivers should focus on maintaining proper posture, keeping both hands on the steering wheel, and observing traffic conditions ahead. Regular practice helps develop muscle memory and quick decision-making skills. Over time, these techniques contribute to safer, more comfortable, and more efficient driving experiences for all road users.</p><p></p>	[]	The Drive Master Technique is a driving method designed to improve vehicle control, safety, and efficiency. By mastering essential skills such as steering, acceleration, and braking, drivers can navigate various road conditions with greater confidence. This technique is especially beneficial for new drivers who want to build strong driving habits and experienced drivers seeking to refine their abilities.	\N									https://ik.imagekit.io/oy4rsvid5/blogs/8b458099-05a6-4b11-a758-02653f665a8b/WhatsApp_Image_2026-06-25_at_21.54.40_6jD3Ypf7e.jpeg?tr=w-700,h-450,fo-auto		\N	\N	cf475ead-91bf-4a55-a47e-cf93279240b6	published	\N	\N	0	0	0	1	f	f	0	f	2026-06-25 01:03:10.313982+07	2026-06-26 14:20:00.159463+07	2026-06-26 18:57:27.889324+07	en	1	0				0	0	0	\N
a0000001-0000-0000-0000-000000000002	Understanding Indonesian Road Signs: A Complete Guide	understanding-indonesian-road-signs-complete-guide	Drive Master Team	Indonesian roads can be challenging for new drivers. Understanding road signs is crucial for safe navigation. This comprehensive guide covers all essential road signs.	[{"url": "https://example.com/images/road-signs.jpg", "name": "road-signs.jpg", "type": "image/jpeg"}]	Navigate Indonesian roads with confidence by mastering these essential traffic signs and their meanings.										https://example.com/images/road-signs-hero.jpg		c0000001-0000-0000-0000-000000000002	{road-safety,traffic-rules,beginners,indonesia}	00000000-0000-0000-0000-000000000000	published	2026-04-21 01:31:15.192638+07	\N	893	56	23	7	f	t	8	f	2026-06-21 01:31:15.202321+07	2026-06-29 11:08:41.864818+07	2026-07-10 16:43:34.342705+07	en	1	0				0	0	0	\N
5afc48dc-c21a-4dbf-882a-75ab145f118e	Tutorial: Drive Master Technique	tutorial-drive-master-technique	\N	<p>Tahap 1 – Pengenalan Kendaraan</p><p>Materi meliputi:</p><p>• Posisi mengemudi.</p><p>• Fungsi pedal dan tuas.</p><p>• Instrumen kendaraan.</p><p>• Pemeriksaan dasar kendaraan.</p><p>• Prosedur keselamatan.</p><p>Tahap 2 – Pengendalian Dasar</p><p>Materi meliputi:</p><p>• Menjalankan kendaraan.</p><p>• Menghentikan kendaraan.</p><p>• Pengendalian kecepatan.</p><p>• Penggunaan spion.</p><p>• Teknik kemudi dasar.</p><p>Tahap 3 – Manuver dan Parkir</p><p>Materi meliputi:</p><p>• Parkir paralel.</p><p>• Parkir mundur.</p><p>• Parkir tegak lurus.</p><p>• Putar balik.</p><p>• Manuver pada area sempit.</p><p>Tahap 4 – Berkendara di Jalan Raya</p><p>Materi meliputi:</p><p>• Perpindahan jalur.</p><p>• Membaca situasi lalu lintas.</p><p>• Pengambilan keputusan.</p><p>• Berkendara pada kondisi ramai.</p><p>Tahap 5 – Persiapan Berkendara Mandiri</p><p>Materi meliputi:</p><p>• Simulasi kondisi nyata.</p><p>• Penguatan kepercayaan diri.</p><p>• Evaluasi akhir.</p><p>• Persiapan graduation.</p><p>Struktur ini memastikan setiap siswa memperoleh pengalaman belajar yang sistematis dan progresif.</p>	[]	DriveMaster dirancang secara bertahap untuk membantu siswa berkembang dari tingkat\npemula hingga siap berkendara secara mandiri.	\N									https://ik.imagekit.io/oy4rsvid5/blogs/5afc48dc-c21a-4dbf-882a-75ab145f118e/3_Ss9SEzyg4.png?tr=w-700,h-450,fo-auto		\N	\N	cf475ead-91bf-4a55-a47e-cf93279240b6	published	\N	\N	15	0	0	1	f	f	0	f	2026-06-26 19:37:43.13784+07	2026-08-04 10:03:14.708505+07	\N	en	1	0				0	0	0	\N
\.


--
-- Data for Name: cars; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.cars (id, brand, model, year, license_plate, color, transmission, status, mileage, image_url, notes, created_at, updated_at) FROM stdin;
898170a5-08db-4467-b33e-049660a4231c	BYD	Atto 1	2024	B 1234 XYZ		manual	available	0			2026-06-19 19:56:29.979452+07	2026-06-19 19:56:29.979452+07
c0000001-0000-0000-0000-000000000001	Toyota	Vios	2023	B 1234 ABC	White	automatic	available	5000	https://example.com/toyota-vios.jpg	Clean condition, AC working	2026-07-13 09:03:10.157535+07	2026-07-13 09:03:10.157535+07
c0000001-0000-0000-0000-000000000002	Honda	Civic	2022	B 5678 DEF	Black	automatic	available	12000	https://example.com/honda-civic.jpg	Full option, leather seats	2026-07-13 09:03:10.157535+07	2026-07-13 09:03:10.157535+07
c0000001-0000-0000-0000-000000000003	Suzuki	Baleno	2023	D 9012 GHI	Silver	manual	available	3000	https://example.com/suzuki-baleno.jpg	Fuel efficient, good for beginner	2026-07-13 09:03:10.157535+07	2026-07-13 09:03:10.157535+07
c0000001-0000-0000-0000-000000000004	Toyota	Yaris	2021	F 3456 JKL	Red	automatic	maintenance	25000	https://example.com/toyota-yaris.jpg	Under maintenance - scheduled service	2026-07-13 09:03:10.157535+07	2026-07-13 09:03:10.157536+07
c0000001-0000-0000-0000-000000000005	Honda	City	2024	B 7890 MNO	Gray	automatic	in_use	1000	https://example.com/honda-city.jpg	Brand new, currently in use for class	2026-07-13 09:03:10.157536+07	2026-07-13 09:03:10.157536+07
c0000001-0000-0000-0000-000000000006	Mitsubishi	Mirage	2022	L 1234 PQR	Blue	manual	available	15000	https://example.com/mitsubishi-mirage.jpg	Compact, easy to park	2026-07-13 09:03:10.157536+07	2026-07-13 09:03:10.157536+07
c0000001-0000-0000-0000-000000000007	Toyota	Avanza	2023	B 5678 STU	White	manual	available	8000	https://example.com/toyota-avanza.jpg	Family car, 7 seats	2026-07-13 09:03:10.157536+07	2026-07-13 09:03:10.157536+07
c0000001-0000-0000-0000-000000000008	Honda	HR-V	2024	D 9012 VWX	Black	automatic	available	500	https://example.com/honda-hr-v.jpg	SUV, great for long distance	2026-07-13 09:03:10.157536+07	2026-07-13 09:03:10.157536+07
\.


--
-- Data for Name: categories; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.categories (id, name, slug, description, parent_id, "order", created_at, updated_at) FROM stdin;
c0000001-0000-0000-0000-000000000001	Driving Tips	driving-tips	Practical advice for better driving	\N	1	2026-06-05 09:08:59.883807+07	2026-06-05 09:08:59.883807+07
c0000001-0000-0000-0000-000000000002	Road Safety	road-safety	Stay safe on the roads	\N	2	2026-06-05 09:08:59.889235+07	2026-06-05 09:08:59.889235+07
c0000001-0000-0000-0000-000000000003	Beginners Guide	beginners-guide	Everything new drivers need to know	\N	3	2026-06-05 09:08:59.893258+07	2026-06-05 09:08:59.893258+07
\.


--
-- Data for Name: contact_inquiries; Type: TABLE DATA; Schema: public; Owner: admin_drive
--

COPY public.contact_inquiries (id, name, email, subject, message, created_at, updated_at, status) FROM stdin;
9de54ccf-b613-467e-8296-ac39f10118c1	John Tester	tester@gmail.com	Package Information	Can I get more info?	2026-07-15 07:51:44.070395+07	2026-07-15 07:51:44.070396+07	unread
73ee1c37-32cf-49b2-a230-83b12a610883	Test User	test@example.com	Test Inquiry	Hello, this is a test message.	2026-07-14 11:00:08.78432+07	2026-07-15 14:32:49.589526+07	read
c9db6ef9-b090-469b-8b5a-74a5c72cef95	Zeta	muhammadrizqiko@gmail.com	General Inquiry	Saya ingin tidur	2026-07-10 17:26:56.164773+07	2026-07-15 14:32:57.750453+07	read
9ea83618-495b-4478-bb00-61927b562a69	afifah	afifahrahmawati221@gmail.com	Scheduling Issue	tolong buatkan saya jadwal kursusnya	2026-07-23 08:40:17.891083+07	2026-07-23 08:40:17.891083+07	unread
6ef3e693-9063-4b07-898a-78e3f456501b	Addzkiya	muhammadrizqiko@gmail.com	Technical Support	Saya ingin menjadi sukses	2026-07-24 16:20:18.717434+07	2026-07-24 16:20:18.717434+07	unread
3dd710d4-29ae-432c-bb52-77b27805f2a6	Wika Salim	wikasalim@gmail.com	Package Information	Saya mau tidur	2026-08-05 23:43:44.629179+07	2026-08-05 23:43:44.629179+07	unread
\.


--
-- Data for Name: districts; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.districts (id, name, regency_id) FROM stdin;
1	2 x 11 Enam Lingkuang	13
2	2 x 11 Kayu Tanam	13
3	Abab	16
4	Abad Selatan	53
5	Abang	51
6	Abeli	74
7	Abenaho	95
8	Abepura	91
9	Abiansemal	51
10	Aboy	95
11	Abuki	74
12	Abun	92
13	Abung Barat	18
14	Abung Kunang	18
15	Abung Pekurun	18
16	Abung Selatan	18
17	Abung Semuli	18
18	Abung Surakarta	18
19	Abung Tengah	18
20	Abung Timur	18
21	Abung Tinggi	18
22	Adian Koting	12
23	Adiluwih (Adi Luwih)	18
24	Adimulyo	33
25	Adipala	33
26	Adiwerna	33
27	Adonara	53
28	Adonara Barat	53
29	Adonara Tengah	53
30	Adonara Timur	53
31	Aek Bilah	12
32	Aek Kuasan	12
33	Aek Kuo	12
34	Aek Ledong	12
35	Aek Nabara Barumun	12
36	Aek Natas	12
37	Aek Songsongan	12
38	Aere	74
39	Aertembaga (Bitung Timur)	71
40	Aesesa	53
41	Aesesa Selatan	53
42	Afulu	12
43	Agandugume	94
44	Agats	93
45	Agimuga	94
46	Agisiga	94
47	Agrabinta	32
48	Aifat	92
49	Aifat Selatan	92
50	Aifat Timur	92
51	Aifat Timur Jauh	92
52	Aifat Timur Selatan	92
53	Aifat Timur Tengah	92
54	Aifat Utara	92
55	Aikmel	52
56	Aimando Padaido	91
57	Aimas	92
58	Aimere	53
59	Air Batu	12
60	Air Besar	61
61	Air Besi	17
62	Air Buaya (Airbuaya)	81
63	Air Dikit	17
64	Air Hangat	15
65	Air Hangat Barat	15
66	Air Hangat Timur	15
67	Air Hitam	15
68	Air Hitam	18
69	Air Joman	12
70	Air Kumbang	16
71	Air Majunto	17
72	Air Naningan	18
73	Air Napal	17
74	Air Nipis	17
75	Air Padang	17
76	Air Periukan	17
77	Air Putih	12
78	Air Rami	17
79	Air Salek	16
80	Air Sugihan	16
81	Air Upas	61
82	Airgaram	95
83	Airgegas (Air Gegas)	19
84	Airmadidi	71
85	Airpura	13
86	Airu	91
87	Aitinyo	92
88	Aitinyo Barat	92
89	Aitinyo Raya	92
90	Aitinyo Tengah	92
91	Aitinyo Utara	92
92	Ajangale	73
93	Ajibarang	33
94	Ajibata	12
95	Ajung	35
96	Akabiluru	13
97	Akat	93
98	Alafan (Alapan)	11
99	Alak	53
100	Alalak	63
101	Alam Barajo	15
102	Alama	95
103	Alama	94
104	Alang Alang Lebar	16
105	Alas	52
106	Alas Barat	52
107	Alasa	12
108	Alasa Talumuzoi	12
109	Alemsom	95
110	Alian (Aliyan)	33
111	Alla	73
112	Allu (Alu)	76
113	Alok	53
114	Alok Barat	53
115	Alok Timur	53
116	Alor Barat Daya	53
117	Alor Barat Laut	53
118	Alor Selatan	53
119	Alor Tengah Utara	53
120	Alor Timur	53
121	Alor Timur Laut	53
122	Aluh Aluh	63
123	Amabi Oefeto	53
124	Amabi Oefeto Timur	53
125	Amahai	81
126	Amalatu	81
127	Amali	73
128	Amanatun Selatan	53
129	Amanatun Utara	53
130	Amandraya	12
131	Amanuban Barat	53
132	Amanuban Selatan	53
133	Amanuban Tengah	53
134	Amanuban Timur	53
135	Amar	94
136	Amarasi	53
137	Amarasi Barat	53
138	Amarasi Selatan	53
139	Amarasi Timur	53
140	Ambal	33
141	Ambalau	61
142	Ambalau	81
143	Ambalawi	52
144	Ambarawa	18
145	Ambarawa	33
146	Ambatkwi (Ambatkui)	93
147	Amberbaken	92
148	Amberbaken Barat	92
149	Ambulu	35
150	Ambunten	35
151	Amen	17
152	Amfoang Barat Daya	53
153	Amfoang Barat Laut	53
154	Amfoang Selatan	53
155	Amfoang Tengah	53
156	Amfoang Timur	53
157	Amfoang Utara	53
158	Amonggedo	74
159	Ampana Kota	72
160	Ampana Tete	72
161	Ampek Angkek (IV Angkat Candung)	13
162	Ampek Nagari (IV Nagari )	13
163	Ampel	33
164	Ampelgading	33
165	Ampelgading	35
166	Ampenan	52
167	Ampibabo	72
168	Amuma	95
169	Amungkalpia	94
170	Amuntai Selatan	63
171	Amuntai Tengah	63
172	Amuntai Utara	63
173	Amurang	71
174	Amurang Barat	71
175	Amurang Timur	71
176	Anak Ratu Aji	18
177	Anak Tuha	18
178	Anawi	95
179	Andam Dewi	12
180	Andey	91
181	Andir	32
182	Andong	33
183	Andoolo	74
184	Andoolo Barat	74
185	Andowia	74
186	Angata	74
187	Anggaberi	74
188	Anggalomoare	74
189	Anggana	64
190	Anggeraja	73
191	Anggi	92
192	Anggi Gida	92
193	Anggrek	75
194	Anggruk	95
195	Angkaisera	91
196	Angkinang	63
197	Angkola Barat	12
198	Angkola Muara Tais	12
199	Angkola Sangkunur	12
200	Angkola Selatan	12
201	Angkola Timur	12
202	Angkona	73
203	Angsana	63
204	Angsana	36
205	Animha	93
206	Anjatan	32
207	Anjir Muara	63
208	Anjir Pasar	63
209	Anjongan	61
210	Anotaurei	91
211	Anreapi	76
212	Antang Kalang	62
213	Antapani (Cicadas)	32
214	Anyar	36
215	Apalapsili	95
216	Apawer Hulu	91
217	Aradide	94
218	Arahan	32
219	Aralle	76
220	Aramo	12
221	Aranday	92
222	Aranio	63
223	Arcamanik	32
224	Argapura	32
225	Argomulyo	33
226	Arguni	92
227	Arimop	93
228	Arjasa	35
229	Arjasari	32
230	Arjawinangun	32
231	Arjosari	35
232	Arma Jaya	17
233	Aroba	92
234	Arongan Lambalek	11
235	Arosbaya	35
236	Arse	12
237	Arso	91
238	Arso Barat	91
239	Arso Timur	91
240	Aru Selatan	81
241	Aru Selatan Timur	81
242	Aru Selatan Utara	81
243	Aru Tengah	81
244	Aru Tengah Selatan	81
245	Aru Tengah Timur	81
246	Aru Utara	81
247	Aru Utara Timur Batuley	81
248	Arungkeke	73
249	Arut Selatan	62
250	Arut Utara	62
251	Asakota	52
252	Asam Jujuhan	13
253	Asem Rowo (Asemrowo)	35
254	Asembagus	35
255	Asera	74
256	Ases	92
257	Asinua	74
258	Asologaima	95
259	Asolokobal	95
260	Asotipo	95
261	Asparaga	75
262	Assue	93
263	Astambul	63
264	Astana Anyar	32
265	Astanajapura	32
266	Aswi	93
267	Atadei	53
268	Atambua Barat	53
269	Atambua Selatan	53
270	Atinggola	75
271	Atsj (Atsy)	93
272	Atu Lintang	11
273	Aur Birugo Tigo Baleh	13
274	Awan Rante Karua	73
275	Awang	62
276	Awangpone	73
277	Awayan	63
278	Aweida	94
279	Aweku	95
280	Awina	95
281	Awinbon	95
282	Awyu	93
283	Ayah	33
284	Ayamaru	92
285	Ayamaru Barat	92
286	Ayamaru Jaya	92
287	Ayamaru Selatan	92
288	Ayamaru Selatan Jaya	92
289	Ayamaru Tengah	92
290	Ayamaru Timur	92
291	Ayamaru Timur Selatan	92
292	Ayamaru Utara	92
293	Ayamaru Utara Timur	92
294	Ayau	92
295	Ayip	93
296	Ayumnati	95
297	Baamang	62
298	Babadan	35
299	Babah Rot	11
300	Babakan	32
301	Babakan Ciparay	32
302	Babakan Madang	32
303	Babakancikao	32
304	Babalan	12
305	Babar Barat (Pulau Pulau Babar)	81
306	Babat	35
307	Babat Supat	16
308	Babat Toman	16
309	Babelan	32
310	Babirik	63
311	Babo	92
312	Babul Makmur	11
313	Babul Rahmah	11
314	Babulu	64
315	Babussalam	11
316	Bacan	82
317	Bacan Barat	82
318	Bacan Barat Utara	82
319	Bacan Selatan	82
320	Bacan Timur	82
321	Bacan Timur Selatan	82
322	Bacan Timur Tengah	82
323	Bacukiki	73
324	Bacukiki Barat	73
325	Badar	11
326	Badas	35
327	Badau	19
328	Badau	61
329	Badegan	35
330	Badiri	12
331	Bae	33
332	Baebunta	73
333	Baebunta Selatan	73
334	Bagan Sinembah Raya	14
335	Bagansinembah (Bagan Sinembah)	14
336	Bagelen	33
337	Bagor	35
338	Baguala	81
339	Bagun	92
340	Bahar Selatan	15
341	Bahar Utara	15
342	Bahau Hulu	65
343	Bahodopi	72
344	Bahorok	12
345	Bahuga	18
346	Baito	74
347	Baiturrahman	11
348	Baitussalam	11
349	Bajawa	53
350	Bajawa Utara	53
351	Bajeng	73
352	Bajeng Barat	73
353	Bajenis	12
354	Bajo	73
355	Bajo Barat	73
356	Bajubang	15
357	Bajuin	63
358	Bakam	19
359	Bakarangan	63
360	Bakauheni	18
361	Baki	33
362	Bakongan	11
363	Bakongan Timur	11
364	Baktiraja (Bakti Raja)	12
365	Baktiya	11
366	Baktiya Barat	11
367	Bakumpai	63
368	Bakung	35
369	Bakung Serumpun	21
370	Balaesang	72
371	Balaesang Tanjung	72
372	Balai	61
373	Balai Jaya	14
374	Balai Riam	62
375	Balanipa	76
376	Balantak	72
377	Balantak Selatan	72
378	Balantak Utara	72
379	Balapulang	33
380	Balaraja	36
381	Baleendah	32
382	Balen	35
383	Balerejo	35
384	Balige	12
385	Balik Bukit	18
386	Balikpapan Barat	64
387	Balikpapan Kota	64
388	Balikpapan Selatan	64
389	Balikpapan Tengah	64
390	Balikpapan Timur	64
391	Balikpapan Utara	64
392	Balingga	95
393	Balingga Barat	95
394	Balinggi	72
395	Balla	76
396	Balocci	73
397	Balong	35
398	Balongan	32
399	Balongbendo	35
400	Balongpanggang (Balong Panggang)	35
401	Balung	35
402	Balusu	73
403	Bambaira	76
404	Bambalamotu	76
405	Bambang	76
406	Bambanglipuro (Bambang Lipuro)	34
407	Bambel	11
408	Bamgi	93
409	Bamusbama	92
410	Banama Tingang	62
411	Banawa	72
412	Banawa Selatan	72
413	Banawa Tengah	72
414	Bancak	33
415	Bancar	35
416	Banda	81
417	Banda Alam	11
418	Banda Baro	11
419	Banda Mulia	11
420	Banda Raya	11
421	Banda Sakti	11
422	Bandar	35
423	Bandar	11
424	Bandar	33
425	Bandar	12
426	Bandar Baru	11
427	Bandar Dua	11
428	Bandar Huluan	12
429	Bandar Khalipah / Khalifah	12
430	Bandar Laksamana	14
431	Bandar Masilam	12
432	Bandar Mataram	18
433	Bandar Negeri Semuong	18
434	Bandar Negeri Suoh	18
435	Bandar Pasir Mandoge	12
436	Bandar Petalangan	14
437	Bandar Pulau	12
438	Bandar Pusaka	11
439	Bandar Sei Kijang	14
440	Bandar Sribhawono (Bandar Sribawono)	18
441	Bandar Surabaya	18
442	Bandarkedungmulyo (Bandar Kedung Mulyo)	35
443	Banding Agung	16
444	Bandongan	33
445	Bandung	36
446	Bandung	35
447	Bandung Kidul	32
448	Bandung Kulon	32
449	Bandung Wetan	32
450	Bandungan	33
451	Bang Haji	17
452	Banggae	76
453	Banggae Timur	76
454	Banggai	72
455	Banggai Selatan	72
456	Banggai Tengah	72
457	Banggai Utara	72
458	Bangil	35
459	Bangilan	35
460	Bangkala	73
461	Bangkala Barat	73
462	Bangkalan	35
463	Bangkelekila	73
464	Bangkinang	14
465	Bangkinang Kota	14
466	Bangko	14
467	Bangko	15
468	Bangko Barat	15
469	Bangko Pusako / Pusaka	14
470	Bangkunat (Bengkunat)	18
471	Bangkurung	72
472	Bangli	51
473	Bangodua	32
474	Bangorejo	35
475	Bangsal	35
476	Bangsalsari	35
477	Bangsri	33
478	Bangun Purba	12
479	Bangun Purba	14
480	Bangun Rejo	18
481	Banguntapan	34
482	Banjang	63
483	Banjar	51
484	Banjar	36
485	Banjar	32
486	Banjar Agung	18
487	Banjar Baru	18
488	Banjar Margo	18
489	Banjaran	32
490	Banjarangkan	51
491	Banjaranyar	32
492	Banjarbaru Selatan (Banjar Baru Selatan)	63
493	Banjarbaru Utara (Banjar Baru Utara)	63
494	Banjarejo	33
495	Banjarharjo	33
496	Banjarmangu	33
497	Banjarmasin Barat	63
498	Banjarmasin Selatan	63
499	Banjarmasin Tengah	63
500	Banjarmasin Timur	63
501	Banjarmasin Utara	63
502	Banjarnegara	33
503	Banjarsari	36
504	Banjarsari	32
505	Banjarsari	33
506	Banjarwangi	32
507	Banjit	18
508	Bansari	33
509	Bantaeng	73
510	Bantan	14
511	Bantaran	35
512	Bantarbolang	33
513	Bantargadung	32
514	Bantargebang (Bantar Gebang)	32
515	Bantarkalong	32
516	Bantarkawung	33
517	Bantarsari	33
518	Bantarujeg	32
519	Bantimurung	73
520	Bantul	34
521	Bantur	35
522	Banua Lawas	63
523	Banua Lima	62
524	Banuhampu	13
525	Banyakan	35
526	Banyuanyar (Banyu Anyar)	35
527	Banyuasin I	16
528	Banyuasin II	16
529	Banyuasin III	16
530	Banyuates	35
531	Banyubiru	33
532	Banyudono	33
533	Banyuglugur	35
534	Banyuke Hulu	61
535	Banyumanik	33
536	Banyumas	33
537	Banyumas	18
538	Banyuputih	35
539	Banyuputih	33
540	Banyuresmi	32
541	Banyusari	32
542	Banyuurip	33
543	Banyuwangi	35
544	Baolan	72
545	Bara	73
546	Barabai	63
547	Baradatu	18
548	Baraka	73
549	Barambai	63
550	Barangin	13
551	Barangka	74
552	Baranti	73
553	Baras	76
554	Barat	35
555	Barebbo	73
556	Baregbeg	32
557	Bareng	35
558	Barito Tuhup Raya	62
559	Baroko	73
560	Barombong	73
561	Baron	35
562	Barong Tongkok	64
563	Baros	36
564	Baros	32
565	Barru	73
566	Baruga	74
567	Barumun	12
568	Barumun Barat	12
569	Barumun Baru	12
570	Barumun Selatan	12
571	Barumun Tengah	12
572	Baruppu	73
573	Barus	12
574	Barus Utara	12
575	Barusjahe (Barus Jahe)	12
576	Basa Ampek Balai Tapan	13
577	Basala	74
578	Basarang	62
579	Basidondo	72
580	Baso	13
581	Basse Sangtempe (Bassesang Tempe / Bastem)	73
582	Basse Sangtempe Utara	73
583	Batabual	81
584	Bataguh	62
585	Batahan	12
586	Batalaiworu (Batalaiwaru)	74
587	Batam Kota	21
588	Batang	33
589	Batang	73
590	Batang Alai Selatan	63
591	Batang Alai Timur	63
592	Batang Alai Utara	63
593	Batang Anai	13
594	Batang Angkola	12
595	Batang Asai	15
596	Batang Asam	15
597	Batang Batang	35
598	Batang Cenaku	14
599	Batang Gangsal (Batang Gansal)	14
600	Batang Gasan	13
601	Batang Hari Leko	16
602	Batang Kapas	13
603	Batang Kawa	62
604	Batang Kuis	12
605	Batang Lubu Sutam	12
606	Batang Lupar	61
607	Batang Masumai	15
608	Batang Merangin	15
609	Batang Natal	12
610	Batang Onang	12
611	Batang Peranap	14
612	Batang Serangan	12
613	Batang Toru	12
614	Batang Tuaka	14
615	Batangan	33
616	Batanghari	18
617	Batanghari Nuban	18
618	Batani	95
619	Batanta Selatan	92
620	Batanta Utara	92
621	Batauga	74
622	Batealit	33
623	Batee	11
624	Bathin II Pelayang	15
625	Bathin III	15
626	Bathin III Ulu	15
627	Bathin Solapan	14
628	Bathin VIII (Batin VIII)	15
629	Bati Bati	63
630	Batik Nau	17
631	Batin II Babeko (Bathin)	15
632	Batin XXIV	15
633	Batipuah Selatan (Batipuh Selatan)	13
634	Batipuh	13
635	Batom	95
636	Batu	35
637	Batu Aji	21
638	Batu Ampar	61
639	Batu Ampar	63
640	Batu Ampar	64
641	Batu Ampar	21
642	Batu Ampar	62
643	Batu Atas	74
644	Batu Benawa	63
645	Batu Brak	18
646	Batu Engau	64
647	Batu Hampar	14
648	Batu Ketulis	18
649	Batu Lanteh (Batulanteh)	52
650	Batu Layar	52
651	Batu Licin (Batulicin)	63
652	Batu Mandi	63
653	Batu Putih	64
654	Batu Putih	53
655	Batu Putih	74
656	Batu Putih	18
657	Batu Sopang	64
658	Batuan	35
659	Batuceper	36
660	Batudaa	75
661	Batudaa Pantai	75
662	Batudaka	72
663	Batui	72
664	Batui Selatan	72
665	Batujajar	32
666	Batujaya	32
667	Batukara	74
668	Batukliang	52
669	Batukliang Utara	52
670	Batulappa (Batu Lappa)	73
671	Batumarmar	35
672	Batununggal	32
673	Batupoaro	74
674	Batuputih	35
675	Batur	33
676	Baturaja Barat	16
677	Baturaja Timur	16
678	Baturetno	33
679	Baturiti	51
680	Baturraden (Baturaden)	33
681	Batuwarno	33
682	Baula	74
683	Baureno	35
684	Bawang	33
685	Bawen	33
686	Bawolato	12
687	Baya Biru	94
688	Bayah	36
689	Bayan	33
690	Bayan	52
691	Bayang	13
692	Bayat	33
693	Bayongbong	32
694	Bayung Lencir	16
695	Bebandem	51
696	Beber	32
697	Bebesen	11
698	Beduai (Beduwai)	61
699	Bejen	33
700	Beji	32
701	Beji	35
702	Bekasi Barat	32
703	Bekasi Selatan	32
704	Bekasi Timur	32
705	Bekasi Utara	32
706	Bekri	18
707	Belakang Padang	21
708	Belalau	18
709	Belang	71
710	Belantikan Raya	62
711	Belat	21
712	Belawa	73
713	Belawang	63
714	Belida Darat	16
715	Belik	33
716	Belimbing	61
717	Belimbing	16
718	Belimbing Hulu	61
719	Belinyu	19
720	Belitang	61
721	Belitang	16
722	Belitang Hilir	61
723	Belitang Hulu	61
724	Belitang II	16
725	Belitang III	16
726	Belitang Jaya	16
727	Belitang Madang Raya	16
728	Belitang Mulya	16
729	Belo	52
730	Belopa	73
731	Belopa Utara	73
732	Benai	14
733	Benakat	16
734	Benawa	95
735	Benda	36
736	Bendahara	11
737	Bendo	35
738	Bendosari	33
739	Bendungan	35
740	Bener	33
741	Bener Kelipah	11
742	Bengalon	64
743	Bengkalis	14
744	Bengkayang	61
745	Bengkong	21
746	Bengo	73
747	Benjeng	35
748	Benowo	35
749	Benteng	73
750	Bentian Besar	64
751	Benua	74
752	Benua Kayong	61
753	Benuki	91
754	Beo	71
755	Beo Selatan	71
756	Beo Utara	71
757	Beoga	94
758	Beoga Barat	94
759	Beoga Timur	94
760	Berampu (Brampu)	12
761	Berandan Barat (Brandan Barat)	12
762	Berastagi (Brastagi)	12
763	Beraur	92
764	Berbah	34
765	Berbak	15
766	Berbek	35
767	Bergas	33
768	Beringin	12
769	Bermani Ilir	17
770	Bermani Ulu	17
771	Bermani Ulu Raya	17
772	Beruntung Baru	63
773	Besitang	12
774	Besuk	35
775	Besuki	35
776	Besulutu	74
777	Betara	15
778	Betayau	65
779	Betcbamu	93
780	Betoambari	74
781	Betung	16
782	Beutong	11
783	Beutong Ateuh Banggalang	11
784	Bewani	95
785	Biak Barat	91
786	Biak Kota	91
787	Biak Timur	91
788	Biak Utara	91
789	Biandoga	94
790	Biaro	71
791	Biatan	64
792	Biau	75
793	Biau	72
794	Bibida	94
795	Biboki Anleu	53
796	Biboki Feotleu	53
797	Biboki Moenleu	53
798	Biboki Selatan	53
799	Biboki Tan Pah	53
800	Biboki Utara	53
801	Biduk-Biduk	64
802	Bies	11
803	Bika	61
804	Bikar	92
805	Bikomi Nilulat	53
806	Bikomi Selatan	53
807	Bikomi Tengah	53
808	Bikomi Utara	53
809	Bilah Barat	12
810	Bilah Hilir	12
811	Bilah Hulu	12
812	Bilalang	71
813	Bilato	75
814	Biluhu	75
815	Bime	95
816	Bina	94
817	Binakal	35
818	Binamu	73
819	Binangun	33
820	Binangun	35
821	Binawidya	14
822	Binduriang	17
823	Bingin Kuning	17
824	Binjai	12
825	Binjai Barat	12
826	Binjai Hulu	61
827	Binjai Kota	12
828	Binjai Selatan	12
829	Binjai Timur	12
830	Binjai Utara	12
831	Binong	32
832	Binongko	74
833	Bintan Pesisir	21
834	Bintan Timur	21
835	Bintan Utara	21
836	Bintang	11
837	Bintang Ara	63
838	Bintang Bayu	12
839	Bintauna	71
840	Bintuni	92
841	Binuang	76
842	Binuang	63
843	Binuang	36
844	Birem Bayeun	11
845	Biringbulu	73
846	Biringkanaya (Biring Kanaya)	73
847	Biru-Biru (Sibiru-biru)	12
848	Biscoop	92
849	Bissappu	73
850	Bittuang	73
851	Biuk	95
852	Bl. Limbangan (Blubur Limbangan)	32
853	Blado	33
854	Blahbatuh (Belah Batuh)	51
855	Blambangan Pagar	18
856	Blambangan Umpu	18
857	Blanakan	32
858	Blang Bintang (Blank Bintang)	11
859	Blang Mangat	11
860	Blangjerango (Blang Jerango)	11
861	Blangkejeren (Blang Kejeren)	11
862	Blangpegayon (Blang Pegayon)	11
863	Blangpidie (Blang Pidie)	11
864	Blega	35
865	Blimbing	35
866	Blimbingsari	35
867	Blora (Blora kota)	33
868	Bluluk	35
869	Bluto	35
870	Boawae	53
871	Bobotsari	33
872	Bodeh	33
873	Bogabaida	94
874	Bogonuk	95
875	Bogor Barat	32
876	Bogor Selatan	32
877	Bogor Tengah	32
878	Bogor Timur	32
879	Bogor Utara	32
880	Bogorejo	33
881	Boja	33
882	Bojonegara	36
883	Bojonegoro	35
884	Bojong	36
885	Bojong	32
886	Bojong	33
887	Bojong Gede (Bojonggede)	32
888	Bojongasih	32
889	Bojonggambir	32
890	Bojonggenteng (Bojong Genteng)	32
891	Bojongloa Kaler	32
892	Bojongloa Kidul	32
893	Bojongmangu	32
894	Bojongmanik	36
895	Bojongpicung	32
896	Bojongsari	33
897	Bojongsari	32
898	Bojongsoang	32
899	Bokan Kepulauan	72
900	Bokat	72
901	Boking	53
902	Bokondini	95
903	Bokoneri	95
904	Bola	53
905	Bola	73
906	Bolaang	71
907	Bolaang Timur	71
908	Bolaang Uki	71
909	Bolakme	95
910	Bolangitang Barat (Bolang Itang Barat)	71
911	Bolangitang Timur (Bolang Itang Timur)	71
912	Bolano	72
913	Bolano Lambunu	72
914	Boleng	53
915	Boliyohuto (Boliohuto)	75
916	Bolo	52
917	Bomakia	93
918	Bomberay	92
919	Bomela	95
920	Bonai Darussalam	14
921	Bonang	33
922	Bonatua Lunasi	12
923	Bondifuar	91
924	Bondoala	74
925	Bondowoso	35
926	Bone	75
927	Bone (Bone Tondo)	74
928	Bone Bone	73
929	Bone Raya	75
930	Bonegunu	74
931	Bonehau	76
932	Bonepantai	75
933	Bongan	64
934	Bongas	32
935	Bonggakaradeng	73
936	Bonggo	91
937	Bonggo Timur	91
938	Bongomeme	75
939	Bonjol	13
940	Bonorowo	33
941	Bontang Barat	64
942	Bontang Selatan	64
943	Bontang Utara	64
944	Bonti	61
945	Bonto Bahari	73
946	Bonto Tiro (Bontotiro)	73
947	Bontoa (Maros Utara)	73
948	Bontoala	73
949	Bontocani	73
950	Bontoharu	73
951	Bontolempangang	73
952	Bontomanai	73
953	Bontomarannu	73
954	Bontomatene	73
955	Bontonompo	73
956	Bontonompo Selatan	73
957	Bontoramba	73
958	Bontosikuyu	73
959	Borbor	12
960	Borme	95
961	Borobudur	33
962	Boronadu	12
963	Borong	53
964	Bosar Maligas	12
965	Botain	92
966	Botin Leobele	53
967	Botolinggo	35
968	Botomuzoi	12
969	Botumoito (Botumoita)	75
970	Botupingge (Botu Pingge)	75
971	Bowobado	94
972	Boyan Tanjung	61
973	Boyolali	33
974	Boyolangu	35
975	Bpiri	95
976	Braja Selebah (Braja Slebah)	18
977	Bram Itam	15
978	Brang Ene	52
979	Brang Rea	52
980	Brangsong	33
981	Brati	33
982	Brebes	33
983	Bringin	33
984	Bringin	35
985	Brondong	35
986	Bruno	33
987	Bruwa	95
988	Bruyadori	91
989	BTS Ulu	16
990	Bua	73
991	Bua Ponrang (Bupon)	73
992	Buahbatu	32
993	Buahdua	32
994	Bualemo (Boalemo)	72
995	Buana Pemaca	16
996	Buaran	33
997	Buay Bahuga	18
998	Buay Madang	16
999	Buay Madang Timur	16
1000	Buay Pemaca	16
1001	Buay Pematang Ribu Ranau Tengah	16
1002	Buay Pemuka Bangsa Raja	16
1003	Buay Pemuka Peliung	16
1004	Buay Rawan	16
1005	Buay Runjung	16
1006	Buay Sandang Aji	16
1007	Buayan	33
1008	Bubon	11
1009	Bubulan	35
1010	Bubutan	35
1011	Budong-Budong	76
1012	Buduran	35
1013	Buer	52
1014	Bugi	95
1015	Buguk Gona	95
1016	Bugul Kidul	35
1017	Buk	92
1018	Bukal	72
1019	Bukateja	33
1020	Buke	74
1021	Buki	73
1022	Bukik Barisan	13
1023	Bukit	11
1024	Bukit Batu	62
1025	Bukit Batu	14
1026	Bukit Bestari	21
1027	Bukit Intan (Bukitintan)	19
1028	Bukit Kapur	14
1029	Bukit Kemuning	18
1030	Bukit Malintang	12
1031	Bukit Raya	62
1032	Bukit Raya	14
1033	Bukit Santuai (Bukit Santuei)	62
1034	Bukit Sundi	13
1035	Bukit Tusam	11
1036	Bukitkecil (Bukit Kecil)	16
1037	Bukitkerman	15
1038	Buko	72
1039	Buko Selatan	72
1040	Bula	81
1041	Bula Barat	81
1042	Bulagi	72
1043	Bulagi Selatan	72
1044	Bulagi Utara	72
1045	Bulak	35
1046	Bulakamba	33
1047	Bulang	21
1048	Bulango Selatan	75
1049	Bulango Timur	75
1050	Bulango Ulu	75
1051	Bulango Utara	75
1052	Bulawa	75
1053	Buleleng	51
1054	Bulik	62
1055	Bulik Timur	62
1056	Bulo	76
1057	Bulok	18
1058	Bulu	33
1059	Bulu Taba	76
1060	Bulukerto	33
1061	Bulukumpa (Bulukumba)	73
1062	Bululawang	35
1063	Bulupoddo	73
1064	Buluspesantren	33
1065	Bumi Agung	18
1066	Bumi Makmur	63
1067	Bumi Nabung	18
1068	Bumi Ratu Nuban	18
1069	Bumi Raya	72
1070	Bumi Waras	18
1071	Bumiaji	35
1072	Bumiayu	33
1073	Bumijawa	33
1074	Bunaken	71
1075	Bunaken Kepulauan	71
1076	Bunga Mas	17
1077	Bunga Mayang	16
1078	Bunga Mayang	18
1079	Bunga Raya	14
1080	Bungah	35
1081	Bungatan	35
1082	Bungaya	73
1083	Bungbulang	32
1084	Bungi	74
1085	Bungin	73
1086	Bungkal	35
1087	Bungku Barat	72
1088	Bungku Pesisir	72
1089	Bungku Selatan	72
1090	Bungku Tengah	72
1091	Bungku Timur	72
1092	Bungku Utara	72
1093	Bungo Dani	15
1094	Bungoro	73
1095	Bungur	63
1096	Bunguran Barat	21
1097	Bunguran Batubi	21
1098	Bunguran Selatan	21
1099	Bunguran Tengah	21
1100	Bunguran Timur	21
1101	Bunguran Timur Laut	21
1102	Bunguran Utara	21
1103	Bungursari	32
1104	Bungus Teluk Kabung	13
1105	Bunobogu	72
1106	Bunta	72
1107	Buntao	73
1108	Buntu Batu	73
1109	Buntu Pane	12
1110	Buntu Pepasan	73
1111	Buntulia	75
1112	Buntumalangka	76
1113	Bunut	14
1114	Bunut Hilir	61
1115	Bunut Hulu	61
1116	Bunyu (Pulau Bunyu)	65
1117	Burau	73
1118	Burneh	35
1119	Buru	21
1120	Buruway	92
1121	Busang	64
1122	Busungbiu (Busung biu)	51
1123	Butuh	33
1124	Buyasuri	53
1125	Cabangbungin	32
1126	Cadasari	36
1127	Cakranegara	52
1128	Cakung	31
1129	Camba	73
1130	Cambai	16
1131	Campaka	32
1132	Campakamulya (Campaka Mulya)	32
1133	Campalagian	76
1134	Camplong	35
1135	Campurdarat (Campur Darat)	35
1136	Candi	35
1137	Candi Laras Selatan	63
1138	Candi Laras Utara	63
1139	Candimulyo	33
1140	Candipuro	35
1141	Candipuro	18
1142	Candiroto	33
1143	Candisari	33
1144	Candung	13
1145	Cangkringan	34
1146	Cangkuang	32
1147	Cantigi	32
1148	Capkala	61
1149	Carenang (Cerenang)	36
1150	Caringin	32
1151	Carita	36
1152	Cariu	32
1153	Catubouw	92
1154	Cawas	33
1155	Celala	11
1156	Cempa	73
1157	Cempaga	62
1158	Cempaga Hulu	62
1159	Cempaka	63
1160	Cempaka	16
1161	Cempaka Putih	31
1162	Cendana	73
1163	Cengal	16
1164	Cengkareng	31
1165	Cenrana	73
1166	Ceper	33
1167	Cepiring	33
1168	Cepogo	33
1169	Cepu	33
1170	Cerbon	63
1171	Cerenti	14
1172	Cerme	35
1173	Cermee	35
1174	Cermin Nan Gedang / Gadang	15
1175	Ciambar	32
1176	Ciamis	32
1177	Ciampea	32
1178	Ciampel	32
1179	Cianjur	32
1180	Ciasem	32
1181	Ciater	32
1182	Ciawi	32
1183	Ciawigebang	32
1184	Cibadak	36
1185	Cibadak	32
1186	Cibal	53
1187	Cibal Barat	53
1188	Cibaliung	36
1189	Cibalong	32
1190	Cibarusah	32
1191	Cibatu	32
1192	Cibeber	36
1193	Cibeber	32
1194	Cibeunying Kaler	32
1195	Cibeunying Kidul	32
1196	Cibeureum	32
1197	Cibingbin	32
1198	Cibinong	32
1199	Cibiru	32
1200	Cibitung	36
1201	Cibitung	32
1202	Cibiuk	32
1203	Cibodas	36
1204	Cibogo	32
1205	Cibuaya	32
1206	Cibugel	32
1207	Cibungbulang	32
1208	Cicalengka	32
1209	Cicantayan	32
1210	Cicendo	32
1211	Cicurug	32
1212	Cidadap	32
1213	Cidahu	32
1214	Cidaun	32
1215	Cidolog	32
1216	Ciemas	32
1217	Cigalontang	32
1218	Cigandamekar	32
1219	Cigasong	32
1220	Cigedug	32
1221	Cigemlong (Cigemblong)	36
1222	Cigeulis	36
1223	Cigombong	32
1224	Cigudeg	32
1225	Cigugur	32
1226	Cihampelas	32
1227	Cihara	36
1228	Cihaurbeuti	32
1229	Cihideung	32
1230	Cihurip	32
1231	Cijaku	36
1232	Cijambe	32
1233	Cijati	32
1234	Cijeruk	32
1235	Cijeungjing	32
1236	Cijulang	32
1237	Cikadu	32
1238	Cikajang	32
1239	Cikakak	32
1240	Cikalong	32
1241	Cikalongkulon	32
1242	Cikalongwetan (Cikalong Wetan)	32
1243	Cikampek	32
1244	Cikancung	32
1245	Cikande	36
1246	Cikarang Barat	32
1247	Cikarang Pusat	32
1248	Cikarang Selatan	32
1249	Cikarang Timur	32
1250	Cikarang Utara	32
1251	Cikatomas	32
1252	Cikaum	32
1253	Cikedal (Cikeudal)	36
1254	Cikedung	32
1255	Cikelet	32
1256	Cikembar	32
1257	Cikeusal	36
1258	Cikeusik	36
1259	Cikidang	32
1260	Cikijing	32
1261	Cikole	32
1262	Cikoneng	32
1263	Cikulur	36
1264	Cikupa	36
1265	Cilacap Selatan	33
1266	Cilacap Tengah	33
1267	Cilacap Utara	33
1268	Cilaku	32
1269	Cilamaya Kulon	32
1270	Cilamaya Wetan	32
1271	Cilandak	31
1272	Cilawu	32
1273	Cilebak	32
1274	Cilebar	32
1275	Ciledug	36
1276	Ciledug	32
1277	Cilegon	36
1278	Cileles	36
1279	Cilengkrang	32
1280	Cileungsi	32
1281	Cileunyi	32
1282	Cililin	32
1283	Cilimus	32
1284	Cilincing	31
1285	Cilodong	32
1286	Cilograng	36
1287	Cilongok	33
1288	Cimahi	32
1289	Cimahi Selatan	32
1290	Cimahi Tengah	32
1291	Cimahi Utara	32
1292	Cimalaka	32
1293	Cimanggis	32
1294	Cimanggu	36
1295	Cimanggu	32
1296	Cimanggu	33
1297	Cimanggung	32
1298	Cimanuk	36
1299	Cimaragas	32
1300	Cimarga	36
1301	Cimaung	32
1302	Cimenyan (Cimeunyan)	32
1303	Cimerak	32
1304	Cina	73
1305	Cinambo	32
1306	Cinangka	36
1307	Cineam	32
1308	Cinere	32
1309	Cingambul	32
1310	Ciniru	32
1311	Cintapuri Darussalam	63
1312	Ciomas	36
1313	Ciomas	32
1314	Cipaku	32
1315	Cipanas	36
1316	Cipanas	32
1317	Ciparay	32
1318	Cipari	33
1319	Cipatat	32
1320	Cipatujah	32
1321	Cipayung	31
1322	Cipayung	32
1323	Cipedes	32
1324	Cipeucang	36
1325	Cipeundeuy	32
1326	Cipicung	32
1327	Cipocok Jaya	36
1328	Cipondoh	36
1329	Cipongkor	32
1330	Cipunagara	32
1331	Ciputat	36
1332	Ciputat Timur	36
1333	Ciracap	32
1334	Ciracas	31
1335	Ciranjang	32
1336	Cireunghas	32
1337	Cirinten	36
1338	Ciruas	36
1339	Cisaat	32
1340	Cisaga	32
1341	Cisalak	32
1342	Cisarua	32
1343	Cisata	36
1344	Cisauk	36
1345	Cisayong	32
1346	Ciseeng	32
1347	Cisewu	32
1348	Cisitu	32
1349	Cisoka	36
1350	Cisolok	32
1351	Cisompet	32
1352	Cisurupan	32
1353	Citak-Mitak (Citakmitak)	93
1354	Citamiang	32
1355	Citangkil	36
1356	Citeureup	32
1357	Citta	73
1358	Ciwandan	36
1359	Ciwaringin	32
1360	Ciwaru	32
1361	Ciwidey	32
1362	Cluring	35
1363	Cluwak	33
1364	Coblong	32
1365	Colomadu	33
1366	Comal	33
1367	Compreng	32
1368	Concong	14
1369	Conggeang	32
1370	Congkar	53
1371	Cot Girek	11
1372	Cugenang	32
1373	Cukuh Balak	18
1374	Culamega	32
1375	Curahdami	35
1376	Curio	73
1377	Curug	36
1378	Curug bitung (Curugbitung)	36
1379	Curugkembar	32
1380	Curup	17
1381	Curup Selatan	17
1382	Curup Tengah	17
1383	Curup Timur	17
1384	Curup Utara	17
1385	Dabun Gelang (Debun Gelang)	11
1386	Dadahup	62
1387	Dagai	94
1388	Dagangan	35
1389	Daha Barat	63
1390	Daha Selatan	63
1391	Daha Utara	63
1392	Dako Pemean (Dako Pamean)	72
1393	Dal	95
1394	Damai	64
1395	Damang Batu	62
1396	Damar	19
1397	Damau (Damao)	71
1398	Damer	81
1399	Dampal Selatan	72
1400	Dampal Utara	72
1401	Dampelas	72
1402	Dampit	35
1403	Danau Kembar	13
1404	Danau Kerinci	15
1405	Danau Kerinci Barat	15
1406	Danau Panggang	63
1407	Danau Paris	11
1408	Danau Seluluk	62
1409	Danau Sembuluh	62
1410	Danau Sipin	15
1411	Danau Teluk	15
1412	Dander	35
1413	Dangia	74
1414	Danime	95
1415	Danurejan	34
1416	Dapurang	76
1417	Darangdan	32
1418	Darma	32
1419	Darmaraja	32
1420	Darul Aman	11
1421	Darul Falah	11
1422	Darul Hasanah	11
1423	Darul Hikmah	11
1424	Darul Ihsan (Iksan)	11
1425	Darul Imarah	11
1426	Darul Kamal	11
1427	Darul Makmur	11
1428	Darussalam	11
1429	Dasuk	35
1430	Dataran Beimes	92
1431	Dataran Isim	92
1432	Datuk Bandar	12
1433	Datuk Bandar Timur	12
1434	Datuk Lima Puluh	12
1435	Datuk Tanah Datar	12
1436	Dau	35
1437	Dawan	51
1438	Dawarblandong (Dawar Blandong)	35
1439	Dawe	33
1440	Dawelor Dawera	81
1441	Dawuan	32
1442	Dayeuhkolot	32
1443	Dayeuhluhur	33
1444	Dayun	14
1445	Dedai	61
1446	Deiyai Miyo	94
1447	Dekai	95
1448	Deket	35
1449	Delang	62
1450	Delanggu	33
1451	Deleng Pokhkisen (Deleng Pokhisen)	11
1452	Deli Tua	12
1453	Delima	11
1454	Delta Pawan	61
1455	Demak	33
1456	Demba	91
1457	Demon Pagong	53
1458	Dempet	33
1459	Dempo Selatan	16
1460	Dempo Tengah	16
1461	Dempo Utara	16
1462	Demta	91
1463	Dendang	19
1464	Dendang	15
1465	Dende' Piongan Napo	73
1466	Dengilo	75
1467	Denpasar Barat	51
1468	Denpasar Selatan	51
1469	Denpasar Timur	51
1470	Denpasar Utara	51
1471	Dente Teladas	18
1472	Depapre	91
1473	Depati Tujuh	15
1474	Depok	34
1475	Depok	32
1476	Der Koumur	93
1477	Dervos	94
1478	Detukeli	53
1479	Detusoko	53
1480	Dewantara	11
1481	Didohu	92
1482	Dimba	95
1483	Dimembe	71
1484	Dipa	94
1485	Dirwemna	95
1486	Diwek	35
1487	Dlanggu	35
1488	Dlingo	34
1489	Dogiyai	94
1490	Dogomo	94
1491	Doko	35
1492	Dokome	94
1493	Dolat Rayat	12
1494	Dolo	72
1495	Dolo Barat	72
1496	Dolo Selatan	72
1497	Dolog Masagal	12
1498	Dolok	12
1499	Dolok Batu Nanggar	12
1500	Dolok Masihul	12
1501	Dolok Merawan	12
1502	Dolok Panribuan	12
1503	Dolok Pardamean	12
1504	Dolok Sanggul	12
1505	Dolok Sigompulon	12
1506	Dolok Silou (Dolok Silau)	12
1507	Dolopo	35
1508	Dompu	52
1509	Dondo	72
1510	Donggo	52
1511	Dongko	35
1512	Donomulyo	35
1513	Donorojo	35
1514	Donorojo	33
1515	Donri Donri	73
1516	Doreng	53
1517	Doro	33
1518	Doufo	94
1519	Dow	95
1520	Dramaga	32
1521	Dringu	35
1522	Driyorejo	35
1523	Dua Boccoe	73
1524	Dua Pitue	73
1525	Duampanua	73
1526	Duduksampeyan (Duduk Sampeyan)	35
1527	Duhiadaa	75
1528	Dukuh Pakis	35
1529	Dukuhseti	33
1530	Dukuhturi	33
1531	Dukuhwaru	33
1532	Dukun	33
1533	Dukun	35
1534	Dukupuntang	32
1535	Dulupi	75
1536	Dumadama	94
1537	Dumai Barat	14
1538	Dumai Kota	14
1539	Dumai Selatan	14
1540	Dumai Timur	14
1541	Dumbo Raya	75
1542	Dumoga	71
1543	Dumoga Barat	71
1544	Dumoga Tengah	71
1545	Dumoga Tenggara	71
1546	Dumoga Timur	71
1547	Dumoga Utara	71
1548	Dundu	95
1549	Dungaliyo	75
1550	Dungingi	75
1551	Dungkek	35
1552	Duo Koto	13
1553	Durai	21
1554	Duram	95
1555	Duren Sawit	31
1556	Durenan	35
1557	Duripoku	76
1558	Duruka	74
1559	Dusun Hilir	62
1560	Dusun Selatan	62
1561	Dusun Tengah	62
1562	Dusun Timur	62
1563	Dusun Utara	62
1564	Ebungfao (Ebungfau / Ebungfa)	91
1565	Edera	93
1566	Egiam	95
1567	Eipumek	95
1568	Ekadide	94
1569	Elar	53
1570	Elar Selatan	53
1571	Elelim	95
1572	Elikobal	93
1573	Ella Hilir	61
1574	Elpaputih	81
1575	Embaloh Hilir	61
1576	Embaloh Hulu	61
1577	Embetpen	95
1578	Empanang	61
1579	Empang	52
1580	Empat Petulai Dangku	16
1581	Enam Lingkung	13
1582	Ende	53
1583	Ende Selatan	53
1584	Ende Tengah	53
1585	Ende Timur	53
1586	Ende Utara	53
1587	Endomen	95
1588	Enggal	18
1589	Enggano	17
1590	Enok	14
1591	Enrekang	73
1592	Entikong	61
1593	Eragayam	95
1594	Erelmakawia	94
1595	Eremerasa	73
1596	Eris	71
1597	Eromoko	33
1598	Essang	71
1599	Essang Selatan	71
1600	Fafurwar (Irorutu)	92
1601	Fajar Timur	94
1602	Fak-Fak (Fakfak)	92
1603	Fak-Fak Barat (Fakfak Barat)	92
1604	Fak-Fak Tengah (Fakfak Tengah)	92
1605	Fak-Fak Timur (Fakfak Timur)	92
1606	Fakfak Timur Tengah	92
1607	Fanayama	12
1608	Fatukopa	53
1609	Fatuleu	53
1610	Fatuleu Barat	53
1611	Fatuleu Tengah	53
1612	Fatumnasi	53
1613	Fautmolo	53
1614	Fawi	94
1615	Fayit	93
1616	Fef	92
1617	Fena Fafan	81
1618	Fena Leisela	81
1619	Firiwage	93
1620	Fofi	93
1621	Fokour	92
1622	Fordata (Yaru)	81
1623	Furwagi	92
1624	Gabek	19
1625	Gabus	33
1626	Gabuswetan	32
1627	Gading	35
1628	Gading Cempaka	17
1629	Gading Rejo	18
1630	Gadingrejo	35
1631	Gadung	72
1632	Gajah	33
1633	Gajah Putih	11
1634	Gajahmungkur (Gajah Mungkur)	33
1635	Galang	72
1636	Galang	21
1637	Galang	12
1638	Galela	82
1639	Galela Barat	82
1640	Galela Selatan	82
1641	Galela Utara	82
1642	Galesong	73
1643	Galesong Selatan	73
1644	Galesong Utara	73
1645	Galing	61
1646	Galis	35
1647	Galur	34
1648	Gambir	31
1649	Gambiran	35
1650	Gambut	63
1651	Gamelia	95
1652	Gampengrejo	35
1653	Gamping	34
1654	Gandangbatu Sillanan (Gandang Batu Sillanan)	73
1655	Gandapura (Ganda Pura)	11
1656	Ganding	35
1657	Gandrungmangu	33
1658	Gandus	16
1659	Gandusari	35
1660	Gane Barat	82
1661	Gane Barat Selatan	82
1662	Gane Barat Utara	82
1663	Gane Timur	82
1664	Gane Timur Selatan	82
1665	Gane Timur Tengah	82
1666	Ganeas	32
1667	Gangga	52
1668	Ganra	73
1669	Gantar	32
1670	Gantarang (Gantorang, Gangking)	73
1671	Gantarang Keke (Gantareng Keke)	73
1672	Gantiwarno	33
1673	Gantung	19
1674	Gapura	35
1675	Garawangi	32
1676	Garoga	12
1677	Garum	35
1678	Garung	33
1679	Garut Kota	32
1680	Gatak	33
1681	Gaung	14
1682	Gaung Anak Serka	14
1683	Gayam	35
1684	Gayamsari	33
1685	Gayungan	35
1686	Gearek	95
1687	Gebang	12
1688	Gebang	32
1689	Gebang	33
1690	Gebog	33
1691	Gedangan	35
1692	Gedangsari (Gedang Sari)	34
1693	Gedebage	32
1694	Gedeg	35
1695	Gedong Tataan	18
1696	Gedongtengen (Gedong Tengen)	34
1697	Gedung Aji	18
1698	Gedung Aji Baru	18
1699	Gedung Meneng	18
1700	Gedung Surian	18
1701	Geger	35
1702	Gegerbitung (Geger Bitung)	32
1703	Gegesik	32
1704	Gekbrong	32
1705	Gelok Beam	95
1706	Gelumbang	16
1707	Gemarang	35
1708	Gemawang	33
1709	Gembong	33
1710	Gemeh	71
1711	Gemolong	33
1712	Gempol	32
1713	Gempol	35
1714	Gemuh	33
1715	Gending	35
1716	Geneng	35
1717	Genteng	35
1718	Gentuma Raya	75
1719	Genuk	33
1720	Geragai	15
1721	Gerih	35
1722	Gerogol	36
1723	Gerokgak	51
1724	Gerung	52
1725	Gerunggang	19
1726	Geselma	95
1727	Gesi	33
1728	Getasan	33
1729	Geumpang	11
1730	Geuredong Pase	11
1731	Geya	95
1732	Geyer	33
1733	Gianyar	51
1734	Gido	12
1735	Gika	95
1736	Giliginting (Gili Ginting)	35
1737	Gilireng	73
1738	Gilubandu	95
1739	Giri	35
1740	Giri Mulia (Giri Mulya)	17
1741	Girian	71
1742	Girimarto	33
1743	Girimaya	19
1744	Girimulyo	34
1745	Girisubo	34
1746	Giritontro	33
1747	Giriwoyo	33
1748	Girsang Sipangan Bolon	12
1749	Gisting	18
1750	Gladagsari	33
1751	Glagah	35
1752	Glenmore	35
1753	Glumpang Baro	11
1754	Glumpang Tiga (Geulumpang Tiga)	11
1755	Gn. Bintang Awai (Gunung Bintang Awai)	62
1756	Goa Balim	95
1757	Godean	34
1758	Godong	33
1759	Golewa	53
1760	Golewa Barat	53
1761	Golewa Selatan	53
1762	Gollo	95
1763	Gombong	33
1764	Gome	94
1765	Gome Utara	94
1766	Gomo	12
1767	Gondang	33
1768	Gondang	35
1769	Gondanglegi	35
1770	Gondangrejo	33
1771	Gondangwetan (Gondang Wetan)	35
1772	Gondokusuman	34
1773	Gondomanan	34
1774	Gorom Timur	81
1775	Goyage	95
1776	Grabag	33
1777	Grabagan	35
1778	Grati	35
1779	Greged (Greget)	32
1780	Gresi Selatan	91
1781	Gresik	35
1782	Gringsing	33
1783	Grobogan	33
1784	Grogol	33
1785	Grogol	35
1786	Grogol Petamburan	31
1787	Grong Grong	11
1788	Grujugan	35
1789	Gu	74
1790	Gubeng	35
1791	Gubug	33
1792	Gubume	94
1793	Gucialit	35
1794	Gudo	35
1795	Guguak (Gugu)	13
1796	Guguak Panjang (Guguk Panjang)	13
1797	Guluk-Guluk (Guluk Guluk)	35
1798	Gumay Talang	16
1799	Gumay Ulu	16
1800	Gumbasa	72
1801	Gumelar	33
1802	Gumukmas (Gumuk Mas)	35
1803	Guna	95
1804	Gundagi	95
1805	Gunem	33
1806	Guntur	33
1807	Gunuang Omeh (Gunung Mas)	13
1808	Gunung Agung	18
1809	Gunung Alip	18
1810	Gunung Anyar (Gununganyar)	35
1811	Gunung Jati (Cirebon Utara)	32
1812	Gunung Kaler	36
1813	Gunung Kerinci	15
1814	Gunung Kijang	21
1815	Gunung Labuhan	18
1816	Gunung Malela	12
1817	Gunung Maligas	12
1818	Gunung Megang	16
1819	Gunung Meriah	12
1820	Gunung Meriah (Mariah)	11
1821	Gunung Pelindung	18
1822	Gunung Purei	62
1823	Gunung Putri	32
1824	Gunung Puyuh (Gunungpuyuh)	32
1825	Gunung Raya	15
1826	Gunung Sahilan	14
1827	Gunung Sindur	32
1828	Gunung Sitember	12
1829	Gunung Sugih	18
1830	Gunung Tabur	64
1831	Gunung Talang	13
1832	Gunung Terang	18
1833	Gunung Timang	62
1834	Gunung Tujuh	15
1835	Gunung Tuleh (Gunungtuleh)	13
1836	Gunungguruh	32
1837	Gununghalu	32
1838	Gunungkencana (Gunung Kencana)	36
1839	Gunungpati	33
1840	Gunungsari	52
1841	Gunungsari (Gunung Sari)	36
1842	Gunungsitoli	12
1843	Gunungsitoli Alo'oa	12
1844	Gunungsitoli Barat	12
1845	Gunungsitoli Idanoi	12
1846	Gunungsitoli Selatan	12
1847	Gunungsitoli Utara	12
1848	Gunungtanjung (Gunung Tanjung)	32
1849	Gunungtoar (Gunung Toar)	14
1850	Gunungwungkal	33
1851	Gupura	95
1852	Gurage	94
1853	Gurah	35
1854	Habinsaran	12
1855	Haharu	53
1856	Haju	93
1857	Halong	63
1858	Halongonan	12
1859	Halongonan Timur	12
1860	Hampang	63
1861	Hamparan Perak	12
1862	Hamparan Rawang	15
1863	Hanau	62
1864	Hantakan	63
1865	Hantara	32
1866	Haranggaol Horisan (Haranggaol Horison)	12
1867	Harau	13
1868	Harian	12
1869	Harjamukti	32
1870	Haruai	63
1871	Haruyan	63
1872	Hatonduhan	12
1873	Hatungun	63
1874	Haur Gading	63
1875	Haurgeulis	32
1876	Haurwangi	32
1877	Hawu Mehara	53
1878	Helumo	71
1879	Heram	91
1880	Hereapini	95
1881	Herlang (Hero Lange Lange)	73
1882	Hewokloang	53
1883	Hibala	12
1884	Hiliduho	12
1885	Hilimegai	12
1886	Hilipuk	95
1887	Hiliran Gumanti	13
1888	Hilisalawa'ahe (Hilisalawaahe)	12
1889	Hiliserangkai (Hili Serangkai / Hilisaranggu)	12
1890	Hinai	12
1891	Hingk	92
1892	Hitadipa	94
1893	Hoat Sorbay	81
1894	Hobard	92
1895	Hogio	95
1896	Holuwon	95
1897	Homeyo	94
1898	Hoya	94
1899	Hu'u	52
1900	Huamual	81
1901	Huamual Belakang	81
1902	Hubikiak	95
1903	Hubikosi	95
1904	Hulonthalangi	75
1905	Hulu Gurung	61
1906	Hulu Kuantan	14
1907	Hulu Palik	17
1908	Hulu Sihapas	12
1909	Hulu Sungai	61
1910	Hulu Sungkai	18
1911	Huristak	12
1912	Huruna	12
1913	Huta Bargot	12
1914	Huta Bayu Raja	12
1915	Hutaraja Tinggi (Huta Raja Tinggi)	12
1916	Ibele	95
1917	Ibu	82
1918	Ibu Selatan	82
1919	Ibu Utara	82
1920	Ibun	32
1921	Idanogawo (Idano Gawo)	12
1922	Idanotae	12
1923	Idi Rayeuk	11
1924	Idi Timur	11
1925	Idi Tunong	11
1926	Ilaga	94
1927	Ilaga Utara	94
1928	Ilamburawi	94
1929	Ile Ape	53
1930	Ile Ape Timur	53
1931	Ile Boleng	53
1932	Ile Bura	53
1933	Ile Mandiri	53
1934	Ilir Barat Dua (Ilir Barat II)	16
1935	Ilir Barat Satu (Ilir Barat I)	16
1936	Ilir Talo	17
1937	Ilir Timur Dua (Ilir Timur II)	16
1938	Ilir Timur Satu (Ilir Timur I)	16
1939	Ilir Timur Tiga	16
1940	Ilu	94
1941	Ilugwa	95
1942	Ilwayab (Ilyawab)	93
1943	Imogiri	34
1944	Inamosol	81
1945	Inanwatan	92
1946	Indihiang	32
1947	Indra Jaya	11
1948	Indra Makmu (Indra Makmur)	11
1949	Indrajaya (Indra Jaya)	11
1950	Indralaya	16
1951	Indralaya Selatan	16
1952	Indralaya Utara	16
1953	Indramayu	32
1954	Indrapuri	11
1955	Inerie	53
1956	Inggerus	91
1957	Ingin Jaya	11
1958	Inikgal	95
1959	Iniyandit	93
1960	Iniye	95
1961	Insana	53
1962	Insana Barat	53
1963	Insana Fafinesu	53
1964	Insana Tengah	53
1965	Insana Utara	53
1966	Inuman	14
1967	Io Kufeu	53
1968	Ipuh (Muko Muko Selatan)	17
1969	Ireres	92
1970	Irimuli	94
1971	Itlay Hisage	95
1972	IV Jurai	13
1973	IV Koto (Ampek Koto)	13
1974	IV Koto Aur Malintang	13
1975	IV Nagari	13
1976	IV Nagari Bayang Utara	13
1977	Iwaka	94
1978	Iwoimendaa	74
1979	Iwur	95
1980	IX Koto Sungai Lasi	13
1981	Jabiren Raya	62
1982	Jabon	35
1983	Jabung	35
1984	Jabung	18
1985	Jagakarsa	31
1986	Jagebob	93
1987	Jagoi Babang	61
1988	Jagong Jeget	11
1989	Jailolo	82
1990	Jailolo Selatan	82
1991	Jair	93
1992	Jakabaring	16
1993	Jaken	33
1994	Jakenan	33
1995	Jalaksana	32
1996	Jalancagak	32
1997	Jamanis	32
1998	Jambangan	35
1999	Jambe	36
2000	Jambesari Darus Sholah	35
2001	Jambi Luar Kota	15
2002	Jambi Selatan	15
2003	Jambi Timur	15
2004	Jamblang	32
2005	Jambon	35
2006	Jambu	33
2007	Jampangkulon (Jampang Kulon)	32
2008	Jampangtengah (Jampang Tengah)	32
2009	Janapria	52
2010	Jangka	11
2011	Jangka Buya (Jangka Buaya)	11
2012	Jangkang	61
2013	Jangkar	35
2014	Jangkat	15
2015	Jangkat Timur (Sungai Tenang)	15
2016	Japah	33
2017	Japara	32
2018	Jarai	16
2019	Jaro	63
2020	Jasinga	32
2021	Jaten	33
2022	Jati	33
2023	Jati Agung	18
2024	Jatiasih	32
2025	Jatibanteng	35
2026	Jatibarang	32
2027	Jatibarang	33
2028	Jatigede	32
2029	Jatikalen	35
2030	Jatilawang	33
2031	Jatiluhur	32
2032	Jatinagara	32
2033	Jatinangor	32
2034	Jatinegara	31
2035	Jatinegara	33
2036	Jatinom	33
2037	Jatinunggal	32
2038	Jatipurno	33
2039	Jatipuro	33
2040	Jatirejo	35
2041	Jatirogo	35
2042	Jatiroto	35
2043	Jatiroto	33
2044	Jatisampurna (Jati Sampurna)	32
2045	Jatisari	32
2046	Jatisrono	33
2047	Jatitujuh	32
2048	Jatiuwung	36
2049	Jatiwangi	32
2050	Jatiwaras	32
2051	Jatiyoso	33
2052	Jawa Maraja Bah Jambi	12
2053	Jawai	61
2054	Jawai Selatan	61
2055	Jawilan	36
2056	Jaya	11
2057	Jaya Baru	11
2058	Jayakerta	32
2059	Jayaloka (Jaya Loka)	16
2060	Jayanti	36
2061	Jayapura	16
2062	Jayapura Selatan	91
2063	Jayapura Utara	91
2064	Jebres	33
2065	Jebus	19
2066	Jejangkit	63
2067	Jejawi	16
2068	Jekan Raya	62
2069	Jekulo	33
2070	Jelai	62
2071	Jelai Hulu	61
2072	Jelbuk	35
2073	Jelimpo	61
2074	Jelutung	15
2075	Jemaja	21
2076	Jemaja Barat	21
2077	Jemaja Timur	21
2078	Jembrana	51
2079	Jempang	64
2080	Jenamas	62
2081	Jenangan	35
2082	Jenar	33
2083	Jenawi	33
2084	Jenggawah	35
2085	Jenu	35
2086	Jepara	33
2087	Jepon	33
2088	Jerebuu	53
2089	Jereweh	52
2090	Jerowaru	52
2091	Jeruklegi	33
2092	Jetfa	95
2093	Jetis	35
2094	Jetis	34
2095	Jetsy	93
2096	Jeumpa	11
2097	Jeunieb	11
2098	Jiken	33
2099	Jila	94
2100	Jiput	36
2101	Jirak Jaya	16
2102	Jita	94
2103	Jiwan	35
2104	Joerat	93
2105	Jogonalan	33
2106	Jogorogo	35
2107	Jogoroto	35
2108	Johan Pahwalan (Johan Pahlawan)	11
2109	Johar Baru	31
2110	Jombang	36
2111	Jombang	35
2112	Jonggat	52
2113	Jonggol	32
2114	Jongkat (Siantan)	61
2115	Jongkong (Jengkong)	61
2116	Jorlang Hataran	12
2117	Jorong	63
2118	Joutu	93
2119	Jrengik	35
2120	Juai	63
2121	Juhar	12
2122	Jujuhan	15
2123	Jujuhan Ilir	15
2124	Juli	11
2125	Julok	11
2126	Jumantono	33
2127	Jumapolo	33
2128	Jumo	33
2129	Junjung Sirih	13
2130	Junrejo	35
2131	Juntinyuat	32
2132	Juwana	33
2133	Juwangi	33
2134	Juwiring	33
2135	Kabaena	74
2136	Kabaena Barat	74
2137	Kabaena Selatan	74
2138	Kabaena Tengah	74
2139	Kabaena Timur	74
2140	Kabaena Utara	74
2141	Kabandungan	32
2142	Kabangka	74
2143	Kabanjahe	12
2144	Kabaruan	71
2145	Kabat	35
2146	Kabawo	74
2147	Kabianggama	95
2148	Kabila	75
2149	Kabila Bone	75
2150	Kabola	53
2151	Kabuh	35
2152	Kabun	14
2153	Kadatua	74
2154	Kademangan	35
2155	Kadia	74
2156	Kadipaten	32
2157	Kadudampit	32
2158	Kadugede	32
2159	Kaduhejo	36
2160	Kadungora	32
2161	Kadupandak	32
2162	Kadur	35
2163	Kahaungu Eti (Kahaunguweti)	53
2164	Kahayan Hilir	62
2165	Kahayan Hulu Utara	62
2166	Kahayan Kuala	62
2167	Kahayan Tengah	62
2168	Kahu	73
2169	Kai	95
2170	Kaibar	93
2171	Kaidipang	71
2172	Kaimana	92
2173	Kairatu	81
2174	Kairatu Barat	81
2175	Kais	92
2176	Kais Darat	92
2177	Kaisenar	91
2178	Kaitaro	92
2179	Kajang	73
2180	Kajen	33
2181	Kajoran	33
2182	Kajuara	73
2183	Kakas	71
2184	Kakas Barat	71
2185	Kakuluk Mesak	53
2186	Kalaena	73
2187	Kalanganyar	36
2188	Kalapanunggal (Kalapa Nunggal)	32
2189	Kalasan	34
2190	Kalawat	71
2191	Kaledupa	74
2192	Kaledupa Selatan	74
2193	Kalianda	18
2194	Kalianget	35
2195	Kaliangkrik	33
2196	Kalibagor	33
2197	Kalibaru	35
2198	Kalibawang	34
2199	Kalibawang	33
2200	Kalibening	33
2201	Kalibunder	32
2202	Kalidawir	35
2203	Kalideres	31
2204	Kalidoni	16
2205	Kaligesing	33
2206	Kaligondang	33
2207	Kalijambe	33
2208	Kalijati	32
2209	Kalikajar	33
2210	Kalikotes	33
2211	Kalimanah	33
2212	Kalimanggis	32
2213	Kalinyamatan	33
2214	Kaliorang	64
2215	Kaliori	33
2216	Kalipare	35
2217	Kalipucang	32
2218	Kalipuro	35
2219	Kalirejo	18
2220	Kalis	61
2221	Kalisat	35
2222	Kalitengah	35
2223	Kalitidu	35
2224	Kaliwates	35
2225	Kaliwedi	32
2226	Kaliwiro	33
2227	Kaliwungu	33
2228	Kaliwungu Selatan	33
2229	Kalomdol	95
2230	Kalome	94
2231	Kalongan	71
2232	Kaloran	33
2233	Kalukku	76
2234	Kalumpang	76
2235	Kalumpang (Kelumpang)	63
2236	Kamal	35
2237	Kamang Baru	13
2238	Kamang Magek	13
2239	Kamanre	73
2240	Kambata Mapambuhang	53
2241	Kambera	53
2242	Kamboneri	95
2243	Kambowa	74
2244	Kambrau (Kambraw / Kamberau)	92
2245	Kambu	74
2246	Kamipang	62
2247	Kampa (Kampar Timur)	14
2248	Kampak	35
2249	Kampar	14
2250	Kampar Kiri	14
2251	Kampar Kiri Hilir	14
2252	Kampar Kiri Hulu	14
2253	Kampar Kiri Tengah	14
2254	Kampar Utara	14
2255	Kampung Laut	33
2256	Kampung Melayu	17
2257	Kampung Rakyat	12
2258	Kamu	94
2259	Kamu Selatan	94
2260	Kamu Timur	94
2261	Kamu Utara	94
2262	Kamundan	92
2263	Kanatang	53
2264	Kandangan	63
2265	Kandangan	33
2266	Kandangan	35
2267	Kandanghaur	32
2268	Kandangserang	33
2269	Kandat	35
2270	Kandeman	33
2271	Kandis	14
2272	Kandis	16
2273	Kangae	53
2274	Kangayan	35
2275	Kanggime	95
2276	Kangkung	33
2277	Kanigaran	35
2278	Kanigoro	35
2279	Kanor	35
2280	Kao	82
2281	Kao Barat	82
2282	Kao Teluk	82
2283	Kao Utara	82
2284	Kapala Pitu (Kapalla Pitu)	73
2285	Kapas	35
2286	Kapetakan	32
2287	Kapiraya	94
2288	Kapoiala	74
2289	Kapongan	35
2290	Kapontori	74
2291	Kaptel	93
2292	Kapuas (Sanggau Kapuas)	61
2293	Kapuas Barat	62
2294	Kapuas Hilir	62
2295	Kapuas Hulu	62
2296	Kapuas Kuala	62
2297	Kapuas Murung	62
2298	Kapuas Tengah	62
2299	Kapuas Timur	62
2300	Kapur IX/Sembilan	13
2301	Karamat	72
2302	Karang Agung Ilir	16
2303	Karang Bahagia (Karangbahagia)	32
2304	Karang Baru	11
2305	Karang Bintang	63
2306	Karang Dapo	16
2307	Karang Intan	63
2308	Karang Jaya	16
2309	Karang Kancana (Karangkancana)	32
2310	Karang Pilang (Karangpilang)	35
2311	Karang Tanjung	36
2312	Karang Tengah	36
2313	Karang Tinggi	17
2314	Karangampel	32
2315	Karangan	35
2316	Karangan	64
2317	Karanganom	33
2318	Karanganyar	33
2319	Karanganyar	35
2320	Karangasem (Karang Asem)	51
2321	Karangawen	33
2322	Karangbinangun	35
2323	Karangdadap	33
2324	Karangdowo	33
2325	Karanggayam	33
2326	Karanggede	33
2327	Karanggeneng (Karang Geneng)	35
2328	Karangjambu	33
2329	Karangjati	35
2330	Karangjaya (Karang Jaya)	32
2331	Karangkobar	33
2332	Karanglewas	33
2333	Karangmalang	33
2334	Karangmojo	34
2335	Karangmoncol	33
2336	Karangnongko	33
2337	Karangnunggal	32
2338	Karangpandan	33
2339	Karangpawitan	32
2340	Karangpenang (Karang Penang)	35
2341	Karangploso (Karang Ploso)	35
2342	Karangpucung	33
2343	Karangrayung	33
2344	Karangreja	33
2345	Karangrejo	35
2346	Karangrejo (Karang Rejo)	35
2347	Karangsambung	33
2348	Karangsembung	32
2349	Karangtengah	32
2350	Karangtengah	33
2351	Karangtengah (Karang Tengah)	33
2352	Karangwareng	32
2353	Karas	35
2354	Karas	92
2355	Karau Kuala	62
2356	Karawaci	36
2357	Karawang Barat	32
2358	Karawang Timur	32
2359	Kare	35
2360	Karera	53
2361	Karimun	21
2362	Karimunjawa (Karimun Jawa)	33
2363	Karossa	76
2364	Kartasura	33
2365	Kartoharjo	35
2366	Kartoharjo (Kertoharjo)	35
2367	Karu	95
2368	Karubaga	95
2369	Karusen Janang	62
2370	Karya Penggawa	18
2371	Kasembon	35
2372	Kasemen	36
2373	Kasi	92
2374	Kasihan	34
2375	Kasiman	35
2376	Kasimbar	72
2377	Kasiruta Barat	82
2378	Kasiruta Timur	82
2379	Kasokandel	32
2380	Kasomalang	32
2381	Kasreman	35
2382	Kasui	18
2383	Katala Hamu Lingu	53
2384	Katang Bidare	21
2385	Katapang	32
2386	Kateman	14
2387	Katibung	18
2388	Katiku Tana	53
2389	Katiku Tana Selatan (Katikutana Selatan)	53
2390	Katingan Hilir	62
2391	Katingan Hulu	62
2392	Katingan Kuala	62
2393	Katingan Tengah	62
2394	Katobu	74
2395	Katoi	74
2396	Kaubun	64
2397	Kauditan	71
2398	Kauman	35
2399	Kaur Selatan	17
2400	Kaur Tengah	17
2401	Kaur Utara	17
2402	Kaureh	91
2403	Kawagit	93
2404	Kawali	32
2405	Kawalu	32
2406	Kawangkoan	71
2407	Kawangkoan Barat	71
2408	Kawangkoan Utara	71
2409	Kaway XVI	11
2410	Kawedanan	35
2411	Kawor	95
2412	Kawunganten	33
2413	Kayan Hilir	61
2414	Kayan Hilir	65
2415	Kayan Hulu	61
2416	Kayan Hulu	65
2417	Kayan Selatan	65
2418	Kayangan	52
2419	Kayauni	92
2420	Kayen	33
2421	Kayen Kidul	35
2422	Kayo	95
2423	Kayoa	82
2424	Kayoa Barat	82
2425	Kayoa Selatan	82
2426	Kayoa Utara	82
2427	Kayu Agung	16
2428	Kayu Aro	15
2429	Kayu Aro Barat	15
2430	Kebakkramat	33
2431	Kebar	92
2432	Kebar Selatan	92
2433	Kebar Timur	92
2434	Kebasen	33
2435	Kebawetan	17
2436	Kebayakan	11
2437	Kebayoran Baru	31
2438	Kebayoran Lama	31
2439	Kebo	94
2440	Kebomas	35
2441	Kebon Jeruk	31
2442	Kebonagung	33
2443	Kebonagung (Kebon Agung)	35
2444	Kebonarum	33
2445	Kebonpedes	32
2446	Kebonsari (Kebon Sari)	35
2447	Kebumen	33
2448	Kebun Tebu	18
2449	Kedamaian	18
2450	Kedamean	35
2451	Kedaton	18
2452	Kedaton Peninjauan Raya	16
2453	Kedawung	32
2454	Kedawung	33
2455	Kedewan	35
2456	Kediri	51
2457	Kediri	52
2458	Kedokan Bunder	32
2459	Kedondong	18
2460	Kedopok (Kedopak)	35
2461	Kedu	33
2462	Kedung	33
2463	Kedung Waringin	32
2464	Kedungadem	35
2465	Kedungbanteng (Kedung Banteng)	33
2466	Kedungdung	35
2467	Kedunggalar	35
2468	Kedungjajang	35
2469	Kedungjati	33
2470	Kedungkandang	35
2471	Kedungpring	35
2472	Kedungreja	33
2473	Kedungtuban	33
2474	Kedungwaru	35
2475	Kedungwuni	33
2476	Kedurang	17
2477	Kedurang Ilir	17
2478	Keera	73
2479	Kegayem	95
2480	Kei Besar	81
2481	Kei Besar Selatan	81
2482	Kei Besar Selatan Barat	81
2483	Kei Besar Utara Barat	81
2484	Kei Besar Utara Timur	81
2485	Kei Kecil	81
2486	Kei Kecil Barat	81
2487	Kei Kecil Timur	81
2488	Kei Kecil Timur Selatan	81
2489	Kejajar	33
2490	Kejaksan	32
2491	Kejayan	35
2492	Kejobong	33
2493	Kejuruan Muda	11
2494	Kelam Permai	61
2495	Kelam Tengah	17
2496	Kelapa	19
2497	Kelapa Dua	36
2498	Kelapa Gading	31
2499	Kelapa Kampit	19
2500	Kelapa Lima	53
2501	Kelara	73
2502	Kelay	64
2503	Kelayang	14
2504	Kelekar	16
2505	Kelila	95
2506	Keliling Danau	15
2507	Kelimutu	53
2508	Keling	33
2509	Kelua (Klua)	63
2510	Keluang	16
2511	Kelubagolit	53
2512	Kelulome	95
2513	Kelumbayan (Klumbayan)	18
2514	Kelumbayan Barat (Klumbayan Barat)	18
2515	Kelumpang Barat	63
2516	Kelumpang Hilir	63
2517	Kelumpang Hulu	63
2518	Kelumpang Selatan	63
2519	Kelumpang Tengah	63
2520	Kelumpang Utara	63
2521	Kema	71
2522	Kemalang	33
2523	Kemang	32
2524	Kemangkon	33
2525	Kemayoran	31
2526	Kembang	33
2527	Kembang Janggut	64
2528	Kembang Tanjong	11
2529	Kembangan	31
2530	Kembangbahu	35
2531	Kembaran	33
2532	Kembayan	61
2533	Kembru	94
2534	Kembu	95
2535	Kemiling	18
2536	Kemiri	36
2537	Kemiri	33
2538	Kemlagi	35
2539	Kempas	14
2540	Kempo	52
2541	Kemranjen	33
2542	Kemtuk	91
2543	Kemtuk Gresi	91
2544	Kemuning	14
2545	Kemuning	16
2546	Kemusu	33
2547	Kencong	35
2548	Kendahe	71
2549	Kendal	35
2550	Kendal	33
2551	Kendari	74
2552	Kendari Barat	74
2553	Kendawangan	61
2554	Kendit	35
2555	Kenduruan	35
2556	Kenjeran	35
2557	Kenohan	64
2558	Kenyam	95
2559	Keo Tengah	53
2560	Kep. Bala Balakang	76
2561	Kepahiang	17
2562	Kepala Madan	81
2563	Kepanjen	35
2564	Kepanjenkidul (Kepanjen Kidul)	35
2565	Kepenuhan	14
2566	Kepenuhan Hulu	14
2567	Kepil	33
2568	Kepohbaru	35
2569	Kepulauan Ambai	91
2570	Kepulauan Aruri	91
2571	Kepulauan Ayau	92
2572	Kepulauan Botanglomang	82
2573	Kepulauan Joronga	82
2574	Kepulauan Karimata	61
2575	Kepulauan Manipa	81
2576	Kepulauan Marore	71
2577	Kepulauan Masaloka Raya	74
2578	Kepulauan Pongok	19
2579	Kepulauan Posek	21
2580	Kepulauan Roma (Romang)	81
2581	Kepulauan Sangkarrang	73
2582	Kepulauan Sembilan	92
2583	Kepulauan Seribu Selatan	31
2584	Kepulauan Seribu Utara	31
2585	Kepulauan Tanakeke	73
2586	Kepung	35
2587	Kerajaan	12
2588	Kerambitan	51
2589	Kerek	35
2590	Kerinci Kanan	14
2591	Keritang	14
2592	Kerjo	33
2593	Kerkap	17
2594	Kersamanah	32
2595	Kersana	33
2596	Kertajati	32
2597	Kertak Hanyar	63
2598	Kertanegara	33
2599	Kertapati	16
2600	Kertasari	32
2601	Kertasemaya	32
2602	Kertek	33
2603	Kertosono	35
2604	Keruak	52
2605	Kerumutan	14
2606	Kesamben	35
2607	Kesambi	32
2608	Kesesi	33
2609	Kesu	73
2610	Kesugihan	33
2611	Ketahun	17
2612	Ketambe	11
2613	Ketanggungan	33
2614	Ketapang	35
2615	Ketapang	18
2616	Ketol	11
2617	Ketungau Hilir	61
2618	Ketungau Hulu	61
2619	Ketungau Tengah	61
2620	Keumala	11
2621	Kewapante	53
2622	Ki	93
2623	KI'E (Kie, Ki'e)	53
2624	Kian Darat	81
2625	Kiaracondong	32
2626	Kiarapedes	32
2627	Kibin	36
2628	Kikim Barat	16
2629	Kikim Selatan	16
2630	Kikim Tengah	16
2631	Kikim Timur	16
2632	Kilmid	95
2633	Kilmury	81
2634	Kilo	52
2635	Kimaam	93
2636	Kinal	17
2637	Kinali	13
2638	Kindang	73
2639	Kinovaro	72
2640	Kintamani	51
2641	Kintap	63
2642	Kintom	72
2643	Kirihi	91
2644	Kisam Ilir	16
2645	Kisam Tinggi	16
2646	Kisar Selatan (Pulau Pulau Terselatan)	81
2647	Kisar Utara	81
2648	Kismantoro	33
2649	Kiwirok	95
2650	Kiwirok Timur	95
2651	Kiyage	94
2652	Klabang	35
2653	Klabot	92
2654	Klakah	35
2655	Klambu	33
2656	Klamono	92
2657	Klampis	35
2658	Klangenan	32
2659	Klapanunggal	32
2660	Klari	32
2661	Klasafet	92
2662	Klaso	92
2663	Klaten Selatan	33
2664	Klaten Tengah	33
2665	Klaten Utara	33
2666	Klaurung	92
2667	Klawak	92
2668	Klayili	92
2669	Kledung	33
2670	Klego	33
2671	Klirong	33
2672	Klojen	35
2673	Kluet Selatan	11
2674	Kluet Tengah	11
2675	Kluet Timur	11
2676	Kluet Utara	11
2677	Klungkung	51
2678	Koba	19
2679	Kobakma	95
2680	Kobalima	53
2681	Kobalima Timur	53
2682	Kodeoha	74
2683	Kodi	53
2684	Kodi Balaghar	53
2685	Kodi Bangedo	53
2686	Kodi Utara	53
2687	Kofiau	92
2688	Koja	31
2689	Kok Baun	53
2690	Kokalukuna	74
2691	Kokap	34
2692	Kokas	92
2693	Kokoda	92
2694	Kokoda Utara	92
2695	Kokop	35
2696	Kolaka	74
2697	Kolang	12
2698	Kolawa	95
2699	Kolbano	53
2700	Kolf Braza	93
2701	Kolono	74
2702	Kolono Timur	74
2703	Kombay	93
2704	Kombeng (Kongbeng)	64
2705	Kombi	71
2706	Kombut	93
2707	Komodo	53
2708	Kona	95
2709	Konang	35
2710	Konawe	74
2711	Konda	74
2712	Konda	92
2713	Konda/ Kondaga	95
2714	Konhir	92
2715	Kontu Kowuna	74
2716	Kontuar	93
2717	Kontunaga	74
2718	Kopang	52
2719	Kopay	93
2720	Kopo	36
2721	Kora	95
2722	Koragi	95
2723	Kormomolin	81
2724	Koroncong	36
2725	Koroptak	95
2726	Koroway Buluanop	93
2727	Korupun	95
2728	Kosambi	36
2729	Kosarek	95
2730	Kosiwo	91
2731	Kot Olin	53
2732	Kota (Kediri Kota)	35
2733	Kota Agung	16
2734	Kota Agung (Kota Agung Pusat)	18
2735	Kota Agung Barat	18
2736	Kota Agung Timur	18
2737	Kota Arga Makmur	17
2738	Kota Atambua (Atambua Kota)	53
2739	Kota Bahagia	11
2740	Kota Baharu	11
2741	Kota Bangun	64
2742	Kota Bangun Darat	64
2743	Kota Barat	75
2744	Kota Baru	15
2745	Kota Baru	53
2746	Kota Baru (Kotabaru)	32
2747	Kota Besi	62
2748	Kota Gajah	18
2749	Kota Jantho	11
2750	Kota Juang	11
2751	Kota Kefamenanu	53
2752	Kota Kisaran Barat	12
2753	Kota Kisaran Timur	12
2754	Kota Komba	53
2755	Kota Komba Utara	53
2756	Kota Kualasinpang (Kota Kuala Simpang)	11
2757	Kota Kudus (Kudus Kota)	33
2758	Kota Lama	53
2759	Kota Maba	82
2760	Kota Manna	17
2761	Kota Masohi	81
2762	Kota Mukomuko (Mukomuko Utara)	17
2763	Kota Padang	17
2764	Kota Raja	53
2765	Kota Selatan	75
2766	Kota Sigli	11
2767	Kota Soe	53
2768	Kota Sumenep	35
2769	Kota Tambolaka	53
2770	Kota Tengah	75
2771	Kota Ternate Selatan	82
2772	Kota Ternate Tengah	82
2773	Kota Ternate Utara	82
2774	Kota Timur	75
2775	Kota Utara	75
2776	Kota Waikabubak	53
2777	Kota Waingapu	53
2778	Kota Waisai	92
2779	Kotaanyar (Kota Anyar)	35
2780	Kotabumi	18
2781	Kotabumi Selatan	18
2782	Kotabumi Utara	18
2783	Kotabunan	71
2784	Kotagede	34
2785	Kotamobagu Barat	71
2786	Kotamobagu Selatan	71
2787	Kotamobagu Timur	71
2788	Kotamobagu Utara	71
2789	Kotanopan	12
2790	Kotapinang (Kota Pinang)	12
2791	Kotarih	12
2792	Kotawaringin Lama	62
2793	Koting	53
2794	Koto Balingka	13
2795	Koto Baru	13
2796	Koto Baru	15
2797	Koto Besar	13
2798	Koto Gasib	14
2799	Koto Kampar Hulu	14
2800	Koto Parik Gadang Diateh	13
2801	Koto Salak	13
2802	Koto Tangah	13
2803	Koto VII	13
2804	Koto XI Tarusan	13
2805	Kouh	93
2806	Kradenan	33
2807	Kragan	33
2808	Kragilan	36
2809	Kraksaan	35
2810	Kramat	33
2811	Kramatjati (Kramat Jati)	31
2812	Kramatmulya (Kramat Mulya)	32
2813	Kramatwatu	36
2814	Kramongmongga	92
2815	Kranggan	35
2816	Kranggan	33
2817	Krangkeng	32
2818	Kras	35
2819	Kraton	34
2820	Kraton	35
2821	Krayan	65
2822	Krayan Barat	65
2823	Krayan Selatan	65
2824	Krayan Tengah	65
2825	Krayan Timur	65
2826	Krejengan	35
2827	Krembangan	35
2828	Krembung	35
2829	Krepkuri	95
2830	Kresek	36
2831	Kretek	34
2832	Krian	35
2833	Kromengan	35
2834	Kronjo	36
2835	Kroya	32
2836	Kroya	33
2837	Krucil	35
2838	Krueng Barona Jaya	11
2839	Krueng Sabee	11
2840	Krui Selatan	18
2841	Kuala	12
2842	Kuala	11
2843	Kuala Baru	11
2844	Kuala Batee	11
2845	Kuala Behe	61
2846	Kuala Betara	15
2847	Kuala Cenaku	14
2848	Kuala Indragiri	14
2849	Kuala Jambi	15
2850	Kuala Kampar	14
2851	Kuala Kencana	94
2852	Kuala Mandor B	61
2853	Kuala Pesisir	11
2854	Kualin	53
2855	Kualuh Hilir	12
2856	Kualuh Hulu	12
2857	Kualuh Leidong	12
2858	Kualuh Selatan	12
2859	Kuanfatu	53
2860	Kuantan Hilir	14
2861	Kuantan Hilir Seberang	14
2862	Kuantan Mudik	14
2863	Kuantan Tengah	14
2864	Kuari	95
2865	Kuaro	64
2866	Kuatnana	53
2867	Kubu	95
2868	Kubu	14
2869	Kubu	51
2870	Kubu	61
2871	Kubu Babussalam	14
2872	Kubung	13
2873	Kubutambahan	51
2874	Kudu	35
2875	Kulawi	72
2876	Kulawi Selatan	72
2877	Kulim	14
2878	Kulisusu (Kalingsusu/Kalisusu)	74
2879	Kulisusu Barat	74
2880	Kulisusu Utara	74
2881	Kulo	73
2882	Kuly Lanny	95
2883	Kumai	62
2884	Kumelembuai	71
2885	Kumpeh	15
2886	Kumpeh Ulu	15
2887	Kumun Debai	15
2888	Kundur	21
2889	Kundur Barat	21
2890	Kundur Utara	21
2891	Kunduran	33
2892	Kuningan	32
2893	Kunir	35
2894	Kunjang	35
2895	Kunto Darussalam	14
2896	Kuok	14
2897	Kupang Barat	53
2898	Kupang Tengah	53
2899	Kupang Timur	53
2900	Kupitan	13
2901	Kur Selatan	81
2902	Kuranji	63
2903	Kuranji	13
2904	Kurau	63
2905	Kuri	92
2906	Kuri Wamesa	92
2907	Kurik	93
2908	Kurima	95
2909	Kuripan	35
2910	Kuripan	63
2911	Kuripan	52
2912	Kurra	73
2913	Kurulu	95
2914	Kurun	62
2915	Kusambi	74
2916	Kusan Hilir	63
2917	Kusan Hulu	63
2918	Kusan Tengah	63
2919	Kuta	51
2920	Kuta Alam	11
2921	Kuta Baro	11
2922	Kuta Blang	11
2923	Kuta Cot Glie (Kota Cot Glie)	11
2924	Kuta Makmur	11
2925	Kuta Malaka (Kota Malaka)	11
2926	Kuta Raja	11
2927	Kuta Selatan	51
2928	Kuta Utara	51
2929	Kutabuluh (Kuta Buluh)	12
2930	Kutalimbaru	12
2931	Kutambaru	12
2932	Kutapanjang (Kuta Panjang)	11
2933	Kutasari	33
2934	Kutawaluya	32
2935	Kutawaringin	32
2936	Kute Panang	11
2937	Kute Siantan	21
2938	Kutoarjo	33
2939	Kutorejo	35
2940	Kutowinangun	33
2941	Kuwarasan	33
2942	Kuwus	53
2943	Kuwus Barat	53
2944	Kuyawage	95
2945	Kwadungan	35
2946	Kwamki Narama	94
2947	Kwandang	75
2948	Kwanyar	35
2949	Kwelamdua	95
2950	Kwesefo	92
2951	Kwikma	95
2952	Kwoor	92
2953	Labakkang	73
2954	Labang	35
2955	Labangka	52
2956	Labobo	72
2957	Laboya Barat (Lamboya Barat)	53
2958	Labuan	72
2959	Labuan	36
2960	Labuan Amas Selatan	63
2961	Labuan Amas Utara	63
2962	Labuapi	52
2963	Labuhan Badas	52
2964	Labuhan Deli	12
2965	Labuhan Haji	11
2966	Labuhan Haji	52
2967	Labuhan Haji Barat	11
2968	Labuhan Haji Timur	11
2969	Labuhan Maringgai	18
2970	Labuhan Ratu	18
2971	Ladongi	74
2972	Lae Parira	12
2973	Laenmanen	53
2974	Laeya	74
2975	Lage	72
2976	Laguboti	12
2977	Laham	64
2978	Lahat	16
2979	Lahat Selatan	16
2980	Lahei	62
2981	Lahei Barat	62
2982	Lahewa	12
2983	Lahewa Timur	12
2984	Lahomi (Gahori)	12
2985	Lahusa	12
2986	Laikang	73
2987	Lainea	74
2988	Lais	16
2989	Lais	17
2990	Lakarsantri	35
2991	Lakbok	32
2992	Lakea (Lipunoto)	72
2993	Lakudo	74
2994	Lalabata	73
2995	Lalan	16
2996	Lalembuu	74
2997	Lalolae	74
2998	Lalonggasumeeto	74
2999	Lamaknen	53
3000	Lamaknen Selatan	53
3001	Lamala	72
3002	Lamandau	62
3003	Lamasi	73
3004	Lamasi Timur	73
3005	Lamba Leda	53
3006	Lamba Leda Selatan (Poco Ranaka)	53
3007	Lamba Leda Timur (Poco Ranaka Timur)	53
3008	Lamba Leda Utara	53
3009	Lambai	74
3010	Lambandia	74
3011	Lambewi	94
3012	Lambitu	52
3013	Lamboya	53
3014	Lambu	52
3015	Lambu Kibang	18
3016	Lambuya	74
3017	Lamongan	35
3018	Lampasio	72
3019	Lampihong	63
3020	Lamposi Tigo Nagori / Nagari	13
3021	Lamuru	73
3022	Landasan Ulin	63
3023	Landawe	74
3024	Landono	74
3025	Landu Leko	53
3026	Langda	95
3027	Langensari	32
3028	Langgam	14
3029	Langgikima	74
3030	Langgudu	52
3031	Langkahan	11
3032	Langkaplancar	32
3033	Langkapura	18
3034	Langke Rembong	53
3035	Langowan Barat	71
3036	Langowan Selatan	71
3037	Langowan Timur	71
3038	Langowan Utara	71
3039	Langsa Barat	11
3040	Langsa Baro	11
3041	Langsa Kota	11
3042	Langsa Lama	11
3043	Langsa Timur	11
3044	Lannyna	95
3045	Lansirang (Lanrisang)	73
3046	Lantari Jaya	74
3047	Lantung	52
3048	Laonti	74
3049	Lapandewa	74
3050	Lapang	11
3051	Lape (Lape Lopok)	52
3052	Lappariaja	73
3053	Larangan	36
3054	Larangan	33
3055	Larangan	35
3056	Larantuka	53
3057	Lareh Sago Halaban	13
3058	Laren	35
3059	Lariang	76
3060	Larompong	73
3061	Larompong Selatan	73
3062	Lasalepa	74
3063	Lasalimu	74
3064	Lasalimu Selatan	74
3065	Lasem	33
3066	Lasiolat	53
3067	Lasolo	74
3068	Lasolo Kepulauan	74
3069	Lasusua	74
3070	Latambaga	74
3071	Latimojong	73
3072	Latoma	74
3073	Lau	73
3074	Laubaleng	12
3075	Laung Tuhup	62
3076	Laut Tador	12
3077	Laut Tawar (Lut Tawar)	11
3078	Lawa	74
3079	Lawang	35
3080	Lawang Kidul	16
3081	Lawang Wetan	16
3082	Lawe Alas	11
3083	Lawe Bulan	11
3084	Lawe Sigala Gala	11
3085	Lawe Sumur	11
3086	Laweyan	33
3087	Lea-Lea	74
3088	Lebak Wangi	36
3089	Lebakbarang	33
3090	Lebakgedong	36
3091	Lebaksiu	33
3092	Lebakwangi	32
3093	Lebatukan	53
3094	Lebong Atas	17
3095	Lebong Sakti	17
3096	Lebong Selatan	17
3097	Lebong Tengah	17
3098	Lebong Utara	17
3099	Leces	35
3100	Lede	82
3101	Ledo	61
3102	Ledokombo	35
3103	Legok	36
3104	Legonkulon	32
3105	Leihitu	81
3106	Leihitu Barat	81
3107	Leitimur Selatan	81
3108	Lekok	35
3109	Leksono	33
3110	Leksula	81
3111	Lela	53
3112	Lelak	53
3113	Lelea	32
3114	Leles	32
3115	Lemahabang	32
3116	Lemahsugih	32
3117	Lemahwungkuk (Lemah Wungkuk)	32
3118	Lembah Bawang	61
3119	Lembah Gumanti	13
3120	Lembah Masurai	15
3121	Lembah Melintang	13
3122	Lembah Sabil	11
3123	Lembah Segar	13
3124	Lembah Seulawah	11
3125	Lembah Sorik Marapi	12
3126	Lembak	16
3127	Lembang	73
3128	Lembang	32
3129	Lembang Jaya	13
3130	Lembar	52
3131	Lembean Timur	71
3132	Lembeh Selatan (Bitung Selatan)	71
3133	Lembeh Utara	71
3134	Lembeyan	35
3135	Lembo	72
3136	Lembo	74
3137	Lembo Raya	72
3138	Lembor	53
3139	Lembor Selatan	53
3140	Lembur	53
3141	Lembursitu	32
3142	Lemito	75
3143	Lemong	18
3144	Lempuing	16
3145	Lempuing Jaya	16
3146	Lenangguar	52
3147	Lendah	34
3148	Lenek	52
3149	Lengayang	13
3150	Lengkiti	16
3151	Lengkong	35
3152	Lengkong	32
3153	Lenteng	35
3154	Lepar (Lepar Pongok)	19
3155	Lepembusu Kelisoke	53
3156	Leupung	11
3157	Leuser	11
3158	Leuwidamar	36
3159	Leuwigoong	32
3160	Leuwiliang	32
3161	Leuwimunding	32
3162	Leuwisadeng	32
3163	Leuwisari	32
3164	Lewa	53
3165	Lewa Tidahu	53
3166	Lewolema	53
3167	Lhoknga (Lho'nga)	11
3168	Lhoksukon	11
3169	Lhoong	11
3170	Li Anogomma	95
3171	Liang	72
3172	Liang Anggang	63
3173	Libarek	95
3174	Libureng	73
3175	Licin	35
3176	Ligung	32
3177	Likupang Barat	71
3178	Likupang Selatan	71
3179	Likupang Timur	71
3180	Lilialy	81
3181	Liliraja (Lili Riaja)	73
3182	Lilirilau (Lili Rilau)	73
3183	Lima Kaum	13
3184	Lima Puluh	14
3185	Lima Puluh (Limapuluh)	12
3186	Lima Puluh Pesisir	12
3187	Limau	18
3188	Limbangan	33
3189	Limboro	76
3190	Limboto	75
3191	Limboto Barat	75
3192	Limbur Lubuk Mengkuang	15
3193	Limo	32
3194	Limpasu	63
3195	Limpung	33
3196	Limun	15
3197	Lindu	72
3198	Linge	11
3199	Lingga	21
3200	Lingga Bayu	12
3201	Lingga Timur	21
3202	Lingga Utara	21
3203	Linggang Bigung	64
3204	Linggo Sari Baganti	13
3205	Lingsar	52
3206	Lintang Kanan	16
3207	Lintau Buo	13
3208	Lintau Buo Utara	13
3209	Lintong Nihuta	12
3210	Lio Timur	53
3211	Lirik	14
3212	Lirung	71
3213	Liukang Kalmas (Kalukuang Masalima)	73
3214	Liukang Tangaya	73
3215	Liukang Tupabbiring	73
3216	Liukang Tupabbiring Utara	73
3217	Loa Janan	64
3218	Loa Janan Ilir	64
3219	Loa Kulu	64
3220	Loaholu	53
3221	Loano	33
3222	Lobalain	53
3223	Lobu	72
3224	Loceret	35
3225	Loea	74
3226	Logas Tanah Darat	14
3227	Lohbener	32
3228	Lohia	74
3229	Lokpaikat	63
3230	Loksado	63
3231	Lolak	71
3232	Lolat	95
3233	Lolayan	71
3234	Loli	53
3235	Loloda	82
3236	Loloda Kepulauan	82
3237	Loloda Tengah	82
3238	Loloda Utara	82
3239	Lolofitu Moi	12
3240	Lolomatua	12
3241	Lolong Guba	81
3242	Lolowau	12
3243	Long Apari	64
3244	Long Bagun	64
3245	Long Hubung	64
3246	Long Ikis	64
3247	Long Iram	64
3248	Long Kali	64
3249	Long Mesangat	64
3250	Long Pahangai	64
3251	Longkib	11
3252	Lopok	52
3253	Lore Barat	72
3254	Lore Piore	72
3255	Lore Selatan	72
3256	Lore Tengah	72
3257	Lore Timur	72
3258	Lore Utara	72
3259	Losarang	32
3260	Losari	32
3261	Losari	33
3262	Lotu	12
3263	Loura	53
3264	Lowokwaru	35
3265	Luahagundre Maniamolo	12
3266	Luak (Luhak)	13
3267	Luas	17
3268	Lubai	16
3269	Lubai Ulu	16
3270	Lubuak Tarok	13
3271	Lubuk Alung	13
3272	Lubuk Baja	21
3273	Lubuk Barumun	12
3274	Lubuk Basung	13
3275	Lubuk Batang	16
3276	Lubuk Batu Jaya	14
3277	Lubuk Begalung	13
3278	Lubuk Besar	19
3279	Lubuk Dalam	14
3280	Lubuk Keliat	16
3281	Lubuk Kilangan	13
3282	Lubuk Linggau Barat Dua (II)	16
3283	Lubuk Linggau Barat Satu (I)	16
3284	Lubuk Linggau Selatan Dua (II)	16
3285	Lubuk Linggau Selatan Satu (I)	16
3286	Lubuk Linggau Timur Dua (II)	16
3287	Lubuk Linggau Timur Satu (I)	16
3288	Lubuk Linggau Utara Dua (II)	16
3289	Lubuk Linggau Utara Satu (I)	16
3290	Lubuk Pakam	12
3291	Lubuk Pinang	17
3292	Lubuk Raja	16
3293	Lubuk Sandi	17
3294	Lubuk Sikaping	13
3295	Lubuk Sikarah	13
3296	Lueng Bata	11
3297	Luhak Nan Duo	13
3298	Lumajang	35
3299	Lumar	61
3300	Lumban Julu	12
3301	Lumbang	35
3302	Lumbir	33
3303	Lumbis	65
3304	Lumbis Hulu	65
3305	Lumbis Ogong	65
3306	Lumbis Pansiangan	65
3307	Lumbok Seminung	18
3308	Lumbung	32
3309	Lumo	94
3310	Lumut	12
3311	Lunang	13
3312	Lungkang Kule	17
3313	Lunyuk	52
3314	Luragung	32
3315	Luwuk	72
3316	Luwuk Selatan	72
3317	Luwuk Timur	72
3318	Luwuk Utara	72
3319	Luyo	76
3320	Ma'u	12
3321	Maba	82
3322	Maba Selatan	82
3323	Maba Tengah	82
3324	Maba Utara	82
3325	Mabugi	94
3326	Macang Pacar	53
3327	Madang Suku I	16
3328	Madang Suku II	16
3329	Madang Suku III	16
3330	Madapangga	52
3331	Madat	11
3332	Madidir (Bitung Tengah)	71
3333	Madiun	35
3334	Madukara	33
3335	Maduran	35
3336	Maesa	71
3337	Maesaan	71
3338	Maesan	35
3339	Mage'abume	94
3340	Magelang Selatan	33
3341	Magelang Tengah	33
3342	Magelang Utara	33
3343	Magepanda	53
3344	Magersari	35
3345	Magetan	35
3346	Maginti	74
3347	Mahu	53
3348	Maima	95
3349	Maiwa	73
3350	Maja	36
3351	Maja	32
3352	Majalaya	32
3353	Majalengka	32
3354	Majasari	36
3355	Majauleng	73
3356	Maje	17
3357	Majenang	33
3358	Makale	73
3359	Makale Selatan	73
3360	Makale Utara	73
3361	Makarti Jaya	16
3362	Makasar	31
3363	Makassar	73
3364	Makbon	92
3365	Makian Barat	82
3366	Makimi	94
3367	Makki	95
3368	Makmur	11
3369	Malabotom	92
3370	Maladum Mes	92
3371	Malaimsimsa	92
3372	Malaka Barat	53
3373	Malaka Tengah	53
3374	Malaka Timur	53
3375	Malalak (Malakak)	13
3376	Malalayang	71
3377	Malangbong	32
3378	Malangke	73
3379	Malangke Barat	73
3380	Malausma	32
3381	Maleber	32
3382	Malifut	82
3383	Maligano	74
3384	Maliku	62
3385	Malili	73
3386	Malimbong Balepe	73
3387	Malin Deman	17
3388	Malinau Barat	65
3389	Malinau Kota	65
3390	Malinau Selatan	65
3391	Malinau Selatan Hilir	65
3392	Malinau Selatan Hulu	65
3393	Malinau Utara	65
3394	Malind	93
3395	Malingping	36
3396	Malllawa (Mallawa)	73
3397	Mallusetasi	73
3398	Malo	35
3399	Malua	73
3400	Maluk	52
3401	Malunda	76
3402	Mam	95
3403	Mamajang	73
3404	Mamasa	76
3405	Mamberamo Hilir	91
3406	Mamberamo Hulu	91
3407	Mamberamo Tengah	91
3408	Mamberamo Tengah Timur	91
3409	Mambi	76
3410	Mambioman Bapai	93
3411	Mamboro	53
3412	Mamosalato	72
3413	Mampang Prapatan	31
3414	Mamuju	76
3415	Mananggu	75
3416	Mancak	36
3417	Mandah	14
3418	Mandai	73
3419	Mandalajati	32
3420	Mandalawangi	36
3421	Mandalle	73
3422	Mandastana	63
3423	Mandau	14
3424	Mandau Talawang	62
3425	Mande	32
3426	Mandiangin	15
3427	Mandiangin Koto Selayan	13
3428	Mandiangin Timur	15
3429	Manding	35
3430	Mandioli Selatan	82
3431	Mandioli Utara	82
3432	Mandiraja	33
3433	Mandirancan	32
3434	Mandobo	93
3435	Mandolang	71
3436	Mandonga	74
3437	Mandor	61
3438	Mandrehe	12
3439	Mandrehe Barat	12
3440	Mandrehe Utara	12
3441	Manduamas	12
3442	Mane	11
3443	Manekar	92
3444	Manganitu	71
3445	Manganitu Selatan	71
3446	Mangarabombang (Mangara Bombang)	73
3447	Mangaran	35
3448	Manggala	73
3449	Manggalewa	52
3450	Manggar	19
3451	Manggelum	93
3452	Manggeng	11
3453	Manggis	51
3454	Mangkubumi	32
3455	Mangkutana	73
3456	Mangoli Barat	82
3457	Mangoli Selatan	82
3458	Mangoli Tengah	82
3459	Mangoli Timur	82
3460	Mangoli Utara	82
3461	Mangoli Utara Timur	82
3462	Manguharjo	35
3463	Mangunjaya	32
3464	Mangunreja	32
3465	Maniamolo	12
3466	Maniangpajo	73
3467	Maniis	32
3468	Manimeri	92
3469	Manis Mata	61
3470	Manisrenggo	33
3471	Manna	17
3472	Mannem	91
3473	Manokwari Barat	92
3474	Manokwari Selatan	92
3475	Manokwari Timur	92
3476	Manokwari Utara	92
3477	Manonjaya	32
3478	Mantang	21
3479	Mantangai	62
3480	Mantewe	63
3481	Mantikulore	72
3482	Mantingan	35
3483	Mantoh	72
3484	Mantrijeron	34
3485	Mantup	35
3486	Manuhing	62
3487	Manuhing Raya (Mahuning Raya)	62
3488	Manuju	73
3489	Manyak Payed	11
3490	Manyar	35
3491	Manyaran	33
3492	Manyeuw	81
3493	Maos	33
3494	Maospati	35
3495	Mapanget	71
3496	Mapat Tunggul	13
3497	Mapat Tunggul Selatan	13
3498	Mapenduma	95
3499	Mapia	94
3500	Mapia Barat	94
3501	Mapia Tengah	94
3502	Mapilli	76
3503	Mapitara	53
3504	Mappak	73
3505	Mappakasunggu	73
3506	Mappedeceng	73
3507	Marabahan	63
3508	Marancar	12
3509	Marang (Ma Rang)	73
3510	Marang Kayu	64
3511	Maratua	64
3512	Marau	61
3513	Marawola	72
3514	Marawola Barat	72
3515	Marbau	12
3516	Mardingding (Mardinding)	12
3517	Mare	92
3518	Mare	73
3519	Mare Selatan	92
3520	Marga	51
3521	Marga Punduh	18
3522	Marga Sakti Sebelat (Marga Sakti)	17
3523	Marga Sekampung	18
3524	Marga Tiga (Margatiga)	18
3525	Margaasih	32
3526	Margadana	33
3527	Margahayu	32
3528	Margasari	33
3529	Margo Tabir	15
3530	Margomulyo	35
3531	Margorejo	33
3532	Margoyoso	33
3533	Mariat	92
3534	Marikit	62
3535	Marioriawa (Mario Riawa)	73
3536	Marioriwawo (Mario Riwawo)	73
3537	Marisa	75
3538	Mariso	73
3539	Maritengngae	73
3540	Maro Sebo	15
3541	Maro Sebo Ilir	15
3542	Maro Sebo Ulu	15
3543	Marobo	74
3544	Maron	35
3545	Maronge	52
3546	Maros Baru	73
3547	Marpoyan Damai	14
3548	Martapura	16
3549	Martapura (Martapura Kota)	63
3550	Martapura Barat	63
3551	Martapura Timur	63
3552	Marusu	73
3553	Masalembu	35
3554	Masalle	73
3555	Masama	72
3556	Masamba	73
3557	Masanda	73
3558	Masaran	33
3559	Masbagik	52
3560	Masirei	91
3561	Masni	92
3562	Masyeta	92
3563	Mata Oleo	74
3564	Mata Usu	74
3565	Matakali	76
3566	Matan Hilir Selatan	61
3567	Matan Hilir Utara	61
3568	Matangkuli	11
3569	Matangnga	76
3570	Mataram	52
3571	Mataram Baru	18
3572	Mataraman	63
3573	Mataru	53
3574	Matawai La Pawu (Lappau)	53
3575	Matemani	92
3576	Matesih	33
3577	Matraman	31
3578	Mattiro Bulu	73
3579	Mattiro Sompe (Matirro Sompe)	73
3580	Matuari (Bitung Barat)	71
3581	Matur	13
3582	Maudus	92
3583	Mauk	36
3584	Maukaro	53
3585	Maulafa	53
3586	Mauponggo	53
3587	Maurole	53
3588	Mawabuan	92
3589	Mawasangka	74
3590	Mawasangka Tengah	74
3591	Mawasangka Timur	74
3592	Mayamuk	92
3593	Mayang	35
3594	Mayangan	35
3595	Mayong	33
3596	Mazino	12
3597	Mazo	12
3598	Mbahamdandara	92
3599	Mbeliling	53
3600	Mbua Tengah	95
3601	Mbulmu Yalma	95
3602	Mbuwa	95
3603	Mdona Hyera (Mndona Hiera)	81
3604	Mebarok	95
3605	Medan Amplas	12
3606	Medan Area	12
3607	Medan Barat	12
3608	Medan Baru	12
3609	Medan Belawan (Medan Belawan Kota)	12
3610	Medan Deli	12
3611	Medan Denai	12
3612	Medan Helvetia	12
3613	Medan Johor	12
3614	Medan Kota	12
3615	Medan Labuhan	12
3616	Medan Maimun	12
3617	Medan Marelan	12
3618	Medan Perjuangan	12
3619	Medan Petisah	12
3620	Medan Polonia	12
3621	Medan Selayang	12
3622	Medan Sunggal	12
3623	Medan Tembung	12
3624	Medan Timur	12
3625	Medan Tuntungan	12
3626	Medang Deras	12
3627	Medang Kampai	14
3628	Medansatria (Medan Satria)	32
3629	Megaluh	35
3630	Megambilis	95
3631	Megamendung	32
3632	Megang Sakti	16
3633	Mego	53
3634	Mehalaan	76
3635	Mejayan	35
3636	Mejobo	33
3637	Mekakau Ilir	16
3638	Mekar Baru	36
3639	Mekarjaya	36
3640	Mekarmukti	32
3641	Mekarsari (Mekar Sari)	63
3642	Melagi	95
3643	Melagineri	95
3644	Melak	64
3645	Melaya	51
3646	Meliau	61
3647	Melinting	18
3648	Melonguane	71
3649	Melonguane Timur	71
3650	Meluhu	74
3651	Membalong	19
3652	Membey	92
3653	Mempawah Hilir	61
3654	Mempawah Hulu	61
3655	Mempawah Timur	61
3656	Mempura	14
3657	Mendahara	15
3658	Mendahara Ulu	15
3659	Mendawai	62
3660	Mendo Barat	19
3661	Mendoyo	51
3662	Menes	36
3663	Menganti	35
3664	Menggala	18
3665	Menggala Timur	18
3666	Mengkendek	73
3667	Mengwi	51
3668	Menjalin	61
3669	Menou	94
3670	Mentarang	65
3671	Mentarang Hulu	65
3672	Mentawa Baru Ketapang	62
3673	Mentaya Hilir Selatan	62
3674	Mentaya Hilir Utara	62
3675	Mentaya Hulu	62
3676	Mentebah	61
3677	Menteng	31
3678	Menthobi Raya	62
3679	Mentok (Muntok)	19
3680	Menui Kepulauan	72
3681	Menukung	61
3682	Menyuke	61
3683	Meos Mansar	92
3684	Mepanga	72
3685	Meraksa Aji	18
3686	Merakurak	35
3687	Meral	21
3688	Meral Barat	21
3689	Meranti	12
3690	Meranti	61
3691	Merapi Barat	16
3692	Merapi Selatan	16
3693	Merapi Timur	16
3694	Merauke	93
3695	Merawang	19
3696	Merbau	14
3697	Merbau Mataram	18
3698	Merdeka	12
3699	Merdey	92
3700	Merek	12
3701	Mergangsan	34
3702	Merigi	17
3703	Merigi Kelindang	17
3704	Merigi Sakti	17
3705	Merlung	15
3706	Mersam	15
3707	Mertoyudan	33
3708	Mesidah	11
3709	Mesjid Raya	11
3710	Messawa	76
3711	Mestong	15
3712	Mesuji	16
3713	Mesuji	18
3714	Mesuji Makmur	16
3715	Mesuji Raya	16
3716	Mesuji Timur	18
3717	Metro Barat	18
3718	Metro Kibang	18
3719	Metro Pusat	18
3720	Metro Selatan	18
3721	Metro Timur	18
3722	Metro Utara	18
3723	Meukek	11
3724	Meurah Dua	11
3725	Meurah Mulia	11
3726	Meuraxa	11
3727	Meureubo	11
3728	Meureudu	11
3729	Mewoluk	94
3730	Meyado (Mayado)	92
3731	Miangas	71
3732	Midai	21
3733	Mihing Raya	62
3734	Mijen	33
3735	Mila	11
3736	Milimbo	95
3737	Mimika Barat	94
3738	Mimika Barat Jauh	94
3739	Mimika Barat Tengah	94
3740	Mimika Baru	94
3741	Mimika Tengah	94
3742	Mimika Timur	94
3743	Mimika Timur Jauh	94
3744	Minas	14
3745	Minasa Tene	73
3746	Mindiptana	93
3747	Minggir	34
3748	Minyambaouw	92
3749	Minyamur	93
3750	Miomaffo Barat (Miomafo Barat)	53
3751	Miomaffo Tengah (Miomafo Tengah)	53
3752	Miomaffo Timur (Miomafo Timur)	53
3753	Miri	33
3754	Miri Manasa	62
3755	Mirit	33
3756	Misool (Misool Utara)	92
3757	Misool Barat	92
3758	Misool Selatan	92
3759	Misool Timur	92
3760	Miyah	92
3761	Miyah Selatan	92
3762	Mlandingan	35
3763	Mlarak	35
3764	Mlati	34
3765	Mlonggo	33
3766	Moa (Moa Lakor)	81
3767	Moba	95
3768	Modayag	71
3769	Modayag Barat	71
3770	Modo	35
3771	Modoinding	71
3772	Modung	35
3773	Mofinop	95
3774	Moga	33
3775	Moilong	72
3776	Moisegen	92
3777	Mojo	35
3778	Mojoagung	35
3779	Mojoanyar	35
3780	Mojogedang	33
3781	Mojolaban	33
3782	Mojoroto	35
3783	Mojosari	35
3784	Mojosongo	33
3785	Mojotengah	33
3786	Mojowarno	35
3787	Mokoni	95
3788	Molagalome	95
3789	Molanikime	94
3790	Molawe	74
3791	Mollo Barat	53
3792	Mollo Selatan	53
3793	Mollo Tengah	53
3794	Mollo Utara	53
3795	Molu Maru	81
3796	Momi Waren	92
3797	Momunu	72
3798	Monano	75
3799	Moncongloe (Moncong Loe)	73
3800	Mondokan	33
3801	Monta	52
3802	Montallat (Montalat)	62
3803	Montasik (Mantasiek)	11
3804	Monterado	61
3805	Montong	35
3806	Montong Gading	52
3807	Mooat	71
3808	Mook Manaar Bulatn	64
3809	Moora	94
3810	Mootilango	75
3811	Moraid	92
3812	Moramo	74
3813	Moramo Utara	74
3814	Mori Atas	72
3815	Mori Utara	72
3816	Moro	21
3817	Moro'o	12
3818	Moronge	71
3819	Morosi	74
3820	Morotai Jaya	82
3821	Morotai Selatan	82
3822	Morotai Selatan Barat	82
3823	Morotai Timur	82
3824	Morotai Utara	82
3825	Moskona Barat	92
3826	Moskona Selatan	92
3827	Moskona Timur	92
3828	Moskona Utara	92
3829	Moswaren	92
3830	Moti	82
3831	Motoling	71
3832	Motoling Barat	71
3833	Motoling Timur	71
3834	Motongkad	71
3835	Motui	74
3836	Moutong	72
3837	Mowewe	74
3838	Mowila	74
3839	Moyo Hilir	52
3840	Moyo Hulu	52
3841	Moyo Utara	52
3842	Moyudan	34
3843	Mpunda	52
3844	Mpur	92
3845	Mranggen	33
3846	Mrebet	33
3847	Muara	12
3848	Muara	95
3849	Muara	94
3850	Muara Ancalong	64
3851	Muara Badak	64
3852	Muara Bangka Hulu	17
3853	Muara Batang Gadis	12
3854	Muara Batang Toru	12
3855	Muara Batu	11
3856	Muara Belida	16
3857	Muara Beliti	16
3858	Muara Bengkal	64
3859	Muara Bulian	15
3860	Muara Dua	11
3861	Muara Dua (Muaradua)	16
3862	Muara Dua Kisam (Muaradua Kisam)	16
3863	Muara Enim	16
3864	Muara Harus	63
3865	Muara Jawa	64
3866	Muara Jaya	16
3867	Muara Kaman	64
3868	Muara Kelingi	16
3869	Muara Kemumu	17
3870	Muara Komam	64
3871	Muara Kuang	16
3872	Muara Lakitan	16
3873	Muara Lawa	64
3874	Muara Muntai	64
3875	Muara Padang	16
3876	Muara Pahu	64
3877	Muara Papalik	15
3878	Muara Pawan	61
3879	Muara Pinang	16
3880	Muara Sabak Barat	15
3881	Muara Sabak Timur	15
3882	Muara Sahung	17
3883	Muara Samu	64
3884	Muara Satu	11
3885	Muara Siau	15
3886	Muara Sipongi	12
3887	Muara Sugihan	16
3888	Muara Sungkai	18
3889	Muara Tabir	15
3890	Muara Tami	91
3891	Muara Telang	16
3892	Muara Tembesi	15
3893	Muara Tiga	11
3894	Muara Uya	63
3895	Muara Wahau	64
3896	Muara Wis	64
3897	Muaragembong (Muara Gembong)	32
3898	Muarapayang	16
3899	Mubrani	92
3900	Mugi	95
3901	Muko-muko Bathin VII	15
3902	Mukok	61
3903	Mulak Sebingkai	16
3904	Mulak Ulu	16
3905	Mulia	94
3906	Muliama	95
3907	Mulyorejo	35
3908	Mumbulsari	35
3909	Muncang	36
3910	Muncar	35
3911	Mundu	32
3912	Mungka	13
3913	Mungkajang	73
3914	Mungkid	33
3915	Munjul	36
3916	Munjungan	35
3917	Munte	12
3918	Muntilan	33
3919	Murhum	74
3920	Murkim	95
3921	Muruk Rian	65
3922	Murung	62
3923	Murung Pudak	63
3924	Musaik	95
3925	Musatfak	95
3926	Musi	53
3927	Mustikajaya (Mustika Jaya)	32
3928	Musuk	33
3929	Mutiara	11
3930	Mutiara Timur	11
3931	Muting	93
3932	Mutis	53
3933	Muye	94
3934	Na IX - X (Na IX-X)	12
3935	Nabire	94
3936	Nabire Barat	94
3937	Nabunage	95
3938	Naga Juang	12
3939	Naga Wutung	53
3940	Nagrak	32
3941	Nagreg	32
3942	Naibenu	53
3943	Naikere	92
3944	Nainggolan	12
3945	Nakama	94
3946	Nalca	95
3947	Nalo Tantan (Nalo Tatan)	15
3948	Nalumsari	33
3949	Naman Teran (Nama Teran)	12
3950	Namang	19
3951	Nambluong	91
3952	Nambo	72
3953	Nambo	74
3954	Namlea	81
3955	Namo Rambe (Namorambe)	12
3956	Namohalu Esiwa	12
3957	Namrole	81
3958	Nan Sabaris	13
3959	Nanaet Duabesi	53
3960	Nanga Mahap	61
3961	Nanga Pinoh	61
3962	Nanga Taman	61
3963	Nanga Tayap	61
3964	Nangapanda	53
3965	Nangaroro	53
3966	Nanggala	73
3967	Nanggalo	13
3968	Nanggulan	34
3969	Nanggung	32
3970	Nanusa	71
3971	Napabalano	74
3972	Napal Putih	17
3973	Napan	94
3974	Napano Kusambi	74
3975	Napua	95
3976	Naringgul	32
3977	Narmada	52
3978	Nasal	17
3979	Nassau	12
3980	Natal	12
3981	Natar	18
3982	Naukenjerai	93
3983	Nawangan	35
3984	Ndao Nuse	53
3985	Ndona	53
3986	Ndona Timur	53
3987	Ndori	53
3988	Ndoso	53
3989	Negara	51
3990	Negara Batin	18
3991	Negeri Agung	18
3992	Negeri Besar	18
3993	Negeri Katon	18
3994	Neglasari	36
3995	Nekamese	53
3996	Nelawi	95
3997	Nelle (Maumerei)	53
3998	Neney	92
3999	Nenggeagin	95
4000	Ngabang	61
4001	Ngablak	33
4002	Ngadiluwih	35
4003	Ngadirejo	33
4004	Ngadirojo	33
4005	Ngadirojo	35
4006	Ngadu Ngala	53
4007	Ngaglik	34
4008	Ngajum (Ngajung)	35
4009	Ngaliyan	33
4010	Ngambon	35
4011	Ngambur	18
4012	Ngampel	33
4013	Ngampilan	34
4014	Ngamprah	32
4015	Ngancar	35
4016	Nganjuk	35
4017	Ngantang	35
4018	Ngantru	35
4019	Ngapa	74
4020	Ngaras (Bengkunat Belimbing)	18
4021	Ngargoyoso	33
4022	Ngariboyo	35
4023	Ngaringan	33
4024	Ngasem	35
4025	Ngawen	34
4026	Ngawen	33
4027	Ngawi	35
4028	Ngebel	35
4029	Ngemplak	34
4030	Ngemplak	33
4031	Ngetos	35
4032	Nggaha Ori Angu (Nggaha Oriangu)	53
4033	Ngguti	93
4034	Ngimbang	35
4035	Nglegok	35
4036	Nglipar	34
4037	Ngluwar	33
4038	Ngluyu	35
4039	Ngombol	33
4040	Ngoro	35
4041	Ngraho	35
4042	Ngrambe	35
4043	Ngrampal	33
4044	Ngrayun	35
4045	Ngronggot	35
4046	Nguling	35
4047	Nguntoronadi	35
4048	Nguntoronadi	33
4049	Ngunut	35
4050	Ngusikan	35
4051	Nguter	33
4052	Nibong	11
4053	Nibung	16
4054	Nibung Hangus	12
4055	Nikiwar	92
4056	Nikogwe	95
4057	Nimbokrang	91
4058	Nimboran	91
4059	Niname	95
4060	Ninati	93
4061	Ninia	95
4062	Nioga	94
4063	Nipah Panjang	15
4064	Nipsan	95
4065	Nirkuri	95
4066	Nirunmas	81
4067	Nisam	11
4068	Nisam Antara	11
4069	Nita	53
4070	Noebana	53
4071	Noebeba	53
4072	Noemuti	53
4073	Noemuti Timur	53
4074	Nogi	95
4075	Nogosari	33
4076	Nokilalaki	72
4077	Nonggunong	35
4078	Nongme	95
4079	Nongsa	21
4080	Nosu	76
4081	Noyan	61
4082	Nuangan	71
4083	Nubatukan	53
4084	Nuha	73
4085	Nuhon	72
4086	Numba	95
4087	Nume	94
4088	Numfor Barat	91
4089	Numfor Timur	91
4090	Nunbena	53
4091	Nunggawi	95
4092	Nunkolo	53
4093	Nunukan	65
4094	Nunukan Selatan	65
4095	Nurussalam	11
4096	Nusa Laut	81
4097	Nusa Penida (Nusapenida)	51
4098	Nusa Tabukan	71
4099	Nusaherang	32
4100	Nusaniwe (Nusanive)	81
4101	Nusawani	91
4102	Nusawungu	33
4103	Nyalindung	32
4104	Nyuatan	64
4105	O'o'u (Oou)	12
4106	Oba	82
4107	Oba Selatan	82
4108	Oba Tengah	82
4109	Oba Utara	82
4110	Obaa	93
4111	Obi	82
4112	Obi Barat	82
4113	Obi Selatan	82
4114	Obi Timur	82
4115	Obi Utara	82
4116	Obio	95
4117	Oebobo	53
4118	Oenino	53
4119	Ogamanim	94
4120	Ogodeide	72
4121	Oheo	74
4122	Ok Aom	95
4123	Okaba	93
4124	Okbab	95
4125	Okbape	95
4126	Okbemtau	95
4127	Okbibab	95
4128	Okhika	95
4129	Oklip	95
4130	Oksamol	95
4131	Oksebang	95
4132	Oksibil	95
4133	Oksop	95
4134	Omben	35
4135	Omesuri	53
4136	Omukia	94
4137	Onan Ganjang	12
4138	Onan Runggu	12
4139	Onembute	74
4140	Oneri	94
4141	Ongka Malino	72
4142	Onohazumba	12
4143	Onolalu	12
4144	Oransbari	92
4145	Oridek	91
4146	Orkeri	91
4147	Orong Telu	52
4148	Oudate	91
4149	Paal Dua	71
4150	Paal Merah	15
4151	Pabean Cantian (Pabean Cantikan)	35
4152	Pabedilan	32
4153	Pabelan	33
4154	Paberiwai	53
4155	Pabuaran	36
4156	Pabuaran	32
4157	Pacar	53
4158	Pace	35
4159	Pacet	35
4160	Pacet	32
4161	Paciran	35
4162	Pacitan	35
4163	Padaherang	32
4164	Padaido	91
4165	Padakembang	32
4166	Padalarang	32
4167	Padamara	33
4168	Padang	35
4169	Padang Barat	13
4170	Padang Batung	63
4171	Padang Bolak	12
4172	Padang Bolak Julu	12
4173	Padang Bolak Tenggara	12
4174	Padang Cermin	18
4175	Padang Ganting	13
4176	Padang Gelugur	13
4177	Padang Guci Hilir	17
4178	Padang Guci Hulu	17
4179	Padang Hilir	12
4180	Padang Hulu	12
4181	Padang Jaya	17
4182	Padang Laweh	13
4183	Padang Panjang Barat	13
4184	Padang Panjang Timur	13
4185	Padang Ratu	18
4186	Padang Sago	13
4187	Padang Selatan	13
4188	Padang Tiji	11
4189	Padang Timur	13
4190	Padang Tualang	12
4191	Padang Ulak Tanding	17
4192	Padang Utara	13
4193	Padangan	35
4194	Padangguni	74
4195	Padangsidimpuan Angkola Julu	12
4196	Padangsidimpuan Batunadua	12
4197	Padangsidimpuan Hutaimbaru	12
4198	Padangsidimpuan Selatan	12
4199	Padangsidimpuan Tenggara	12
4200	Padangsidimpuan Utara	12
4201	Padarincang	36
4202	Padas	35
4203	Pademangan	31
4204	Pademawu	35
4205	Padua	93
4206	Padureso	33
4207	Paga	53
4208	Pagaden	32
4209	Pagaden Barat	32
4210	Pagai Selatan	13
4211	Pagai Utara	13
4212	Pagak	35
4213	Pagaleme	94
4214	Pagar Alam Selatan	16
4215	Pagar Alam Utara	16
4216	Pagar Dewa	18
4217	Pagar Gunung	16
4218	Pagar Jati	17
4219	Pagar Merbau	12
4220	Pagaran	12
4221	Pagaran Tapah Darussalam	14
4222	Pagedangan	36
4223	Pagedongan	33
4224	Pagelaran	36
4225	Pagelaran	32
4226	Pagelaran	35
4227	Pagelaran	18
4228	Pagelaran Utara	18
4229	Pagentan	33
4230	Pagerageung	32
4231	Pagerbarang	33
4232	Pageruyung (Pagerruyung)	33
4233	Pagerwojo	35
4234	Pagimana	72
4235	Pagindar	12
4236	Pagu	35
4237	Paguat	75
4238	Paguyaman	75
4239	Paguyaman Pantai	75
4240	Paguyangan	33
4241	Pahae Jae	12
4242	Pahae Julu	12
4243	Pahandut	62
4244	Pahunga Lodu	53
4245	Paiton	35
4246	Pajangan	34
4247	Pajar Bulan	16
4248	Pajo	52
4249	Paju Epat	62
4250	Pajukukang	73
4251	Pakal	35
4252	Pakantan	12
4253	Pakel	35
4254	Pakem	34
4255	Pakem	35
4256	Pakenjeng	32
4257	Pakis	33
4258	Pakis	35
4259	Pakis Aji	33
4260	Pakisaji	35
4261	Pakisjaya	32
4262	Pakkat	12
4263	Pakong	35
4264	Paku	62
4265	Pakualaman	34
4266	Pakuan Ratu	18
4267	Pakue	74
4268	Pakue Tengah	74
4269	Pakue Utara	74
4270	Pakuhaji	36
4271	Pakuniran	35
4272	Pakusari	35
4273	Pal Merah (Palmerah)	31
4274	Palabuhanratu (Pelabuhanratu)	32
4275	Palakka	73
4276	Palang	35
4277	Palangga	74
4278	Palangga (Pallangga)	73
4279	Palangga Selatan	74
4280	Palaran	64
4281	Palas	18
4282	Palasa	72
4283	Palasah	32
4284	Paleleh	72
4285	Paleleh Barat	72
4286	Palembayan (Pelembayan)	13
4287	Palengaan (Palenggaan, Palenga'an)	35
4288	Paleteang	73
4289	Palibelo	52
4290	Palimanan	32
4291	Palipi	12
4292	Paliyan	34
4293	Palmatak	21
4294	Paloh	61
4295	Palolo	72
4296	Palu Barat	72
4297	Palu Selatan	72
4298	Palu Timur	72
4299	Palu Utara	72
4300	Palue	53
4301	Palupuh	13
4302	Pamanukan	32
4303	Pamarayan	36
4304	Pamarican	32
4305	Pamatang Sidamanik (Pematang Sidamanik)	12
4306	Pamatang Silima Huta (Pematang)	12
4307	Pamboang	76
4308	Pamek	95
4309	Pamekasan	35
4310	Pamenang	15
4311	Pamenang Barat	15
4312	Pamenang Selatan	15
4313	Pameungpeuk	32
4314	Pamijahan	32
4315	Paminggir	63
4316	Pammana	73
4317	Pamona Barat	72
4318	Pamona Puselemba	72
4319	Pamona Selatan	72
4320	Pamona Tenggara	72
4321	Pamona Timur	72
4322	Pamona Utara	72
4323	Pamotan	33
4324	Pampangan	16
4325	Pamukan Barat	63
4326	Pamukan Selatan	63
4327	Pamukan Utara	63
4328	Pamulang	36
4329	Pamulihan	32
4330	Pana	76
4331	Panaga	95
4332	Panai Hilir	12
4333	Panai Hulu	12
4334	Panai Tengah	12
4335	Panakkukang	73
4336	Panang Enim	16
4337	Panarukan	35
4338	Panawangan	32
4339	Panca Jaya	18
4340	Panca Lautang	73
4341	Panca Rijang	73
4342	Pancalang	32
4343	Pancatengah	32
4344	Panceng	35
4345	Pancoran	31
4346	Pancoran Mas	32
4347	Pancung Soal	13
4348	Pancur	33
4349	Pancur Batu	12
4350	Pandaan	35
4351	Pandak	34
4352	Pandan	12
4353	Pandanarum	33
4354	Pandawai	53
4355	Pandawan	63
4356	Pandeglang	36
4357	Pandih Batu	62
4358	Pandrah	11
4359	Panei	12
4360	Panekan	35
4361	Panga (Keude Panga)	11
4362	Pangale	76
4363	Pangalengan	32
4364	Pangandaran	32
4365	Pangarengan	35
4366	Pangaribuan	12
4367	Pangatikan	32
4368	Pangean	14
4369	Pangenan	32
4370	Panggang	34
4371	Panggarangan	36
4372	Panggema	95
4373	Panggul	35
4374	Panggungrejo	35
4375	Pangkah	33
4376	Pangkajene	73
4377	Pangkal Balam	19
4378	Pangkalan	32
4379	Pangkalan Banteng	62
4380	Pangkalan Baru	19
4381	Pangkalan Jambu	15
4382	Pangkalan Kerinci	14
4383	Pangkalan Koto Baru	13
4384	Pangkalan Kuras	14
4385	Pangkalan Lada	62
4386	Pangkalan Lampam	16
4387	Pangkalan Lesung	14
4388	Pangkalan Susu	12
4389	Pangkatan	12
4390	Pangkur	35
4391	Panguragan	32
4392	Pangururan	12
4393	Paniai Barat	94
4394	Paniai Timur	94
4395	Panimbang	36
4396	Paninggaran	33
4397	Panjalu	32
4398	Panjang	18
4399	Panjatan	34
4400	Panji	35
4401	Panombeian Panei / Pane	12
4402	Panongan	36
4403	Pantai Barat	91
4404	Pantai Baru	53
4405	Pantai Cermin	13
4406	Pantai Cermin	12
4407	Pantai Kasuari	93
4408	Pantai Labu	12
4409	Pantai Lunci	62
4410	Pantai Timur	91
4411	Pantai Timur Bagian Barat	91
4412	Pantan Cuaca	11
4413	Pantar	53
4414	Pantar Barat	53
4415	Pantar Baru Laut	53
4416	Pantar Tengah	53
4417	Pantar Timur	53
4418	Pante Bidari	11
4419	Pante Ceureumen (Pantai Ceuremen)	11
4420	Panteraja	11
4421	Panti	35
4422	Panti	13
4423	Panton Reu	11
4424	Panumbangan	32
4425	Panyabungan	12
4426	Panyabungan Barat	12
4427	Panyabungan Selatan	12
4428	Panyabungan Timur	12
4429	Panyabungan Utara	12
4430	Panyileukan	32
4431	Panyingkiran	32
4432	Panyipatan	63
4433	Papalang	76
4434	Papar	35
4435	Parado	52
4436	Parakan	33
4437	Parakansalak (Parakan Salak)	32
4438	Paramasan (Peramasan)	63
4439	Parang	35
4440	Paranggupito	33
4441	Paranginan	12
4442	Parangloe	73
4443	Parbuluan	12
4444	Pardasuka	18
4445	Pare	35
4446	Parengan	35
4447	Parenggean	62
4448	Pariaman Selatan	13
4449	Pariaman Tengah	13
4450	Pariaman Timur	13
4451	Pariaman Utara	13
4452	Pariangan	13
4453	Parigi	73
4454	Parigi	72
4455	Parigi	74
4456	Parigi	32
4457	Parigi Barat	72
4458	Parigi Selatan	72
4459	Parigi Tengah	72
4460	Parigi Utara	72
4461	Parindu	61
4462	Paringin	63
4463	Paringin Selatan	63
4464	Parittiga	19
4465	Pariwari	92
4466	Parlilitan	12
4467	Parmaksian	12
4468	Parmonangan	12
4469	Paro	95
4470	Paron	35
4471	Parongpong	32
4472	Parung	32
4473	Parung Panjang	32
4474	Parungkuda (Parung Kuda)	32
4475	Parungponteng	32
4476	Pasak Talawang	62
4477	Pasaleman	32
4478	Pasaman	13
4479	Pasan	71
4480	Pasangkayu	76
4481	Pasar Jambi	15
4482	Pasar Kemis	36
4483	Pasar Kliwon	33
4484	Pasar Manna	17
4485	Pasar Minggu	31
4486	Pasar Muaro Bungo (Pasar Muara Bungo)	15
4487	Pasar Rebo	31
4488	Pasaribu Tobing	12
4489	Pasarwajo (Pasar Wajo)	74
4490	Pasawahan	32
4491	Pasean	35
4492	Paseh	32
4493	Pasekan	32
4494	Pasema	95
4495	Pasemah Air Keruh	16
4496	Paser Belengkong (Pasir Belengkong)	64
4497	Pasi Kolaga	74
4498	Pasi Raja (Pasie Raja)	11
4499	Pasie Raya	11
4500	Pasilambena	73
4501	Pasimarannu	73
4502	Pasimasunggu (Pasimassunggu)	73
4503	Pasimasunggu Timur	73
4504	Pasir Limau Kapas	14
4505	Pasir Penyu	14
4506	Pasir Putih	74
4507	Pasir Putih	95
4508	Pasir Sakti	18
4509	Pasirian	35
4510	Pasirjambu	32
4511	Pasirkuda	32
4512	Pasirwangi	32
4513	Pasongsongan	35
4514	Pasrepan	35
4515	Pasrujambe (Pasujambe)	35
4516	Passi Barat	71
4517	Passi Timur	71
4518	Passue	93
4519	Passue Bawah	93
4520	Patampanua	73
4521	Patamuan	13
4522	Patangkep Tutui	62
4523	Patani	82
4524	Patani Barat	82
4525	Patani Timur	82
4526	Patani Utara	82
4527	Pataruman	32
4528	Patean	33
4529	Patebon	33
4530	Pati	33
4531	Patia	36
4532	Patianrowo	35
4533	Patikraja	33
4534	Patilanggio	75
4535	Patimpeng	73
4536	Patimuan	33
4537	Patokbeusi	32
4538	Patrang	35
4539	Patrol	32
4540	Pattalasang (Pattallassang)	73
4541	Pattallassang (Patallassang)	73
4542	Patuk	34
4543	Patumbak	12
4544	Pauh	15
4545	Pauh	13
4546	Pauh Duo	13
4547	Paya Bakong	11
4548	Payakumbuh	13
4549	Payakumbuh Barat	13
4550	Payakumbuh Selatan	13
4551	Payakumbuh Timur	13
4552	Payakumbuh Utara	13
4553	Payangan	51
4554	Payaraman	16
4555	Payung	19
4556	Payung	12
4557	Payung Sekaki	14
4558	Payung Sekaki	13
4559	Pebayuran	32
4560	Pecalungan	33
4561	Pecangaan	33
4562	Pedamaran	16
4563	Pedamaran Timur	16
4564	Pedan	33
4565	Pedes	32
4566	Pedongga	76
4567	Pedurungan	33
4568	Pegagan Hilir	12
4569	Pegajahan	12
4570	Pegandon	33
4571	Pegantenan	35
4572	Pegasing	11
4573	Pejagoan	33
4574	Pejarakan (Pajarakan)	35
4575	Pejawaran	33
4576	Pekaitan	14
4577	Pekalipan	32
4578	Pekalongan	18
4579	Pekalongan Barat	33
4580	Pekalongan Selatan	33
4581	Pekalongan Timur	33
4582	Pekalongan Utara	33
4583	Pekanbaru Kota	14
4584	Pekat	52
4585	Pekuncen	33
4586	Pekutatan	51
4587	Pelaihari	63
4588	Pelalawan	14
4589	Pelangiran	14
4590	Pelawan	15
4591	Pelayangan	15
4592	Pelebaga	95
4593	Pelepat	15
4594	Pelepat Ilir	15
4595	Peling Tengah	72
4596	Pemahan	61
4597	Pemalang	33
4598	Pemali	19
4599	Pemangkat	61
4600	Pematang Bandar	12
4601	Pematang Jaya	12
4602	Pematang Karau	62
4603	Pematang Sawa	18
4604	Pematang Tiga	17
4605	Pemayung	15
4606	Pemenang	52
4607	Pemulutan	16
4608	Pemulutan Barat	16
4609	Pemulutan Selatan	16
4610	Penajam	64
4611	Penanggalan	11
4612	Penarik	17
4613	Penawangan	33
4614	Penawar Aji	18
4615	Penawar Tama	18
4616	Pendalian IV Koto	14
4617	Pendopo	16
4618	Pendopo Barat	16
4619	Penebel	51
4620	Penengahan	18
4621	Pengabuan	15
4622	Pengadegan	33
4623	Pengandonan	16
4624	Pengaron (Pengarom)	63
4625	Pengasih	34
4626	Pengkadan (Batu Datu)	61
4627	Peninjauan	16
4628	Penjaringan	31
4629	Penrang	73
4630	Penukal	16
4631	Penukal Utara	16
4632	Penyinggahan	64
4633	Pepera	95
4634	Perak	35
4635	Peranap	14
4636	Perbaungan	12
4637	Percut Sei Tuan	12
4638	Pergetteng Getteng Sengkut	12
4639	Perhentian Raja	14
4640	Periuk	36
4641	Permata	11
4642	Permata Intan	62
4643	Permata Kecubung	62
4644	Pesanggaran	35
4645	Pesanggrahan	31
4646	Pesantren	35
4647	Pesisir Bukit	15
4648	Pesisir Selatan	18
4649	Pesisir Tengah	18
4650	Pesisir Utara	18
4651	Peso	65
4652	Peso Hilir (Ilir)	65
4653	Petak Malai	62
4654	Petanahan	33
4655	Petang	51
4656	Petarukan	33
4657	Petasia	72
4658	Petasia Barat	72
4659	Petasia Timur	72
4660	Peterongan	35
4661	Petir	36
4662	Petungkriyono (Petungkriono)	33
4663	Peudada	11
4664	Peudawa	11
4665	Peukan Bada	11
4666	Peukan Baro	11
4667	Peulimbang (Plimbang)	11
4668	Peunaron	11
4669	Peundeuy	32
4670	Peureulak	11
4671	Peureulak Barat	11
4672	Peureulak Timur	11
4673	Peusangan	11
4674	Peusangan Selatan	11
4675	Peusangan Siblah Krueng	11
4676	Piani	63
4677	Picung	36
4678	Pidie	11
4679	Pija	95
4680	Pilangkenceng	35
4681	Pinang (Penang)	36
4682	Pinang Belapis	17
4683	Pinang Raya	17
4684	Pinangsori	12
4685	Pineleng	71
4686	Pinembani (Panembani)	72
4687	Pinggir	14
4688	Pining (Pinding)	11
4689	Pino	17
4690	Pino Raya (Pinoraya)	17
4691	Pinogaluman	71
4692	Pinogu	75
4693	Pinoh Selatan	61
4694	Pinoh Utara	61
4695	Pinolosian	71
4696	Pinolosian Tengah	71
4697	Pinolosian Timur	71
4698	Pintu Pohan Meranti	12
4699	Pintu Rime Gayo	11
4700	Pinu Pahar (Pinupahar / Pirapahar)	53
4701	Pipikoro	72
4702	Pirak Timur	11
4703	Piramid	95
4704	Pirime	95
4705	Pisugi	95
4706	Pitu	35
4707	Pitu Riase	73
4708	Pitu Riawa	73
4709	Pitumpanua	73
4710	Pituruh	33
4711	Piyaiye	94
4712	Piyungan	34
4713	Plaju	16
4714	Plakat Tinggi (Pelakat Tinggi)	16
4715	Plampang	52
4716	Plandaan	35
4717	Plantungan	33
4718	Plaosan	35
4719	Playen	34
4720	Plemahan	35
4721	Plered	32
4722	Pleret	34
4723	Ploso	35
4724	Plosoklaten	35
4725	Plumbon	32
4726	Plumpang	35
4727	Plupuh	33
4728	Poasia	74
4729	Poga	95
4730	Pogalan	35
4731	Poganeri	95
4732	Pogoma	94
4733	Pohjentrek	35
4734	Poigar	71
4735	Poiru	91
4736	Polanharjo	33
4737	Poleang	74
4738	Poleang Barat	74
4739	Poleang Selatan	74
4740	Poleang Tengah	74
4741	Poleang Tenggara	74
4742	Poleang Timur	74
4743	Poleang Utara	74
4744	Polen	53
4745	Polewali	76
4746	Poli Polia	74
4747	Polinggona	74
4748	Pollung	12
4749	Polokarto	33
4750	Polongbangkeng Selatan (Polombangkeng)	73
4751	Polongbangkeng Timur	73
4752	Polongbangkeng Utara (Polombangkeng)	73
4753	Pomalaa	74
4754	Poncokusumo	35
4755	Poncol	35
4756	Poncowarno	33
4757	Pondidaha	74
4758	Pondok Aren	36
4759	Pondok Kelapa	17
4760	Pondok Kubang	17
4761	Pondok Suguh	17
4762	Pondok Tinggi	15
4763	Pondokgede (Pondok Gede)	32
4764	Pondokmelati (Pondok Melati)	32
4765	Pondoksalam	32
4766	Ponelo Kepulauan	75
4767	Ponggok	35
4768	Ponjong	34
4769	Ponorogo	35
4770	Ponrang	73
4771	Ponrang Selatan	73
4772	Ponre	73
4773	Pontang	36
4774	Pontianak Barat	61
4775	Pontianak Kota	61
4776	Pontianak Selatan	61
4777	Pontianak Tenggara	61
4778	Pontianak Timur	61
4779	Pontianak Utara	61
4780	Poom	91
4781	Popayato	75
4782	Popayato Barat	75
4783	Popayato Timur	75
4784	Popugoba	95
4785	Porehu	74
4786	Porong	35
4787	Porsea	12
4788	Portibi	12
4789	Posigadan	71
4790	Poso Kota	72
4791	Poso Kota Selatan	72
4792	Poso Kota Utara	72
4793	Poso Pesisir	72
4794	Poso Pesisir Selatan	72
4795	Poso Pesisir Utara	72
4796	Poto Tano	52
4797	Prabumulih Barat	16
4798	Prabumulih Selatan	16
4799	Prabumulih Timur	16
4800	Prabumulih Utara	16
4801	Pracimantoro	33
4802	Prafi	92
4803	Pragaan	35
4804	Prajekan	35
4805	Prajuritkulon (Prajurit Kulon)	35
4806	Prambanan	34
4807	Prambanan	33
4808	Prambon	35
4809	Praya	52
4810	Praya Barat	52
4811	Praya Barat Daya	52
4812	Praya Tengah	52
4813	Praya Timur	52
4814	Prembun	33
4815	Prigen	35
4816	Pringapus	33
4817	Pringgabaya	52
4818	Pringgarata	52
4819	Pringgasela	52
4820	Pringkuku	35
4821	Pringsewu	18
4822	Pringsurat	33
4823	Pronggoli	95
4824	Pronojiwo	35
4825	Proppo	35
4826	Pseksu	16
4827	Pubian	18
4828	Pucakwangi	33
4829	Pucanglaban	35
4830	Pucuk	35
4831	Pucuk Rantau	14
4832	Pudak	35
4833	Puding Besar	19
4834	Pugaan	63
4835	Puger	35
4836	Pugo Dagi	94
4837	Pugung	18
4838	Puhpelem	33
4839	Pujananting	73
4840	Pujer	35
4841	Pujon	35
4842	Pujud	14
4843	Pujungan	65
4844	Pujut	52
4845	Pulau Banyak	11
4846	Pulau Banyak Barat	11
4847	Pulau Batang Dua	82
4848	Pulau Beringin	16
4849	Pulau Besar (Pulaubesar)	19
4850	Pulau Burung	14
4851	Pulau Derawan	64
4852	Pulau Dullah Selatan	81
4853	Pulau Dullah Utara	81
4854	Pulau Ende	53
4855	Pulau Gebe	82
4856	Pulau Gorom	81
4857	Pulau Hanaut	62
4858	Pulau Haruku	81
4859	Pulau Hiri	82
4860	Pulau Kurudu	91
4861	Pulau Lakor	81
4862	Pulau Laut	21
4863	Pulau Laut Barat (Pulaulaut Barat)	63
4864	Pulau Laut Kepulauan (Pulaulaut Kepulauan)	63
4865	Pulau Laut Selatan (Pulaulaut Selatan)	63
4866	Pulau Laut Tanjung Selayar	63
4867	Pulau Laut Tengah (Pulaulaut Tengah)	63
4868	Pulau Laut Timur (Pulaulaut Timur)	63
4869	Pulau Laut Utara (Pulaulaut Utara)	63
4870	Pulau Leti (Letti Moa Lakor)	81
4871	Pulau Makian	82
4872	Pulau Malan	62
4873	Pulau Masela	81
4874	Pulau Maya (Pulau Maya Karimata)	61
4875	Pulau Panggung	18
4876	Pulau Panjang	81
4877	Pulau Panjang	21
4878	Pulau Petak	62
4879	Pulau Pisang (Pulaupisang)	18
4880	Pulau Punjung	13
4881	Pulau Pura	53
4882	Pulau Rakyat	12
4883	Pulau Rao	82
4884	Pulau Rimau	16
4885	Pulau Sebuku (Pulausebuku)	63
4886	Pulau Seluan	21
4887	Pulau Sembilan	73
4888	Pulau Sembilan (Pulausembilan)	63
4889	Pulau Ternate	82
4890	Pulau Tiga	21
4891	Pulau Tiga	93
4892	Pulau Tiga Barat	21
4893	Pulau Wetang	81
4894	Pulau Yerui	91
4895	Pulau-Pulau Aru	81
4896	Pulau-pulau Babar Timur	81
4897	Pulau-Pulau Batu	12
4898	Pulau-Pulau Batu Barat	12
4899	Pulau-Pulau Batu Timur	12
4900	Pulau-Pulau Batu Utara	12
4901	Pulau-Pulau Kur	81
4902	Pulaulaut Sigam	63
4903	Pulaumerbau	14
4904	Pulaupinang (Pulau Pinang)	16
4905	Puldama	95
4906	Pule	35
4907	Pulo Aceh	11
4908	Pulo Ampel	36
4909	Pulo Bandring	12
4910	Pulogadung (Pulo Gadung)	31
4911	Pulokulon	33
4912	Pulomerak	36
4913	Pulosari	36
4914	Pulosari	33
4915	Pulubala	75
4916	Pulung	35
4917	Pulutan	71
4918	Puncak Sorik Marapi	12
4919	Puncu	35
4920	Pundong	34
4921	Punduh Pidada	18
4922	Punggelan	33
4923	Pungging	35
4924	Punggur	18
4925	Punung	35
4926	Pupuan	51
4927	Purabaya	32
4928	Purba	12
4929	Purba Tua (Purbatua)	12
4930	Purbalingga	33
4931	Purbaratu	32
4932	Purbolinggo	18
4933	Pureman	53
4934	Puri	35
4935	Puriala	74
4936	Puring	33
4937	Puring Kencana	61
4938	Purwadadi	32
4939	Purwaharja	32
4940	Purwakarta	36
4941	Purwakarta	32
4942	Purwanegara (Purwonegoro)	33
4943	Purwantoro	33
4944	Purwasari	32
4945	Purwoasri	35
4946	Purwodadi	33
4947	Purwodadi	35
4948	Purwodadi	16
4949	Purwoharjo	35
4950	Purwojati	33
4951	Purwokerto Barat	33
4952	Purwokerto Selatan	33
4953	Purwokerto Timur	33
4954	Purwokerto Utara	33
4955	Purworeja Klampok (Purworejo Klampok)	33
4956	Purworejo	33
4957	Purworejo	35
4958	Purwosari	35
4959	Purwosari	34
4960	Pusakajaya	32
4961	Pusakanagara	32
4962	Pusako	14
4963	Pusomaen	71
4964	Puspahiang	32
4965	Puspo	35
4966	Puteri Betung (Putri Betung)	11
4967	Putra Rumbia	18
4968	Putri Hijau	17
4969	Putussibau Selatan	61
4970	Putussibau Utara	61
4971	Puuwatu	74
4972	Ra'as (Raas)	35
4973	Raba	52
4974	Rahong Utara	53
4975	Rahuning (Rahunig)	12
4976	Raihat	53
4977	Raijua	53
4978	Raimanuk	53
4979	Raimbawi	91
4980	Rainis	71
4981	Raja Basa (Rajabasa)	18
4982	Rajabasa	18
4983	Rajadesa	32
4984	Rajagaluh	32
4985	Rajapolah	32
4986	Rajeg	36
4987	Rakit	33
4988	Rakit Kulim	14
4989	Rakumpit	62
4990	Raman Utara	18
4991	Rambah	14
4992	Rambah Hilir	14
4993	Rambah Samo	14
4994	Rambang	16
4995	Rambang Kapak Tengah	16
4996	Rambang Kuang	16
4997	Rambang Niru (Rambang Dangku)	16
4998	Rambatan	13
4999	Rambipuji	35
5000	Rambutan	16
5001	Rambutan	12
5002	Rampi	73
5003	Rana Mese	53
5004	Ranah Ampek Hulu Tapan	13
5005	Ranah Batahan	13
5006	Ranah Pesisir	13
5007	Ranca Bungur	32
5008	Rancabali (Ranca Bali)	32
5009	Rancaekek	32
5010	Rancah	32
5011	Rancakalong	32
5012	Rancasari	32
5013	Randangan	75
5014	Randuagung	35
5015	Randublatung	33
5016	Randudongkal	33
5017	Rangkasbitung	36
5018	Rangkui	19
5019	Rangsang	14
5020	Rangsang Barat	14
5021	Rangsang Pesisir	14
5022	Rano	73
5023	Ranomeeto	74
5024	Ranomeeto Barat	74
5025	Ranowulu (Bitung Utara)	71
5026	Ranoyapo	71
5027	Ransiki	92
5028	Rantau	11
5029	Rantau Alai	16
5030	Rantau Badauh	63
5031	Rantau Bayur	16
5032	Rantau Kopar	14
5033	Rantau Pandan	15
5034	Rantau Panjang	16
5035	Rantau Peureulak (Ranto Peureulak)	11
5036	Rantau Pulung	64
5037	Rantau Rasau	15
5038	Rantau Selamat	11
5039	Rantau Selatan	12
5040	Rantau Utara	12
5041	Rante Angin	74
5042	Rantebua	73
5043	Rantebulahan Timur	76
5044	Rantepao	73
5045	Rantetayo	73
5046	Ranto Baek	12
5047	Ranuyoso	35
5048	Rao	13
5049	Rao Selatan	13
5050	Rao Utara	13
5051	Rappocini	73
5052	Raren Batuah	62
5053	Rarowatu	74
5054	Rarowatu Utara	74
5055	Rasanae Barat	52
5056	Rasanae Timur	52
5057	Rasau Jaya	61
5058	Rasiei	92
5059	Ratahan	71
5060	Ratahan Timur	71
5061	Ratatotok	71
5062	Ratolindo	72
5063	Ratu Agung	17
5064	Ratu Samban	17
5065	Raveni Rara (Ravenirara)	91
5066	Rawa Jitu Selatan (Rawajitu Selatan)	18
5067	Rawa Jitu Timur (Rawajitu Timur)	18
5068	Rawa Jitu Utara	18
5069	Rawa Pitu	18
5070	Rawalo	33
5071	Rawalumbu	32
5072	Rawamerta	32
5073	Rawang Panca Arga	12
5074	Rawas Ilir	16
5075	Rawas Ulu	16
5076	Raya	12
5077	Raya Kahean	12
5078	Reban	33
5079	Rebang Tangkas	18
5080	Regol	32
5081	Rejoso	35
5082	Rejotangan	35
5083	Rembang	35
5084	Rembang	33
5085	Remboken	71
5086	Rembon	73
5087	Renah Mendaluh	15
5088	Renah Pamenang (Renah Pemenang)	15
5089	Renah Pembarap	15
5090	Rendang	51
5091	Rengasdengklok	32
5092	Rengat	14
5093	Rengat Barat	14
5094	Rengel	35
5095	Reok	53
5096	Reok Barat	53
5097	Reteh	14
5098	Rhee	52
5099	Riau Silip	19
5100	Rikit Gaib	11
5101	Rilauale (Rilau Ale)	73
5102	Rimba Melintang	14
5103	Rimbo Bujang	15
5104	Rimbo Ilir	15
5105	Rimbo Pengadang	17
5106	Rimbo Tengah	15
5107	Rimbo Ulu	15
5108	Rindi	53
5109	Rindingallo	73
5110	Ringinarum	33
5111	Ringinrejo	35
5112	Rinhat	53
5113	Rio Pakava	72
5114	Risei Sayati	91
5115	Riung	53
5116	Riung Barat	53
5117	Robatal	35
5118	Rogojampi	35
5119	Rokan IV Koto	14
5120	Rongga	32
5121	Ronggur Nihuta	12
5122	Rongkong (Limbong)	73
5123	Rongkop	34
5124	Roon	92
5125	Ropang	52
5126	Roswar	92
5127	Rote Barat	53
5128	Rote Barat Daya	53
5129	Rote Barat Laut	53
5130	Rote Selatan	53
5131	Rote Tengah	53
5132	Rote Timur	53
5133	Routa	74
5134	Rowokangkung	35
5135	Rowokele	33
5136	Rowosari	33
5137	Rubaru	35
5138	Rufaer	91
5139	Rumbai (Rumbai Pesisir)	14
5140	Rumbai Barat	14
5141	Rumbai Timur	14
5142	Rumberpon	92
5143	Rumbia	73
5144	Rumbia	74
5145	Rumbia	18
5146	Rumbia Tengah	74
5147	Rumbio Jaya	14
5148	Rumpin	32
5149	Rundeng	11
5150	Rungan	62
5151	Rungan Barat	62
5152	Rungan Hulu	62
5153	Rungkut	35
5154	Runjung Agung	16
5155	Rupat	14
5156	Rupat Utara	14
5157	Rupit	16
5158	Rusip Antara	11
5159	Ruteng	53
5160	Sa'dan	73
5161	Sabak Auh	14
5162	Sabangau (Sebangau)	62
5163	Sabangparu	73
5164	Sabbang	73
5165	Sabbang Selatan	73
5166	Sabu Barat	53
5167	Sabu Liae	53
5168	Sabu Tengah	53
5169	Sabu Timur	53
5170	Sabulakoa	74
5171	Sadananya	32
5172	Sadang	33
5173	Sadaniang	61
5174	Sadu	15
5175	Saengkeduk	92
5176	Safan	93
5177	Sagalaherang	32
5178	Sagaranten	32
5179	Saguling	32
5180	Sagulung	21
5181	Sahu	82
5182	Sahu Timur	82
5183	Saifi	92
5184	Sail	14
5185	Saipar Dolok Hole	12
5186	Sajad	61
5187	Sajingan Besar	61
5188	Sajira	36
5189	Sajoanging	73
5190	Saketi	36
5191	Sako	16
5192	Sakra	52
5193	Sakra Barat	52
5194	Sakra Timur	52
5195	Sakti	11
5196	Salahutu	81
5197	Salak	12
5198	Salam	33
5199	Salam Babaris	63
5200	Salaman	33
5201	Salang	11
5202	Salapian	12
5203	Salatiga	61
5204	Salawati	92
5205	Salawati Barat	92
5206	Salawati Selatan	92
5207	Salawati Tengah	92
5208	Salawati Utara	92
5209	Salawu	32
5210	Sale	33
5211	Salem	33
5212	Salibabu	71
5213	Salimpaung (Salimpauang)	13
5214	Saling	16
5215	Salkma	92
5216	Salo	14
5217	Salomekko	73
5218	Salopa	32
5219	Saluputi (Saluputti)	73
5220	Samadua (Sama Dua)	11
5221	Samalanga	11
5222	Samalantan	61
5223	Samarang	32
5224	Samarinda Ilir	64
5225	Samarinda Kota	64
5226	Samarinda Seberang	64
5227	Samarinda Ulu	64
5228	Samarinda Utara	64
5229	Samatiga	11
5230	Samaturu	74
5231	Sambaliung	64
5232	Sambas	61
5233	Sambelia (Sambalia)	52
5234	Sambeng	35
5235	Sambi	33
5236	Sambi Rampas	53
5237	Sambikerep	35
5238	Sambirejo	33
5239	Sambit	35
5240	Samboja (Semboja)	64
5241	Samboja Barat	64
5242	Sambong	33
5243	Sambung Makmur	63
5244	Sambungmacan (Sambung Macan)	33
5245	Sambutan	64
5246	Samenage	95
5247	Samigaluh	34
5248	Samofa	91
5249	Sampaga	76
5250	Sampanahan	63
5251	Sampang	35
5252	Sampang	33
5253	Sampara	74
5254	Sampoi Niet (Sampoiniet)	11
5255	Sampolawa	74
5256	Sampung	35
5257	Samudera	11
5258	Sanaman Mantikei (Senamang Mantikei)	62
5259	Sanana	82
5260	Sanana Utara	82
5261	Sanankulon (Sanan Kulon)	35
5262	Sananwetan (Sanan Wetan)	35
5263	Sandai	61
5264	Sandaran	64
5265	Sanden	34
5266	Sandubaya (Sandujaya)	52
5267	Sang Tombolang	71
5268	Sanga Desa	16
5269	Sanga Sanga	64
5270	Sangalla (Sanggala)	73
5271	Sangalla Selatan	73
5272	Sangalla Utara	73
5273	Sangatta Selatan	64
5274	Sangatta Utara	64
5275	Sanggalangi	73
5276	Sanggar	52
5277	Sanggau Ledo	61
5278	Sangia Wambulu	74
5279	Sangir	13
5280	Sangir Balai Janggo	13
5281	Sangir Batang Hari	13
5282	Sangir Jujuan	13
5283	Sangkapura	35
5284	Sangkub	71
5285	Sangkulirang	64
5286	Sano Nggoang	53
5287	Sanrobone	73
5288	Santian	53
5289	Saparua	81
5290	Saparua Timur	81
5291	Sape	52
5292	Sapeken	35
5293	Saptosari (Sapto Sari)	34
5294	Sapuran	33
5295	Saradan	35
5296	Sarang	33
5297	Sario	71
5298	Sarirejo	35
5299	Sariwangi	32
5300	Sarjo	76
5301	Sarmi	91
5302	Sarmi Selatan	91
5303	Sarmi Timur	91
5304	Sarolangun	15
5305	Saronggi	35
5306	Sarudik	12
5307	Sarudu	76
5308	Sasak Ranah Pasisie (Pesisir, Pasisir, Pesisie)	13
5309	Sasitamean	53
5310	Satar Mese	53
5311	Satar Mese Barat	53
5312	Satar Mese Utara	53
5313	Satui	63
5314	Sausapor	92
5315	Sausu	72
5316	Sawa	74
5317	Sawa Erma	93
5318	Sawah Besar	31
5319	Sawahan	35
5320	Sawai	91
5321	Sawan	51
5322	Sawang	11
5323	Sawangan	32
5324	Sawangan	33
5325	Sawerigadi	74
5326	Sawiat	92
5327	Sawit	33
5328	Sawit Seberang	12
5329	Sawo	12
5330	Sawoo	35
5331	Sayan	61
5332	Sayosa	92
5333	Sayosa Timur	92
5334	Sayung	33
5335	Sayur Matinggi	12
5336	Sebangau Kuala	62
5337	Sebangki	61
5338	Sebatik	65
5339	Sebatik Barat	65
5340	Sebatik Tengah	65
5341	Sebatik Timur	65
5342	Sebatik Utara	65
5343	Sebawi	61
5344	Seberang Kota	15
5345	Seberang Musi	17
5346	Seberang Ulu Dua (Seberang Ulu II)	16
5347	Seberang Ulu Satu (Seberang Ulu I)	16
5348	Seberuang	61
5349	Sebuku	65
5350	Sebulu	64
5351	Secang	33
5352	Secanggang	12
5353	Sedan	33
5354	Sedati	35
5355	Sedayu	34
5356	Sedong	32
5357	Segah	64
5358	Segedong	61
5359	Segeri	73
5360	Seget	92
5361	Seginim	17
5362	Segun	92
5363	Sei Balai	12
5364	Sei Bamban	12
5365	Sei Beduk (Sungai Beduk)	21
5366	Sei Bingai (Sei Binge / Bingei)	12
5367	Sei Dadap	12
5368	Sei Kepayang	12
5369	Sei Kepayang Barat	12
5370	Sei Kepayang Timur	12
5371	Sei Lepan	12
5372	Sei Menggaris	65
5373	Sei Rampah	12
5374	Sei Suka	12
5375	Sei Tualang Raso	12
5376	Sejangkung	61
5377	Sekadau Hilir	61
5378	Sekadau Hulu	61
5379	Sekampung	18
5380	Sekampung Udik	18
5381	Sekar	35
5382	Sekaran	35
5383	Sekarbela	52
5384	Sekatak	65
5385	Sekayam	61
5386	Sekayu	16
5387	Sekerak	11
5388	Sekernan	15
5389	Sekincau	18
5390	Seko	73
5391	Sekolaq Darat	64
5392	Sekongkang	52
5393	Sekotong	52
5394	Sekupang	21
5395	Sela	95
5396	Selaawi	32
5397	Selagai Lingga	18
5398	Selagan Raya	17
5399	Selajambe	32
5400	Selakau	61
5401	Selakau Timur	61
5402	Selangit	16
5403	Selaprang (Selaparang)	52
5404	Selaru	81
5405	Selat	62
5406	Selat	51
5407	Selat Gelam	21
5408	Selat Nasik	19
5409	Selat Penuguan	16
5410	Selayar	21
5411	Selebar	17
5412	Selemadeg	51
5413	Selemadeg Barat (Salemadeg Barat)	51
5414	Selemadeg Timur (Salamadeg Timur, Salemadeg Timur)	51
5415	Selemkai	92
5416	Selesai	12
5417	Selimbau	61
5418	Selo	33
5419	Selogiri	33
5420	Selomerto	33
5421	Selong	52
5422	Selopampang	33
5423	Selopuro	35
5424	Selorejo	35
5425	Seluas	61
5426	Seluma	17
5427	Seluma Barat	17
5428	Seluma Selatan	17
5429	Seluma Timur	17
5430	Seluma Utara	17
5431	Selupu Rejang	17
5432	Semadam	11
5433	Semaka	18
5434	Semampir	35
5435	Semanding	35
5436	Semangga	93
5437	Semanu	34
5438	Semarang Barat	33
5439	Semarang Selatan	33
5440	Semarang Tengah	33
5441	Semarang Timur	33
5442	Semarang Utara	33
5443	Sematangborang (Sematang Borang)	16
5444	Sematu Jaya	62
5445	Semau	53
5446	Semau Selatan	53
5447	Sembakung	65
5448	Sembakung Atulai	65
5449	Sembalun	52
5450	Sembawa	16
5451	Sembilan Koto	13
5452	Semboro	35
5453	Semen	35
5454	Semendawai Barat	16
5455	Semendawai Suku III	16
5456	Semendawai Timur	16
5457	Semende Darat Laut	16
5458	Semende Darat Tengah	16
5459	Semende Darat Ulu	16
5460	Semidang Aji	16
5461	Semidang Alas	17
5462	Semidang Alas Maras	17
5463	Semidang Gumay / Gumai	17
5464	Semidang Lagan	17
5465	Semin	34
5466	Semitau	61
5467	Semparuk	61
5468	Sempol	35
5469	Sempor	33
5470	Sempu	35
5471	Senapelan	14
5472	Senayang	21
5473	Sendana	76
5474	Sendana	73
5475	Sendang	35
5476	Sendang Agung	18
5477	Senduro	35
5478	Senen	31
5479	Sengah Temila	61
5480	Senggi	91
5481	Senopi	92
5482	Senori	35
5483	Sentajo Raya	14
5484	Sentani	91
5485	Sentani Barat	91
5486	Sentani Timur	91
5487	Sentolo	34
5488	Senyerang	15
5489	Sepaku	64
5490	Sepang (Sepang Simin)	62
5491	Sepatan	36
5492	Sepatan Timur	36
5493	Sepauk	61
5494	Seponti	61
5495	Sepulu	35
5496	Seputih Agung	18
5497	Seputih Banyak	18
5498	Seputih Mataram	18
5499	Seputih Raman	18
5500	Seputih Surabaya	18
5501	Seradala	95
5502	Serai Serumpun	15
5503	Seram Barat	81
5504	Seram Timur	81
5505	Seram Utara	81
5506	Seram Utara Barat	81
5507	Seram Utara Timur Kobi	81
5508	Seram Utara Timur Seti	81
5509	Serambakon	95
5510	Seranau	62
5511	Serang	36
5512	Serang Baru	32
5513	Serangpanjang	32
5514	Serasan	21
5515	Serasan Timur	21
5516	Serawai	61
5517	Serba Jadi	12
5518	Serbajadi	11
5519	Seremuk	92
5520	Serengan	33
5521	Seri Kuala Lobam (Sri)	21
5522	Seribu Riam	62
5523	Seririt	51
5524	Serpong	36
5525	Serpong Utara	36
5526	Seruway	11
5527	Seruyan Hilir	62
5528	Seruyan Hilir Timur	62
5529	Seruyan Hulu	62
5530	Seruyan Raya	62
5531	Seruyan Tengah	62
5532	Sesayap	65
5533	Sesayap Hilir	65
5534	Sesean	73
5535	Sesean Suloara	73
5536	Sesenapadang	76
5537	Sesnuk	93
5538	Seteluk (Sateluk)	52
5539	Setia	11
5540	Setia Bhakti (Setia Bakti)	11
5541	Setia Janji	12
5542	Setiabudi (Setia Budi)	31
5543	Setu	36
5544	Setu	32
5545	Seulimeum	11
5546	Seunagan	11
5547	Seunagan Timur	11
5548	Seunuddon (Seunudon)	11
5549	Sewon	34
5550	Seyegan	34
5551	Siabu	12
5552	Siak	14
5553	Siak Hulu	14
5554	Siak Kecil	14
5555	Sianjar Mula Mula (Sianjur)	12
5556	Siantan	21
5557	Siantan Selatan	21
5558	Siantan Tengah	21
5559	Siantan Timur	21
5560	Siantan Utara	21
5561	Siantar	12
5562	Siantar Barat	12
5563	Siantar Marihat	12
5564	Siantar Marimbun	12
5565	Siantar Martoba	12
5566	Siantar Narumonda	12
5567	Siantar Selatan	12
5568	Siantar Sitalasari	12
5569	Siantar Timur	12
5570	Siantar Utara	12
5571	Siatas Barita	12
5572	Siau Barat	71
5573	Siau Barat Selatan	71
5574	Siau Barat Utara	71
5575	Siau Tengah	71
5576	Siau Timur	71
5577	Siau Timur Selatan	71
5578	Sibabangun	12
5579	Siberida (Seberida)	14
5580	Siberut Barat	13
5581	Siberut Barat Daya	13
5582	Siberut Selatan	13
5583	Siberut Tengah	13
5584	Siberut Utara	13
5585	Sibolangit	12
5586	Sibolga Kota	12
5587	Sibolga Sambas	12
5588	Sibolga Selatan	12
5589	Sibolga Utara	12
5590	Siborong-Borong	12
5591	Sibulue	73
5592	Sidamanik	12
5593	Sidamulih	32
5594	Sidareja	33
5595	Sidayu	35
5596	Sidemen	51
5597	Sidey	92
5598	Sidikalang	12
5599	Siding	61
5600	Sidoan	72
5601	Sidoarjo	35
5602	Sidoharjo	33
5603	Sidomukti	33
5604	Sidomulyo	18
5605	Sidorejo	35
5606	Sidorejo	33
5607	Sidua'ori	12
5608	Siempat Nempu	12
5609	Siempat Nempu Hilir	12
5610	Siempat Nempu Hulu	12
5611	Siempat Rube	12
5612	Siepkosi	95
5613	Sigaluh	33
5614	Sigi Biromaru	72
5615	Sigumpar	12
5616	Sihapas Barumun	12
5617	Sijamapolang (Sijama Polang)	12
5618	Sijuk	19
5619	Sijunjung	13
5620	Sikakap	13
5621	Sikap Dalam	16
5622	Sikur	52
5623	Silaen	12
5624	Silahisabungan (Silahi Sabungan)	12
5625	Silangkitang	12
5626	Silat Hilir	61
5627	Silat Hulu	61
5628	Silau Laut	12
5629	Silaut	13
5630	Silian Raya	71
5631	Silih Nara	11
5632	Silima Pungga Pungga	12
5633	Silimakuta	12
5634	Silimo	95
5635	Silinda	12
5636	Siliragung	35
5637	Silo	35
5638	Silo Karno Doga	95
5639	Silou Kahean	12
5640	Silungkang	13
5641	Siluq Ngurai	64
5642	Siman	35
5643	Simangambat	12
5644	Simangumban	12
5645	Simanindo	12
5646	Simbang	73
5647	Simboro (Simboro dan Kepulauan)	76
5648	Simbuang	73
5649	Simeulue Barat (Simeuleu Barat)	11
5650	Simeulue Cut	11
5651	Simeulue Tengah (Simeuleu Tengah)	11
5652	Simeulue Timur (Simeuleu Timur)	11
5653	Simo	33
5654	Simokerto	35
5655	Simpang	16
5656	Simpang Alahan Mati	13
5657	Simpang Dua	61
5658	Simpang Empat	63
5659	Simpang Empat	12
5660	Simpang Hilir	61
5661	Simpang Hulu	61
5662	Simpang Jernih	11
5663	Simpang Kanan	11
5664	Simpang Kanan	14
5665	Simpang Katis	19
5666	Simpang Kiri	11
5667	Simpang Kramat (Keramat)	11
5668	Simpang Mamplam	11
5669	Simpang Pematang	18
5670	Simpang Pesak	19
5671	Simpang Raya	72
5672	Simpang Renggiang	19
5673	Simpang Rimba	19
5674	Simpang Teritip	19
5675	Simpang Tiga	11
5676	Simpang Ulim	11
5677	Simpenan	32
5678	Simpur	63
5679	Simuk	12
5680	Sinaboi (Senaboi)	14
5681	Sinak	94
5682	Sinak Barat	94
5683	Sinar Peninjauan	16
5684	Sindang	32
5685	Sindang Beliti Ilir	17
5686	Sindang Beliti Ulu	17
5687	Sindang Danau	16
5688	Sindang Dataran (Sindang Daratan)	17
5689	Sindang Jaya	36
5690	Sindang Kelingi	17
5691	Sindangagung (Sindang Agung)	32
5692	Sindangbarang	32
5693	Sindangkasih	32
5694	Sindangkerta	32
5695	Sindangresmi	36
5696	Sindangwangi	32
5697	Sindue	72
5698	Sindue Tobata	72
5699	Sindue Tombusabora	72
5700	Sine	35
5701	Singajaya	32
5702	Singaparna	32
5703	Singaran Pati	17
5704	Singgahan	35
5705	Singingi	14
5706	Singingi Hilir	14
5707	Singkawang Barat	61
5708	Singkawang Selatan	61
5709	Singkawang Tengah	61
5710	Singkawang Timur	61
5711	Singkawang Utara	61
5712	Singkep	21
5713	Singkep Barat	21
5714	Singkep Pesisir	21
5715	Singkep Selatan	21
5716	Singkil	11
5717	Singkil	71
5718	Singkil Utara	11
5719	Singkohor	11
5720	Singkup	61
5721	Singkut	15
5722	Singojuruh	35
5723	Singorojo	33
5724	Singosari	35
5725	Siniu	72
5726	Sinjai Barat	73
5727	Sinjai Borong	73
5728	Sinjai Selatan	73
5729	Sinjai Tengah	73
5730	Sinjai Timur	73
5731	Sinjai Utara	73
5732	Sinoa	73
5733	Sinonsayang	71
5734	Sintang	61
5735	Sintuak Toboh Gadang	13
5736	Sinunukan	12
5737	Siompu	74
5738	Siompu Barat	74
5739	Siotapina (Siontapia / Siontapina)	74
5740	Sipahutar	12
5741	Sipatana	75
5742	Sipirok	12
5743	Sipispis	12
5744	Sipoholon	12
5745	Sipora Selatan	13
5746	Sipora Utara	13
5747	Sir-Sir	81
5748	Sirah Pulau Padang	16
5749	Sirampog	33
5750	Sirandorung	12
5751	Sirapit (Serapit)	12
5752	Sirenja	72
5753	Sirets	93
5754	Sirimau	81
5755	Siritaun Wida Timur	81
5756	Siriwo	94
5757	Sirombu	12
5758	Sitahuis	12
5759	Sitelu Tali Urang Jehe (Sitellu)	12
5760	Sitelu Tali Urang Julu (Sitellu)	12
5761	Sitinjau Laut	15
5762	Sitinjo	12
5763	Sitio-tio	12
5764	Sitiung	13
5765	Sitolu Ori	12
5766	Situbondo	35
5767	Situjuah Limo Nagari (Situjuah Lima Nagari)	13
5768	Situraja	32
5769	Siulak	15
5770	Siulak Mukai	15
5771	Siwalalat	81
5772	Siwalan	33
5773	Skanto	91
5774	Slahung	35
5775	Slawi	33
5776	Sleman	34
5777	Sliyeg	32
5778	Slogohimo	33
5779	Sluke	33
5780	Soa	53
5781	Soba	95
5782	Sobaham	95
5783	Sobang	36
5784	Socah	35
5785	Sodonghilir	32
5786	Sogae'adu (Sogae Adu / Sogaeadu)	12
5787	Sojol	72
5788	Sojol Utara	72
5789	Sokan	61
5790	Sokaraja	33
5791	Soko	35
5792	Sokobanah	35
5793	Solear	36
5794	Soloikma	95
5795	Solokanjeruk (Solokan Jeruk)	32
5796	Solokuro	35
5797	Solor Barat	53
5798	Solor Selatan	53
5799	Solor Timur	53
5800	Somagede	33
5801	Somambawa	12
5802	Somba Opu (Upu)	73
5803	Somolo-Molo (Samolo)	12
5804	Sompak	61
5805	Sonder	71
5806	Songgom	33
5807	Songgon	35
5808	Sooko	35
5809	Sopai	73
5810	Soppeng Riaja	73
5811	Sor Ep	93
5812	Sorawolio (Sora Walio / Sorowalio)	74
5813	Soreang	73
5814	Soreang	32
5815	Sorkam	12
5816	Sorkam Barat	12
5817	Soromandi	52
5818	Sorong	92
5819	Sorong Barat	92
5820	Sorong Kepulauan	92
5821	Sorong Kota	92
5822	Sorong Manoi	92
5823	Sorong Timur	92
5824	Sorong Utara	92
5825	Soropia	74
5826	Sosa	12
5827	Sosa Julu	12
5828	Sosa Timur	12
5829	Sosoh Buay Rayap	16
5830	Sosopan	12
5831	Sosorgadong (Sosor Gadong)	12
5832	Sota	93
5833	Soug Jaya	92
5834	Soyo Jaya	72
5835	Soyoi Mambai	91
5836	Sragen	33
5837	Sragi	33
5838	Sragi	18
5839	Srandakan	34
5840	Srengat	35
5841	Sreseh	35
5842	Srono	35
5843	Srumbung	33
5844	Sruweng	33
5845	Stabat	12
5846	STL Ulu Terawas (Suku Tengah Lakitan Ulu Terawas)	16
5847	STM Hilir (Sinembah Tanjung Muda Hilir)	12
5848	STM Hulu (Sinembah Tanjung Muda Hulu)	12
5849	Suak Midai	21
5850	Suak Tapeh	16
5851	Suator	93
5852	Subah	33
5853	Subah	61
5854	Subang	32
5855	Subi	21
5856	Suboh	35
5857	Subur	93
5858	Sucinaraja	32
5859	Sudimoro	35
5860	Sugapa	94
5861	Sugie Besar	21
5862	Sugihwaras	35
5863	Sugio	35
5864	Suhaid	61
5865	Suka Bangun	12
5866	Suka Karya (Sukakarya)	16
5867	Suka Makmue	11
5868	Sukabumi	18
5869	Sukabumi	32
5870	Sukadana	32
5871	Sukadana	61
5872	Sukadana	18
5873	Sukadiri	36
5874	Sukagumiwang	32
5875	Sukahaji	32
5876	Sukahening	32
5877	Sukajadi	32
5878	Sukajadi	14
5879	Sukajaya	32
5880	Sukajaya	11
5881	Sukakarya	11
5882	Sukakarya	32
5883	Sukalarang	32
5884	Sukaluyu	32
5885	Sukamaju	73
5886	Sukamaju Selatan	73
5887	Sukamakmue	11
5888	Sukamakmur	32
5889	Sukamakmur (Suka Makmur)	11
5890	Sukamantri	32
5891	Sukamara	62
5892	Sukamerindu	16
5893	Sukamulia	52
5894	Sukamulya	36
5895	Sukanagara	32
5896	Sukapura	35
5897	Sukaraja	17
5898	Sukaraja	32
5899	Sukarame	32
5900	Sukarame	18
5901	Sukarami	16
5902	Sukaratu	32
5903	Sukaresik	32
5904	Sukaresmi	36
5905	Sukaresmi	32
5906	Sukasada	51
5907	Sukasari	32
5908	Sukatani	32
5909	Sukau	18
5910	Sukawangi	32
5911	Sukawati	51
5912	Sukawening	32
5913	Sukikai Selatan	94
5914	Sukmajaya	32
5915	Sukodadi	35
5916	Sukodono	35
5917	Sukodono	33
5918	Sukoharjo	33
5919	Sukoharjo	18
5920	Sukolilo	35
5921	Sukolilo	33
5922	Sukomanunggal	35
5923	Sukomoro	35
5924	Sukorambi	35
5925	Sukorame	35
5926	Sukorejo	35
5927	Sukorejo	33
5928	Sukosari	35
5929	Sukosewu	35
5930	Sukowono	35
5931	Sukra	32
5932	Sukun	35
5933	Sulabesi Barat	82
5934	Sulabesi Selatan	82
5935	Sulabesi Tengah	82
5936	Sulabesi Timur	82
5937	Sulamu	53
5938	Sulang	33
5939	Suli	73
5940	Suli Barat	73
5941	Suliki	13
5942	Suling Tambun	62
5943	Sultan Daulat	11
5944	Suluun Tareran	71
5945	Sumalata	75
5946	Sumalata Timur	75
5947	Sumarorong	76
5948	Sumay	15
5949	Sumbang	33
5950	Sumbawa	52
5951	Sumber	35
5952	Sumber	32
5953	Sumber	33
5954	Sumber Barito	62
5955	Sumber Harta	16
5956	Sumber Jaya	18
5957	Sumber Marga Telang	16
5958	Sumberasih	35
5959	Sumberbaru (Sumber Baru)	35
5960	Sumberejo	35
5961	Sumberejo (Sumber Rejo)	18
5962	Sumbergempol	35
5963	Sumberjambe (Sumber Jambe)	35
5964	Sumberjaya	32
5965	Sumberlawang	33
5966	Sumbermalang	35
5967	Sumbermanjing Wetan	35
5968	Sumberpucung	35
5969	Sumbersari (Sumber Sari)	35
5970	Sumbersuko	35
5971	Sumberwringin (Sumber Wringin)	35
5972	Sumbul	12
5973	Sumedang Selatan	32
5974	Sumedang Utara	32
5975	Sumo	95
5976	Sumobito	35
5977	Sumowono	33
5978	Sumpiuh	33
5979	Sumpur Kudus	13
5980	Sumur	36
5981	Sumur Bandung	32
5982	Sumuri (Simuri)	92
5983	Sungai Ambawang	61
5984	Sungai Apit	14
5985	Sungai Are	16
5986	Sungai Aur (Sungaiaur)	13
5987	Sungai Babuat	62
5988	Sungai Bahar	15
5989	Sungai Batang	14
5990	Sungai Beremas (Sei Beremas)	13
5991	Sungai Betung	61
5992	Sungai Boh	65
5993	Sungai Bungkal	15
5994	Sungai Durian (Sungaidurian)	63
5995	Sungai Garingging	13
5996	Sungai Gelam	15
5997	Sungai Kakap	61
5998	Sungai Kanan (Sei)	12
5999	Sungai Keruh	16
6000	Sungai Kunjang	64
6001	Sungai Kunyit	61
6002	Sungai Lala	14
6003	Sungai Laur	61
6004	Sungai Lilin	16
6005	Sungai Limau	13
6006	Sungai Loban	63
6007	Sungai Manau	15
6008	Sungai Mandau	14
6009	Sungai Mas	11
6010	Sungai Melayu Rayak	61
6011	Sungai Menang	16
6012	Sungai Pagu	13
6013	Sungai Pandan	63
6014	Sungai Penuh	15
6015	Sungai Pinang	63
6016	Sungai Pinang	64
6017	Sungai Pinang	16
6018	Sungai Pinyuh (Sei Pinyuh)	61
6019	Sungai Pua (Puar)	13
6020	Sungai Raya	61
6021	Sungai Raya	63
6022	Sungai Raya	11
6023	Sungai Raya Kepulauan	61
6024	Sungai Rotan	16
6025	Sungai Rumbai	13
6026	Sungai Rumbai	17
6027	Sungai Selan	19
6028	Sungai Sembilan	14
6029	Sungai Serut	17
6030	Sungai Tabuk	63
6031	Sungai Tabukan	63
6032	Sungai Tarab	13
6033	Sungai Tebelian	61
6034	Sungai Tubu	65
6035	Sungailiat (Sungai Liat)	19
6036	Sungayang	13
6037	Sunggal	12
6038	Sungkai Barat	18
6039	Sungkai Jaya	18
6040	Sungkai Selatan	18
6041	Sungkai Tengah	18
6042	Sungkai Utara	18
6043	Sunook	92
6044	Suntamon	95
6045	Suoh	18
6046	Supiori Barat	91
6047	Supiori Selatan	91
6048	Supiori Timur	91
6049	Supiori Utara	91
6050	Supnin	92
6051	Suppa	73
6052	Suradadi (Surodadi)	33
6053	Surade	32
6054	Suralaga	52
6055	Suranenggala	32
6056	Surian	32
6057	Suro Makmur	11
6058	Suru Suru	95
6059	Suru-suru	93
6060	Suruh	35
6061	Suruh	33
6062	Sururey	92
6063	Susoh	11
6064	Susua	12
6065	Susukan	32
6066	Susukan	33
6067	Susukan Lebak	32
6068	Susut	51
6069	Sutera	13
6070	Suti Semarang	61
6071	Sutojayan	35
6072	Suwawa	75
6073	Suwawa Selatan	75
6074	Suwawa Tengah	75
6075	Suwawa Timur	75
6076	Suwela (Suela)	52
6077	Swandiwe	91
6078	Syahcame	93
6079	Syamtalira Aron	11
6080	Syamtalira Bayu	11
6081	Syiah Kuala	11
6082	Syiah Utama	11
6083	Syujak	92
6084	T. Jambo Aye (Tanah Jambo Aye)	11
6085	Taba Penanjung	17
6086	Tabalar	64
6087	Tabanan	51
6088	Tabang	76
6089	Tabang	64
6090	Tabir	15
6091	Tabir Barat	15
6092	Tabir Ilir	15
6093	Tabir Lintas	15
6094	Tabir Selatan	15
6095	Tabir Timur	15
6096	Tabir Ulu	15
6097	Tabona	82
6098	Tabongo	75
6099	Tabonji	93
6100	Tabukan	63
6101	Tabukan Selatan	71
6102	Tabukan Selatan Tengah	71
6103	Tabukan Selatan Tenggara	71
6104	Tabukan Tengah	71
6105	Tabukan Utara	71
6106	Tabulahan	76
6107	Tabundung	53
6108	Tabunganen	63
6109	Tadu Raya	11
6110	Taebenu	53
6111	Taelarek	95
6112	Taganombak	94
6113	Tagime	95
6114	Tagineri	95
6115	Tagulandang	71
6116	Tagulandang Selatan	71
6117	Tagulandang Utara	71
6118	Tahota	92
6119	Tahuna	71
6120	Tahuna Barat	71
6121	Tahuna Timur	71
6122	Tahunan	33
6123	Taige	92
6124	Tajinan	35
6125	Tajurhalang	32
6126	Taka Bonerate (Takabonerate)	73
6127	Takari	53
6128	Takeran	35
6129	Takisung	63
6130	Takkalalla	73
6131	Takokak	32
6132	Taktakan	36
6133	Talaga	32
6134	Talaga Jaya (Telaga Jaya)	75
6135	Talaga Raya	74
6136	Talamau	13
6137	Talambo	95
6138	Talang	33
6139	Talang Empat	17
6140	Talang Kelapa	16
6141	Talang Muandau	14
6142	Talang Padang	16
6143	Talang Padang	18
6144	Talang Ubi	16
6145	Talango	35
6146	Talatako	72
6147	Talawaan	71
6148	Talawi	13
6149	Talawi	12
6150	Talegong	32
6151	Taliabu Barat	82
6152	Taliabu Barat Laut	82
6153	Taliabu Selatan	82
6154	Taliabu Timur	82
6155	Taliabu Timur Selatan	82
6156	Taliabu Utara	82
6157	Talibura	53
6158	Talisayan	64
6159	Taliwang	52
6160	Tallo	73
6161	Tallunglipu	73
6162	Talo	17
6163	Talo Kecil	17
6164	Taluditi (Taluduti)	75
6165	Talun	33
6166	Talun	35
6167	Talun (Cirebon Selatan)	32
6168	Tamako	71
6169	Tamalanrea	73
6170	Tamalate	73
6171	Tamalatea	73
6172	Taman	33
6173	Taman	35
6174	Taman Krocok	35
6175	Taman Rajo	15
6176	Taman Sari	19
6177	Taman Sari	31
6178	Tamanan	35
6179	Tamansari	32
6180	Tamansari	33
6181	Tambak	33
6182	Tambak	35
6183	Tambakboyo	35
6184	Tambakdahan	32
6185	Tambakrejo	35
6186	Tambakromo	33
6187	Tambaksari	32
6188	Tambaksari	35
6189	Tamban	63
6190	Tamban Catur	62
6191	Tambang	14
6192	Tambang Ulang	63
6193	Tambangan	12
6194	Tambelan	21
6195	Tambelang	32
6196	Tambelangan	35
6197	Tambora	52
6198	Tambora	31
6199	Tambun Selatan	32
6200	Tambun Utara	32
6201	Tambusai	14
6202	Tambusai Utara	14
6203	Tamiang Hulu	11
6204	Tammerodo Sendana (Tammeredo Sendana)	76
6205	Tampahan	12
6206	Tampaksiring (Tampak Siring)	51
6207	Tampan' Amma (Tampan Amma)	71
6208	Tana Lia	65
6209	Tana Lili	73
6210	Tana Righu	53
6211	Tana Wawo	53
6212	Tanah Abang	16
6213	Tanah Abang	31
6214	Tanah Cogok	15
6215	Tanah Grogot	64
6216	Tanah Jawa	12
6217	Tanah Kampung	15
6218	Tanah Luas	11
6219	Tanah Masa	12
6220	Tanah Merah	35
6221	Tanah Merah	14
6222	Tanah Miring	93
6223	Tanah Pasir	11
6224	Tanah Pinem	12
6225	Tanah Pinoh	61
6226	Tanah Pinoh Barat	61
6227	Tanah Putih	14
6228	Tanah Putih Tanjung Melawan	14
6229	Tanah Rubuh	92
6230	Tanah Sareal (Tanah Sereal)	32
6231	Tanah Sepenggal	15
6232	Tanah Sepenggal Lintas	15
6233	Tanah Siang	62
6234	Tanah Siang Selatan	62
6235	Tanah Tumbuh	15
6236	Tanambulava	72
6237	Tanantovea	72
6238	Tanara	36
6239	Tanasitolo	73
6240	Tandes	35
6241	Tanduk Kalua	76
6242	Tandun	14
6243	Tanete Riaja	73
6244	Tanete Riattang	73
6245	Tanete Riattang Barat	73
6246	Tanete Riattang Timur	73
6247	Tanete Rilau	73
6248	Tangan-Tangan	11
6249	Tangaran	61
6250	Tangen	33
6251	Tangerang	36
6252	Tanggetada	74
6253	Tanggeung	32
6254	Tanggul	35
6255	Tanggulangin	35
6256	Tanggunggunung (Tanggung Gunung)	35
6257	Tanggungharjo	33
6258	Tangma	95
6259	Tangse	11
6260	Tanimbar Selatan	81
6261	Tanimbar Utara	81
6262	Taniwel	81
6263	Taniwel Timur	81
6264	Tanjuang Baru (Tanjung Baru)	13
6265	Tanjung	63
6266	Tanjung	52
6267	Tanjung	33
6268	Tanjung Agung	16
6269	Tanjung Agung Palik	17
6270	Tanjung Balai	12
6271	Tanjung Batu	16
6272	Tanjung Beringin	12
6273	Tanjung Bintang	18
6274	Tanjung Bumi (Tanjungbumi)	35
6275	Tanjung Bunga	53
6276	Tanjung Emas	13
6277	Tanjung Gadang	13
6278	Tanjung Harapan	13
6279	Tanjung Harapan	64
6280	Tanjung Kemuning	17
6281	Tanjung Lago	16
6282	Tanjung Lubuk	16
6283	Tanjung Medan	14
6284	Tanjung Morawa	12
6285	Tanjung Mutiara	13
6286	Tanjung Palas	65
6287	Tanjung Palas Barat	65
6288	Tanjung Palas Tengah	65
6289	Tanjung Palas Timur	65
6290	Tanjung Palas Utara	65
6291	Tanjung Pandan	19
6292	Tanjung Pinang Barat	21
6293	Tanjung Pinang Kota	21
6294	Tanjung Pinang Timur	21
6295	Tanjung Priok	31
6296	Tanjung Pura (Tanjungpura)	12
6297	Tanjung Raja	18
6298	Tanjung Raja	16
6299	Tanjung Raya	13
6300	Tanjung Raya	18
6301	Tanjung Redeb	64
6302	Tanjung Sakti Pumi	16
6303	Tanjung Sari	18
6304	Tanjung Selor	65
6305	Tanjung Senang	18
6306	Tanjung Tiram	12
6307	Tanjunganom	35
6308	Tanjungbalai Selatan (Tanjung Balai Selatan)	12
6309	Tanjungbalai Utara (Tanjung Balai Utara)	12
6310	Tanjungjaya	32
6311	Tanjungkarang Barat (Tanjung Karang Barat)	18
6312	Tanjungkarang Pusat (Tanjung Karang Pusat)	18
6313	Tanjungkarang Timur (Tanjung Karang Timur)	18
6314	Tanjungkerta	32
6315	Tanjungmedar	32
6316	Tanjungsakti Pumu (Tanjung Sakti Pumu)	16
6317	Tanjungsari	34
6318	Tanjungsari	32
6319	Tanjungsiang	32
6320	Tanjungtebat (Tanjung Tebat)	16
6321	Tano Tombangan Angkola	12
6322	Tanoh Alas (Tanah Alas)	11
6323	Tanon	33
6324	Tanralili	73
6325	Tanta	63
6326	Taopa	72
6327	Tapa	75
6328	Tapaktuan (Tapak Tuan)	11
6329	Tapalang	76
6330	Tapalang Barat	76
6331	Tapango	76
6332	Tapen	35
6333	Tapian Dolok	12
6334	Tapian Nauli	12
6335	Tapin Selatan	63
6336	Tapin Tengah	63
6337	Tapin Utara	63
6338	Tapos	32
6339	Tapung	14
6340	Tapung Hilir	14
6341	Tapung Hulu	14
6342	Tarabintang (Tara Bintang)	12
6343	Taraju	32
6344	Tarakan Barat	65
6345	Tarakan Tengah	65
6346	Tarakan Timur	65
6347	Tarakan Utara	65
6348	Tarano	52
6349	Tareran	71
6350	Tarik	35
6351	Tarogong Kaler	32
6352	Tarogong Kidul	32
6353	Tarokan	35
6354	Tarowang	73
6355	Tarub	33
6356	Tarumajaya	32
6357	Tarup	95
6358	Tarutung	12
6359	Tasifeto Barat	53
6360	Tasifeto Timur	53
6361	Tasik Payawan	62
6362	Tasik Putri Puyu	14
6363	Tasikmadu	33
6364	Tatah Makmur	63
6365	Tatanga	72
6366	Tatapaan	71
6367	Tatoareng	71
6368	Tawaeli	72
6369	Tawalian	76
6370	Tawang	32
6371	Tawangharjo	33
6372	Tawangmangu	33
6373	Tawangsari	33
6374	Tayan Hilir	61
6375	Tayan Hulu	61
6376	Tayando Tam	81
6377	Tayu	33
6378	Tebas	61
6379	Tebat Karai	17
6380	Tebet	31
6381	Tebing	21
6382	Tebing Syahbandar	12
6383	Tebing Tinggi	12
6384	Tebing Tinggi	15
6385	Tebing Tinggi	14
6386	Tebing Tinggi	16
6387	Tebing Tinggi	63
6388	Tebing Tinggi Barat	14
6389	Tebing Tinggi Kota	12
6390	Tebing Tinggi Timur	14
6391	Tebo Ilir	15
6392	Tebo Tengah	15
6393	Tebo Ulu	15
6394	Tegal Barat	33
6395	Tegal Selatan	33
6396	Tegal Timur	33
6397	Tegalampel	35
6398	Tegalbuleud	32
6399	Tegaldlimo	35
6400	Tegallalang	51
6401	Tegalombo	35
6402	Tegalrejo	34
6403	Tegalrejo	33
6404	Tegalsari	35
6405	Tegalsiwalan (Tegal Siwalan)	35
6406	Tegalwaru	32
6407	Tegalwaru (Tegal Waru)	32
6408	Tegineneng	18
6409	Tegowanu	33
6410	Tehoru	81
6411	Teiraplu	95
6412	Tejakula	51
6413	Tekarang	61
6414	Tekung	35
6415	Telaga	75
6416	Telaga Antang	62
6417	Telaga Bauntung	63
6418	Telaga Biru	75
6419	Telaga Langsat	63
6420	Telagasari (Talagasari)	32
6421	Telanaipura	15
6422	Telawang	62
6423	Telen	64
6424	Telenggeme	95
6425	Tellu Limpoe	73
6426	Tellu Siattinge	73
6427	Tellulimpoe (Tellu Limpoe)	73
6428	Telluwanua	73
6429	Telok Sebong (Teluk Sebong)	21
6430	Teluk Ambon	81
6431	Teluk Ampimoi	91
6432	Teluk Arguni Atas	92
6433	Teluk Arguni Bawah (Yerusi)	92
6434	Teluk Batang	61
6435	Teluk Bayur	64
6436	Teluk Belengkong	14
6437	Teluk Bintan	21
6438	Teluk Dalam	11
6439	Teluk Dalam	12
6440	Teluk Deya	94
6441	Teluk Duairi	92
6442	Teluk Elpaputih	81
6443	Teluk Etna	92
6444	Teluk Gelam	16
6445	Teluk Kaiely	81
6446	Teluk Kepayang	63
6447	Teluk Keramat	61
6448	Teluk Kimi	94
6449	Teluk Mayalibit	92
6450	Teluk Mengkudu	12
6451	Teluk Meranti	14
6452	Teluk Mutiara	53
6453	Teluk Nibung	12
6454	Teluk Pakedai	61
6455	Teluk Pandan	64
6456	Teluk Pandan	18
6457	Teluk Patipi	92
6458	Teluk Sampit	62
6459	Teluk Segara	17
6460	Teluk Umar	94
6461	Teluk Waru	81
6462	Telukbetung Barat	18
6463	Telukbetung Selatan	18
6464	Telukbetung Timur	18
6465	Telukbetung Utara	18
6466	Telukjambe Barat	32
6467	Telukjambe Timur	32
6468	Teluknaga	36
6469	Telutih	81
6470	Temanggung	33
6471	Temayang	35
6472	Tembagapura	94
6473	Tembalang	33
6474	Tembarak	33
6475	Tembelang	35
6476	Tembilahan	14
6477	Tembilahan Hulu	14
6478	Tembuku	51
6479	Tembuni	92
6480	Temiang Pesisir	21
6481	Teminabuan	92
6482	Temon	34
6483	Tempe	73
6484	Tempeh	35
6485	Tempel	34
6486	Tempilang	19
6487	Tempuling	14
6488	Tempunak	61
6489	Tempuran	32
6490	Tempuran	33
6491	Tempurejo	35
6492	Tempursari	35
6493	Tenayan Raya	14
6494	Tenga	71
6495	Tengah Ilir	15
6496	Tengah Tani	32
6497	Tengaran	33
6498	Tenggarang	35
6499	Tenggarong	64
6500	Tenggarong Seberang	64
6501	Tenggilis Mejoyo	35
6502	Tenggulun	11
6503	Tenjo	32
6504	Tenjolaya	32
6505	Teon Nila Serua	81
6506	Teor	81
6507	Tepus	34
6508	Teramang Jaya	17
6509	Terangun (Terangon)	11
6510	Terara	52
6511	Teras	33
6512	Teras Terunjam	17
6513	Terbanggi Besar	18
6514	Terentang	61
6515	Teriak	61
6516	Tering	64
6517	Teripe Jaya (Tripe Jaya)	11
6518	Terisi (Trisi)	32
6519	Ternate Barat	82
6520	Tersono	33
6521	Terusan Nunyai	18
6522	Testega	92
6523	Tetap (Muara Tetap)	17
6524	Teunom	11
6525	Teupah Barat	11
6526	Teupah Selatan	11
6527	Teupah Tengah	11
6528	Tewah	62
6529	Tewang Sangalang Garing	62
6530	Teweh Baru	62
6531	Teweh Selatan	62
6532	Teweh Tengah	62
6533	Teweh Timur	62
6534	Ti Zain	93
6535	Tiang Pumpung	15
6536	Tiang Pumpung Kepungut	16
6537	Tibawa	75
6538	Tidore	82
6539	Tidore Selatan	82
6540	Tidore Timur	82
6541	Tidore Utara	82
6542	Tiga Dihaji	16
6543	Tigabinanga (Tiga Binanga)	12
6544	Tigalingga (Tiga Lingga)	12
6545	Tiganderket	12
6546	Tigapanah (Tiga Panah)	12
6547	Tigaraksa	36
6548	Tigi	94
6549	Tigi Barat	94
6550	Tigi Timur	94
6551	Tigo Lurah	13
6552	Tigo Nagari	13
6553	Tikala	73
6554	Tikala	71
6555	Tikke Raya	76
6556	Tikung	35
6557	Tilamuta	75
6558	Tilango	75
6559	Tilatang Kamang	13
6560	Tiloan	72
6561	Tilongkabila	75
6562	Timang Gajah	11
6563	Timori	95
6564	Timpah	62
6565	Timpeh	13
6566	Tinada	12
6567	Tinambung	76
6568	Tinanggea	74
6569	Tinangkung	72
6570	Tinangkung Selatan	72
6571	Tinangkung Utara	72
6572	Tinggi Raja	12
6573	Tinggimoncong	73
6574	Tingginambut	94
6575	Tinggouw	92
6576	Tingkir	33
6577	Tinombo	72
6578	Tinombo Selatan	72
6579	Tinondo	74
6580	Tiom	95
6581	Tiom Ollo	95
6582	Tiomneri	95
6583	Tiplol Mayalibit	92
6584	Tirawuta	74
6585	Tiris	35
6586	Tiro/Truseb	11
6587	Tiroang	73
6588	Tirtajaya	32
6589	Tirtamulya	32
6590	Tirtayasa	36
6591	Tirto	33
6592	Tirtomoyo	33
6593	Tirtoyudo	35
6594	Titehena	53
6595	Titeue	11
6596	Tiumang	13
6597	Tiworo Kepulauan	74
6598	Tiworo Selatan	74
6599	Tiworo Tengah	74
6600	Tiworo Utara	74
6601	Tiwu	74
6602	Tlanakan	35
6603	Tlogomulyo	33
6604	Tlogosari	35
6605	Tlogowungu	33
6606	Toapaya	21
6607	Toari	74
6608	Toba	61
6609	Tobadak	76
6610	Tobelo	82
6611	Tobelo Barat	82
6612	Tobelo Selatan	82
6613	Tobelo Tengah	82
6614	Tobelo Timur	82
6615	Tobelo Utara	82
6616	Toboali	19
6617	Tobouw	92
6618	Tobu	53
6619	Todanan	33
6620	Togean	72
6621	Togo Binongko	74
6622	Toho	61
6623	Toianas	53
6624	Toili	72
6625	Toili Barat	72
6626	Tojo	72
6627	Tojo Barat	72
6628	Tolala	74
6629	Tolangohula	75
6630	Toli-Toli Utara (Tolitoli Utara)	72
6631	Tolinggula	75
6632	Toma	12
6633	Tomage	92
6634	Tombariri	71
6635	Tombariri Timur	71
6636	Tombatu	71
6637	Tombatu Timur	71
6638	Tombatu Utara	71
6639	Tombolopao (Tombolo Pao)	73
6640	Tombulu	71
6641	Tomia	74
6642	Tomia Timur	74
6643	Tomilito (Tomolito)	75
6644	Tomini	71
6645	Tomini	72
6646	Tommo	76
6647	Tomo	32
6648	Tomohon Barat	71
6649	Tomohon Selatan	71
6650	Tomohon Tengah	71
6651	Tomohon Timur	71
6652	Tomohon Utara	71
6653	Tomoni	73
6654	Tomoni Timur	73
6655	Tomor Birip	93
6656	Tomosiga	94
6657	Tompaso	71
6658	Tompaso Barat	71
6659	Tompaso Baru	71
6660	Tompo Bulu (Tompobulu)	73
6661	Tompobullu (Tompobulu)	73
6662	Tompobulu (Tompu Bulu)	73
6663	Tomu	92
6664	Tondano Barat	71
6665	Tondano Selatan	71
6666	Tondano Timur	71
6667	Tondano Utara	71
6668	Tondon	73
6669	Tondong Tallasa	73
6670	Tongas	35
6671	Tongauna	74
6672	Tongauna Utara	74
6673	Tongkuno	74
6674	Tongkuno Selatan	74
6675	Tonjong	33
6676	Tonra	73
6677	Tontonunu	74
6678	Topiyai	94
6679	Topos	17
6680	Topoyo	76
6681	Tor Atas	91
6682	Torere	94
6683	Torgamba	12
6684	Toribulu	72
6685	Torjun	35
6686	Toroh	33
6687	Torue	72
6688	Tosari	35
6689	Totikum (Totikung)	72
6690	Totikum Selatan	72
6691	Touluaan	71
6692	Touluaan Selatan	71
6693	Towe	91
6694	Towea	74
6695	Towuti	73
6696	Tragah	35
6697	Trangkil	33
6698	Trawas	35
6699	Trenggalek	35
6700	Tretep	33
6701	Trienggadeng	11
6702	Trikora	95
6703	Trimurjo	18
6704	Tripa Makmur	11
6705	Trowulan	35
6706	Trucuk	33
6707	Trucuk	35
6708	Trumon	11
6709	Trumon Tengah	11
6710	Trumon Timur	11
6711	Tuah Negeri	16
6712	Tuahmadani	14
6713	Tualan Hulu	62
6714	Tualang	14
6715	Tuban	35
6716	Tubang	93
6717	Tubei (Pelabai)	17
6718	Tubo Sendana	76
6719	Tugala Oyo	12
6720	Tugu	33
6721	Tugu	35
6722	Tugumulyo	16
6723	Tuhemberua	12
6724	Tuhiba	92
6725	Tujuh Belas	61
6726	Tukak Sadai	19
6727	Tukdana	32
6728	Tukka	12
6729	Tulakan	35
6730	Tulang Bawang Tengah	18
6731	Tulang Bawang Udik	18
6732	Tulangan	35
6733	Tulin Onsoi	65
6734	Tulis	33
6735	Tulung	33
6736	Tulung Selapan	16
6737	Tulungagung	35
6738	Tumbang Titi	61
6739	Tumijajar	18
6740	Tuminting (Tuminiting)	71
6741	Tumpaan	71
6742	Tumpang	35
6743	Tungkal Ilir	15
6744	Tungkal Ilir	16
6745	Tungkal Jaya	16
6746	Tungkal Ulu	15
6747	Tunjung Teja	36
6748	Tunjungan	33
6749	Tuntang	33
6750	Turatea	73
6751	Turen	35
6752	Turi	35
6753	Turi	34
6754	Turikale	73
6755	Tutar (Tubbi Taramanu, Tutallu)	76
6756	Tutuk Tolu	81
6757	Tutur	35
6758	Tutuyan	71
6759	Ubahak	95
6760	Ubalihi	95
6761	Ubud	51
6762	Udanawu	35
6763	Ueesi	74
6764	Uepai	74
6765	Ugimba	94
6766	Ujan Mas	16
6767	Ujan Mas	17
6768	Ujung	73
6769	Ujung Batu	14
6770	Ujung Batu	12
6771	Ujung Bulu	73
6772	Ujung Padang	12
6773	Ujung Pandang	73
6774	Ujung Tanah	73
6775	Ujungberung (Ujung Berung)	32
6776	Ujungjaya	32
6777	Ujungloe (Ujung Loe)	73
6778	Ujungpangkah (Ujung Pangkah)	35
6779	Ukha	95
6780	Ukui	14
6781	Ulakan Tapakih	13
6782	Ulaweng	73
6783	Ulee Kareng	11
6784	Ulilin	93
6785	Ulim	11
6786	Ulok Kupai	17
6787	Ulu Barumun	12
6788	Ulu Belu (Ulubelu)	18
6789	Ulu Idanotae	12
6790	Ulu Manna	17
6791	Ulu Moro'o (Ulu Narwo)	12
6792	Ulu Musi	16
6793	Ulu Ogan	16
6794	Ulu Pungkut	12
6795	Ulu Rawas	16
6796	Ulu Sosa	12
6797	Ulu Talo	17
6798	Uluan	12
6799	Ulubongka	72
6800	Uluere	73
6801	Ulugawo	12
6802	Uluiwoi	74
6803	Ulujadi	72
6804	Ulujami	33
6805	Ulumanda (Ulumunda)	76
6806	Ulunoyo	12
6807	Ulususua	12
6808	Umagi	95
6809	Umalulu	53
6810	Umbu Ratu Nggay	53
6811	Umbu Ratu Nggay Barat	53
6812	Umbu Ratu Nggay Tengah	53
6813	Umbulharjo	34
6814	Umbulsari	35
6815	Umbunasi	12
6816	Umpu Semenguk	18
6817	Una Una	72
6818	Unaaha	74
6819	Undaan	33
6820	Ungar	21
6821	Ungaran Barat	33
6822	Ungaran Timur	33
6823	Unir Sirau	93
6824	Unter Iwes (Unterwiris)	52
6825	Unurum Guay	91
6826	Upau	63
6827	Uram Jaya	17
6828	Urei Faisei	91
6829	Usilimo	95
6830	Utan	52
6831	Uut Murung	62
6832	Uwapa	94
6833	V Koto	17
6834	V Koto Kampung Dalam	13
6835	V Koto Timur	13
6836	Venaha	93
6837	VII Koto	15
6838	VII Koto Ilir	15
6839	VII Koto Sungai Sarik	13
6840	Waan	93
6841	Wabula	74
6842	Wadaga	74
6843	Wadangku	95
6844	Wadaslintang	33
6845	Wado	32
6846	Wae Rii	53
6847	Waeapo	81
6848	Waegi	94
6849	Waelata	81
6850	Waesama	81
6851	Wagir	35
6852	Waiblama	53
6853	Waibu	91
6854	Waigeo Barat	92
6855	Waigeo Barat Kepulauan	92
6856	Waigeo Selatan	92
6857	Waigeo Timur	92
6858	Waigeo Utara	92
6859	Waigete	53
6860	Wajak	35
6861	Wajo	73
6862	Wakate	81
6863	Wakorumba Selatan	74
6864	Wakorumba Utara	74
6865	Wakuwo	95
6866	Walaik	95
6867	Walantaka	36
6868	Walea Besar	72
6869	Walea Kepulauan	72
6870	Waled	32
6871	Walelagama	95
6872	Walenrang	73
6873	Walenrang Barat	73
6874	Walenrang Timur	73
6875	Walenrang Utara	73
6876	Walma	95
6877	Waluran	32
6878	Wame	95
6879	Wamena	95
6880	Wamesa	92
6881	Wamesa (Idoor)	92
6882	Wampu	12
6883	Wanadadi (Wonodadi)	33
6884	Wanaraja	32
6885	Wanaraya	63
6886	Wanareja	33
6887	Wanasaba	52
6888	Wanasalam	36
6889	Wanasari	33
6890	Wanayasa	32
6891	Wanayasa	33
6892	Wandai	94
6893	Wanea	71
6894	Wangbe	94
6895	Wanggar	94
6896	Wanggarasi	75
6897	Wangi Wangi	74
6898	Wangi Wangi Selatan	74
6899	Wangon	33
6900	Wania	94
6901	Wano Barat	95
6902	Wanokaka	53
6903	Wanwi	94
6904	Waplau	81
6905	Wapoga	94
6906	Wapoga	91
6907	Wara	73
6908	Wara Barat	73
6909	Wara Selatan	73
6910	Wara Timur	73
6911	Wara Utara	73
6912	Wari/Taiyeve II	95
6913	Waringinkurung (Waringin Kurung)	36
6914	Waris	91
6915	Warkuk Ranau Selatan	16
6916	Warmare	92
6917	Waropen Atas	91
6918	Waropen Bawah	91
6919	Waropko	93
6920	Warsa	91
6921	Wartutin	92
6922	Waru	35
6923	Waru	64
6924	Warudoyong	32
6925	Warungasem	33
6926	Warunggunung	36
6927	Warungkiara	32
6928	Warungkondang	32
6929	Warungpring	33
6930	Warureja (Warurejo)	33
6931	Warwarbomi	92
6932	Wasile	82
6933	Wasile Selatan	82
6934	Wasile Tengah	82
6935	Wasile Timur	82
6936	Wasile Utara	82
6937	Wasior	92
6938	Wasuponda	73
6939	Watang Pulu	73
6940	Watang Sawitto (Watang Sawito)	73
6941	Watang Sidenreng (Wattang Sidenreng)	73
6942	Wates	34
6943	Wates	35
6944	Watopute	74
6945	Watubangga	74
6946	Watukumpul	33
6947	Watulimo	35
6948	Watumalang	33
6949	Watunohu	74
6950	Waway Karya	18
6951	Wawo	52
6952	Wawo	74
6953	Wawolesea	74
6954	Wawonii Barat	74
6955	Wawonii Selatan	74
6956	Wawonii Tengah	74
6957	Wawonii Tenggara	74
6958	Wawonii Timur	74
6959	Wawonii Timur Laut	74
6960	Wawonii Utara	74
6961	Wawotobi	74
6962	Way Bungur	18
6963	Way Halim	18
6964	Way Jepara	18
6965	Way Kenanga	18
6966	Way Khilau	18
6967	Way Krui	18
6968	Way Lima	18
6969	Way Pangubuan	18
6970	Way Panji	18
6971	Way Ratai	18
6972	Way Seputih	18
6973	Way Serdang	18
6974	Way Sulan	18
6975	Way Tenong	18
6976	Way Tuba	18
6977	Wayer	92
6978	Web	91
6979	Weda	82
6980	Weda Selatan	82
6981	Weda Tengah	82
6982	Weda Timur	82
6983	Weda Utara	82
6984	Wedarijaksa	33
6985	Wedi	33
6986	Wedung	33
6987	Wegee Bino	94
6988	Wegee Muka	94
6989	Weime	95
6990	Welahan	33
6991	Welak	53
6992	Welarek	95
6993	Weleri	33
6994	Welesi	95
6995	Weliman	53
6996	Wemak	92
6997	Wenam	95
6998	Wenang	71
6999	Wera	52
7000	Wereka	95
7001	Weriagar	92
7002	Werima	95
7003	Werinama	81
7004	Wermaktian (Wer Maktian)	81
7005	Wertamrian (Wer Tamrian)	81
7006	Weru	32
7007	Weru	33
7008	Wesaput	95
7009	Wetar Barat	81
7010	Wetar Selatan	81
7011	Wetar Timur	81
7012	Wetar Utara	81
7013	Wewaria	53
7014	Wewewa Barat	53
7015	Wewewa Selatan	53
7016	Wewewa Tengah	53
7017	Wewewa Timur	53
7018	Wewewa Utara	53
7019	Wewiku	53
7020	Widang	35
7021	Widasari	32
7022	Widodaren	35
7023	Wih Pesam	11
7024	Wilangan	35
7025	Wilhem Roumbouts	92
7026	Wina	95
7027	Windesi	91
7028	Windesi	92
7029	Windusari	33
7030	Winong	33
7031	Winongan	35
7032	Wiradesa	33
7033	Wiringgambut	95
7034	Wirobrajan	34
7035	Wirosari	33
7036	Wita Ponda	72
7037	Wita Waya	95
7038	Witihama	53
7039	Wiwirano	74
7040	Wiyung	35
7041	Wlingi	35
7042	Woha	52
7043	Woja	52
7044	Wolasi	74
7045	Wolio	74
7046	Wolo	95
7047	Wolo	74
7048	Wolojita	53
7049	Wolomeze (Riung Selatan)	53
7050	Wolowa	74
7051	Wolowae	53
7052	Wolowaru	53
7053	Wonawa	91
7054	Wondiboy	92
7055	Wonggeduku	74
7056	Wonggeduku Barat	74
7057	Wongsorejo	35
7058	Woniki	95
7059	Wonoasih	35
7060	Wonoasri	35
7061	Wonoayu	35
7062	Wonoboyo	33
7063	Wonocolo	35
7064	Wonodadi	35
7065	Wonogiri	33
7066	Wonokerto	33
7067	Wonokromo	35
7068	Wonomerto	35
7069	Wonomulyo	76
7070	Wonopringgo	33
7071	Wonorejo	35
7072	Wonosalam	33
7073	Wonosalam	35
7074	Wonosamodro	33
7075	Wonosari	34
7076	Wonosari	75
7077	Wonosari	33
7078	Wonosari	35
7079	Wonosegoro	33
7080	Wonosobo	33
7081	Wonosobo	18
7082	Wonotirto	35
7083	Wonotunggal	33
7084	Wonti	91
7085	Wori	71
7086	Wosak	95
7087	Wotan Ulumando	53
7088	Wotu	73
7089	Wouma	95
7090	Woyla	11
7091	Woyla Barat	11
7092	Woyla Timur	11
7093	Wringin	35
7094	Wringinanom (Wringin Anom)	35
7095	Wua-Wua	74
7096	Wuar Labobar	81
7097	Wugi	95
7098	Wulandoni	53
7099	Wulanggitang	53
7100	Wulla Waijelu (Wula Waijelu)	53
7101	Wuluhan	35
7102	Wundulako	74
7103	Wungu	35
7104	Wunim	95
7105	Wuryantoro	33
7106	Wusama	95
7107	Wusi	95
7108	Wutpaga	95
7109	X Koto (Sepuluh Koto)	13
7110	X Koto Diatas	13
7111	X Koto Singkarak	13
7112	XIII Koto Kampar	14
7113	XIV Koto	17
7114	Yaffi	91
7115	Yagai	94
7116	Yahuliambut	95
7117	Yakomi	93
7118	Yal	95
7119	Yalengga	95
7120	Yambi	94
7121	Yamo	94
7122	Yamoneri	94
7123	Yamor	92
7124	Yaniruma	93
7125	Yapen Barat	91
7126	Yapen Selatan	91
7127	Yapen Timur	91
7128	Yapen Utara	91
7129	Yapsi	91
7130	Yaro	94
7131	Yatamo	94
7132	Yaur	94
7133	Yawakukat	91
7134	Yawosi	91
7135	Yembun	92
7136	Yendidori	91
7137	Yenggelo	95
7138	Yigi	95
7139	Yiginua	95
7140	Yiluk	95
7141	Yogosem	95
7142	Yokari	91
7143	Yosowilangun	35
7144	Youtadi	94
7145	Yugumuak	94
7146	Yugungwi	95
7147	Yuko	95
7148	Yuneri	95
\.


--
-- Data for Name: faqs; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.faqs (id, question, answer, "order", category, is_active, created_at, updated_at, deleted_at) FROM stdin;
00000000-0000-0000-0000-000000000001	Apakah ada layanan Antar-Jemput?	Ya, tersedia layanan antar-jemput untuk area yang termasuk dalam cakupan layanan DriveMaster.	2	general	t	2026-06-11 14:14:58.896662+07	2026-07-16 09:42:42.513647+07	\N
13d5a284-721a-4b0f-aad0-f68b8c165853	Berapa lama satu sesi pertemuan?	Setiap sesi berlangsung selama 1 jam, maksimal 2 jam dengan instruktur profesional.	1	general	t	2026-06-13 14:57:14.901499+07	2026-07-16 09:42:52.297929+07	\N
850cbc4f-0ea3-4e7e-aa13-e3fc6d2e29b1	Apakah bisa kursus di Weekend?	Bisa, karena kami menyediakan layanan di weekend jam 08.00-17.00 sesuai ketersediaan jadwal.	4	general	t	2026-06-13 14:51:38.322204+07	2026-07-16 10:24:43.43594+07	\N
00000000-0000-0000-0000-000000000002	Apakah bisa kursus di malam hari?	Bisa, kami menyediakan layanan kursus di malam hari jam 18.00-20.00 sesuai ketersediaan jadwal.	3	general	t	2026-06-11 14:14:58.896662+07	2026-07-16 10:28:26.491849+07	\N
9bd03e34-6c08-4c09-bd15-611ac12ac7ec	Apakah bisa mengambil paket tambahan jika merasa belum cukup?	Tentu bisa. Peserta dapat menambah sesi latihan sesuai kebutuhan agar semakin percaya diri dalam berkendara.	0	general	t	2026-07-16 10:31:07.700222+07	2026-07-16 10:31:07.700222+07	\N
321712f5-0616-4f36-92f2-211bdf405c1f	Apakah ada sertifikat?	Ya. Peserta akan mendapatkan Sertifikat DriveMaster setelah menyelesaikan program kursus mengemudi	0	general	t	2026-07-16 10:33:03.455046+07	2026-07-16 10:33:03.455046+07	\N
\.


--
-- Data for Name: general_settings; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.general_settings (id, business_name, email, phone, fax, whats_app, address, hours_mon_fri, hours_sat_sun, hours_night_shift, promo_end_date, notify_email, notify_sms, notify_whats_app, created_at, updated_at, instagram, youtube, map_direction) FROM stdin;
00000000-0000-0000-0000-000000000001	Drive Master Indonesia	drivemaster.admin@gmail.com	+62 85286160029	+62 21 1234 5679	085286160029	The Smith Office, 9th Floor, Unit 0902 Jl. Jalur Sutera Timur, RT 002/003, Kunciran, Kec. Pinang, Kota Tangerang, Provinsi Banten 15144	08:00 - 17:00	08:00 - 17:00	18:00 - 20:00	\N	t	f	f	2026-06-05 11:29:39.285281+07	2026-07-16 15:39:16.387821+07	https://www.instagram.com/drivemasterkursus/		https://maps.app.goo.gl/qGngC2sF4G3jt8Vs8
\.


--
-- Data for Name: monthly_sales; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.monthly_sales (year, month, total_sales, total_revenue, total_discount, total_refunds, net_revenue, avg_order_value, canceled_sales, pending_sales, completed_sales, source_breakdown, created_at, updated_at) FROM stdin;
\.


--
-- Data for Name: package_benefits; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.package_benefits (id, package_id, title, description, icon, sort_order, created_at) FROM stdin;
42f6944e-34bb-4e0b-acd6-171848c019c2	49e284e0-193e-4707-b226-f7a4ed9f0fcd	4 Sessions with Instructor	Learn basics with experienced instructors	book-open	1	2026-06-01 10:47:42.551994+07
2b5cad2e-58e7-4ef4-b7aa-dfa173e4f5ae	49e284e0-193e-4707-b226-f7a4ed9f0fcd	Basic Traffic Rules	Comprehensive traffic rules and signs training	traffic-cone	2	2026-06-01 10:47:42.552+07
303ddd0f-c009-4a22-8a1c-a24aca7b1487	49e284e0-193e-4707-b226-f7a4ed9f0fcd	90-Minute Per Session	Adequate time for practical learning	clock	3	2026-06-01 10:47:42.552003+07
2f169a88-8be1-4a3c-a5c0-24b3b39e67d7	1ae7b1f4-2fbc-41e8-a473-cd1f9e900d5c	8 Sessions with Instructor	More practice time with professional instructors	users	1	2026-06-01 10:47:42.552013+07
eca478a2-8728-41c3-945e-8066b0591892	1ae7b1f4-2fbc-41e8-a473-cd1f9e900d5c	City Driving Skills	Master driving in urban environments	building	2	2026-06-01 10:47:42.552015+07
1ddd149e-d8b0-4002-acea-2d77ef862675	1ae7b1f4-2fbc-41e8-a473-cd1f9e900d5c	Parking Practice	Learn parallel parking and perpendicular parking	car	3	2026-06-01 10:47:42.552016+07
9f6675bc-7a70-48c9-8fef-044ce84682e2	1ae7b1f4-2fbc-41e8-a473-cd1f9e900d5c	Night Driving Training	Safe night driving techniques	moon	4	2026-06-01 10:47:42.55204+07
e1e94080-9a1b-46b9-965a-589517a5a568	5b0f238f-72b5-44eb-b975-197dd10549e7	10 Sessions with Instructor	Extensive practice with certified instructors	award	1	2026-06-01 10:47:42.552046+07
bcf09d75-e694-4644-a3cd-acf63927d5c1	5b0f238f-72b5-44eb-b975-197dd10549e7	Highway Driving	Learn safe highway driving techniques	road	2	2026-06-01 10:47:42.552047+07
0a838a9e-431d-4c55-a427-80f1e5ebd7a4	5b0f238f-72b5-44eb-b975-197dd10549e7	Defensive Driving	Master defensive driving techniques	shield	3	2026-06-01 10:47:42.552052+07
3645f534-ef19-48ac-8eee-3faf8bbd2bbe	5b0f238f-72b5-44eb-b975-197dd10549e7	Exam Preparation	Complete preparation for driving test	clipboard-check	4	2026-06-01 10:47:42.552053+07
a921609f-717c-4a2e-b3ef-b2cd91a33616	5b0f238f-72b5-44eb-b975-197dd10549e7	Free Study Materials	Access to digital learning resources	book	5	2026-06-01 10:47:42.552055+07
325fd715-6919-4441-9c16-39f1c74d1f92	156e94ec-0b53-4a7f-b552-a5e4e7756726	12 Sessions 1-on-1 Coaching	Personalized instruction with senior instructors	user-check	1	2026-06-01 10:47:42.552059+07
5c697478-5086-4ec8-af7e-5ccd3def409d	156e94ec-0b53-4a7f-b552-a5e4e7756726	All Driving Scenarios	City, highway, mountain, and night driving	map	2	2026-06-01 10:47:42.552061+07
c18e4ea2-91cb-4d6b-a4b1-29c13060658d	156e94ec-0b53-4a7f-b552-a5e4e7756726	Guaranteed Pass	If you don't pass, get free retake	check-circle	3	2026-06-01 10:47:42.552062+07
7d1caf70-0a4f-4d3c-985d-5b8922f1c865	156e94ec-0b53-4a7f-b552-a5e4e7756726	Priority Scheduling	Flexible booking with priority slots	calendar	4	2026-06-01 10:47:42.552064+07
ea64c417-80b0-4f44-8b79-25216382b675	156e94ec-0b53-4a7f-b552-a5e4e7756726	Exclusive Study Materials	Premium video lessons and mock tests	video	5	2026-06-01 10:47:42.552065+07
cadc18ce-00de-4297-8221-9b69e48e09c3	156e94ec-0b53-4a7f-b552-a5e4e7756726	Post-Training Support	3 months mentorship after completion	heart	6	2026-06-01 10:47:42.552067+07
\.


--
-- Data for Name: packages; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.packages (id, name, description, package_type, price, discount_price, duration_minutes, total_sessions, status, image_url, created_at, updated_at, features, highlight, student_count) FROM stdin;
11111111-1111-1111-1111-111111111201	8x	8 training sessions with SIM A certification	silver	2350000.00	2350000.00	60	8	active		2026-07-02 16:06:16.263037+07	2026-07-14 09:36:55.296981+07	{"Free Trial","8 training sessions","SIM A"}	f	0
11111111-1111-1111-1111-111111111301	10x	10 training sessions with SIM A certification	gold	2750000.00	2750000.00	60	10	active		2026-07-02 16:06:16.263038+07	2026-07-14 09:38:10.904816+07	{"Free Trial","10 training sessions","SIM A"}	f	0
11111111-1111-1111-1111-111111111402	12x + Night Session	12 training sessions including night driving with SIM A certification	platinum	3400000.00	3400000.00	60	12	inactive		2026-07-10 08:22:21.828559+07	2026-07-14 10:04:15.041961+07	{"Free Trial","12 training sessions","Night Session","SIM A"}	f	0
11111111-1111-1111-1111-111111111404	12x + Weekend & Night Session	12 training sessions including weekend and night sessions with SIM A certification	platinum	3400000.00	3400000.00	60	12	inactive		2026-07-10 08:22:21.828559+07	2026-07-14 10:04:33.114955+07	{"Free Trial","12 training sessions","Weekend Session","Night Session","SIM A"}	f	0
11111111-1111-1111-1111-111111111304	10x + Weekend & Night Session	10 training sessions including weekend and night sessions with SIM A certification	gold	2950000.00	2950000.00	60	10	inactive		2026-07-10 08:22:21.828558+07	2026-07-14 10:04:53.422061+07	{"Free Trial","10 training sessions","Weekend Session","Night Session","SIM A"}	f	0
11111111-1111-1111-1111-111111111204	8x + Weekend & Night Session	8 training sessions including weekend and night sessions with SIM A certification	silver	2250000.00	2050000.00	60	8	inactive		2026-07-10 08:22:21.828557+07	2026-07-10 10:03:56.110085+07	{"Free Trial","8 training sessions","Weekend Session","Night Session","SIM A"}	f	0
11111111-1111-1111-1111-111111111104	6x + Weekend & Night Session	6 training sessions including weekend and night sessions with SIM A certification	bronze	2250000.00	2250000.00	60	6	inactive		2026-07-10 08:22:21.828556+07	2026-07-14 10:05:35.407599+07	{"Free Trial","6 training sessions","Weekend Session","Night Session","SIM A"}	f	0
11111111-1111-1111-1111-111111111401	12x	12 training sessions with SIM A certification	platinum	3150000.00	3150000.00	60	12	active		2026-07-02 16:06:16.263039+07	2026-07-14 09:38:43.233884+07	{"Free Trial","12 training sessions","SIM A"}	f	0
11111111-1111-1111-1111-111111111101	6x	6 training sessions with SIM A certification	bronze	2150000.00	2150000.00	60	6	active		2026-07-10 08:22:21.828545+07	2026-07-30 08:27:47.112012+07	{"Free Trial","6 training sessions","SIM A"}	f	0
11111111-1111-1111-1111-111111111102	6x + Night Session	6 training sessions including night driving with SIM A certification	bronze	2250000.00	2250000.00	60	6	inactive		2026-07-10 08:22:21.828545+07	2026-07-14 09:59:55.381925+07	{"Free Trial","6 training sessions","Night Session","SIM A"}	f	0
11111111-1111-1111-1111-111111111103	6x + Weekend Session	6 training sessions including weekend sessions with SIM A certification	bronze	2250000.00	2250000.00	60	6	inactive		2026-07-10 08:22:21.828556+07	2026-07-14 10:00:15.232107+07	{"Free Trial","6 training sessions","Weekend Session","SIM A"}	f	0
11111111-1111-1111-1111-111111111203	8x + Weekend Session	8 training sessions including weekend sessions with SIM A certification	silver	2500000.00	2500000.00	60	8	inactive		2026-07-10 08:22:21.828557+07	2026-07-14 10:01:27.007446+07	{"Free Trial","8 training sessions","Weekend Session","SIM A"}	f	0
11111111-1111-1111-1111-111111111202	8x + Night Session	8 training sessions including night driving with SIM A certification	silver	2500000.00	2500000.00	60	8	inactive		2026-07-10 08:22:21.828556+07	2026-07-14 10:01:53.377062+07	{"Free Trial","8 training sessions","Night Session","SIM A"}	f	0
11111111-1111-1111-1111-111111111303	10x + Weekend Session	10 training sessions including weekend sessions with SIM A certification	gold	2950000.00	2950000.00	60	10	inactive		2026-07-10 08:22:21.828558+07	2026-07-14 10:02:42.905979+07	{"Free Trial","10 training sessions","Weekend Session","SIM A"}	f	0
11111111-1111-1111-1111-111111111302	10x + Night Session	10 training sessions including night driving with SIM A certification	gold	2950000.00	2950000.00	60	10	inactive		2026-07-10 08:22:21.828558+07	2026-07-14 10:02:51.64287+07	{"Free Trial","10 training sessions","Night Session","SIM A"}	f	0
11111111-1111-1111-1111-111111111403	12x + Weekend Session	12 training sessions including weekend sessions with SIM A certification	platinum	3400000.00	3400000.00	60	12	inactive		2026-07-10 08:22:21.828559+07	2026-07-14 10:04:00.024057+07	{"Free Trial","12 training sessions","Weekend Session","SIM A"}	f	0
4601ad1b-4dc9-4727-b229-43cdf6fddecb	Paket Percobaan	Dokuritsu Junbii Cosakai	bronze	10000.00	0.00	60	2	inactive		2026-07-30 18:23:54.953434+07	2026-08-04 11:03:49.925912+07	\N	f	0
\.


--
-- Data for Name: pages; Type: TABLE DATA; Schema: public; Owner: admin_drive
--

COPY public.pages (id, title, slug, status, sections, created_at, updated_at, deleted_at) FROM stdin;
00000000-0000-0000-0000-000000000002	About Us	/about	draft	[{"id":"2gfdwlez9","type":"hero","data":{"heading":"Langkah Maju, Belajar dari Masa Depan","subheading":"Drive Master bukan hanya tentang mengajar cara mengemudi, ini tentang mendefinisikan ulang standar pendidikan mengemudi di Indonesia. Sebagai pelopor sekolah mengemudi Kendaraan Listrik, kami percaya bahwa pengemudi masa depan harus lahir dari teknologi masa depan—modern, cerdas, dan ramah lingkungan.","ctaText":"Mulai Perjalanan Anda","ctaLink":"/auth/register","secondaryCtaText":"Hubungi Kami","secondaryCtaLink":"/contact","features":[{"title":"Pelopor EV Bebas Emisi","icon":"i-lucide-leaf"},{"title":"Instruktur Bersertifikat","icon":"i-lucide-award"}]}},{"id":"r8t6g2kmg","type":"specifications","data":{"headline":"Keselamatan Utama","title":"Prioritas Kami adalah Keselamatan Anda","description":"Dalam bisnis kursus mengemudi, keselamatan bukan hanya sebuah fitur; itu adalah fondasi inti kami.","items":[{"title":"Instruktur Bersertifikat","subtitle":"Instruktur kami adalah profesional berlisensi yang bersertifikat khusus untuk mengoperasikan kendaraan listrik premium.","icon":"i-lucide-award","description":[]},{"title":"Teknologi Keselamatan Aktif","subtitle":"Memanfaatkan fitur keselamatan bawaan EV seperti Collision Avoidance dan Blind Spot Monitoring untuk meminimalkan risiko.","icon":"i-lucide-radar","description":[]}]}},{"id":"xj64r6k4g","type":"quote","data":{"quote":"Visi kami bukan hanya untuk menghasilkan pengemudi yang bisa memutar kemudi, tetapi untuk membina pengemudi yang cerdas dan aman yang siap merangkul era elektrifikasi.","description":"Di Drive Master Indonesia, kami percaya bahwa cara kita belajar mengemudi harus berevolusi seiring dengan evolusi teknologi otomotif. Kami berkomitmen untuk menjadi standar baru dalam pendidikan mengemudi yang ramah lingkungan, memastikan bahwa setiap lulusan memiliki keterampilan mengemudi tingkat tinggi serta kesadaran akan masa depan mobilitas yang berkelanjutan.","ctaText":"Mulai Perjalanan Anda","ctaLink":"/auth/register","secondaryCtaText":"Chat WhatsApp","secondaryCtaLink":"https://wa.me/6285286160029?text=Halo%20Drive%20Master%2C%20saya%20ingin%20bertanya%20tentang%20kursus%20mengemudi"}}]	2026-07-10 09:49:21.414193+07	2026-08-05 23:38:35.078313+07	\N
00000000-0000-0000-0000-000000000003	Services	/services	published	[{"id":"hero-services","type":"cta","data":{"heading":"Layanan","description":"Kursus mengemudi komprehensif yang dirancang untuk masa depan listrik. Dari pemula hingga pengemudi mahir, kami memiliki program yang sempurna untuk Anda.","icon":"i-lucide-award","links":[{"label":"Lihat Paket","to":"/packages","color":"warning","icon":"i-lucide-package"},{"label":"Pesan Konsultasi","to":"https://wa.me/628119124848?text=Halo%20Drive%20Master%2C%20saya%20ingin%20bertanya%20tentang%20kursus%20mengemudi","color":"primary","variant":"outline","icon":"i-simple-icons-whatsapp"}],"secondaryButtonText":"Chat WhatsApp","secondaryButtonLink":"https://wa.me/6285286160029?text=Halo%20Drive%20Master%2C%20saya%20ingin%20bertanya%20tentang%20kursus%20mengemudi","secondaryButtonIcon":"i-simple-icons-whatsapp","buttonText":"Lihat Paket","buttonLink":"/packages","buttonIcon":"i-lucide-package"}},{"id":"specs-services","type":"specifications","data":{"headline":"Layanan Khusus","items":[{"title":"Layanan Khusus","subtitle":"Keuntungan memilih Kami","icon":"i-lucide-star","description":["SIM","Antar-jemput gratis (Alam Sutera, BSD, Gading Serpong)","Materi teori termasuk","Sertifikat penyelesaian"]},{"title":"Waktu Sesi","subtitle":"Waktu sesi asli kami","icon":"i-lucide-car","description":["Sesi di hari kerja buka dari jam 08:00 sampai 17:00","60 menit per sesi","Sesi Malam: 18:00 - 20:00","Sabtu - Minggu: 08:00 - 17:00"]},{"title":"Transmisi Mobil","subtitle":"Mobil kami menggunakan transmisi matic","icon":"i-lucide-car","description":["Transmisi Matic"]},{"title":"Sesi Malam","subtitle":"Harga tambahan untuk sesi malam","icon":"i-lucide-moon","description":["Sesi di malam hari buka dari jam 18:00 sampai 20:00","Harga dasar 6x + Rp.100.000 untuk sesi malam","Harga dasar 8x + Rp.150.000 untuk sesi malam","Harga dasar 10x + Rp. 200.000 untuk sesi malam","Harga dasar 12x + Rp. 250.000 untuk sesi malam"]},{"title":"Sesi Akhir Pekan","subtitle":"Harga tambahan untuk sesi akhir pekan","icon":"i-lucide-clock","description":["Sesi di akhir pekan buka dari jam 08:00 sampai 17:00","Harga dasar 6x + Rp.100.000 untuk sesi akhir pekan","Harga dasar 8x + Rp.150.000 untuk sesi akhir pekan","Harga dasar 10x + Rp. 200.000 untuk sesi akhir pekan","Harga dasar 12x + Rp. 250.000 untuk sesi akhir pekan"]}]}},{"id":"material-services","type":"course_material","data":{"headline":"Materi kursus yang akan Anda pelajari","title":"Materi Kursus","description":"Materi kursus yang Anda pelajari akan membuat Anda lebih percaya diri dalam mengemudi.","materials":[{"title":"Pengenalan Kendaran dan Kontrol Dasar","icon":"i-lucide-book-open","description":["1) Latihan kokpit (posisi duduk yang ergonomis, pengaturan spion tengah dan samping, serta penggunaan sabuk pengaman)","2) Pengenalan instrumen (pedal gas, rem, tuas transmisi, rem tangan, lampu indikator dalam dashboard)","3) Pemeriksaan keselamatan (mengecek kondisi ban, oli, dan air radiator sebelum berkendara)"]},{"title":"Pengendalian Awal","icon":"i-lucide-shield-check","description":["1) Menghidupkan & mematikan mesin (prosedur standar menyalakan mesin dengan aman)","2) Teknik pedal gas dengan aman (seimbang dan halus)","3) Teknik pengereman dan berhenti (pengereman yang halus dan cara berhenti di titik tertentu secara presisi)"]},{"title":"Teknik Manuver Dasar","icon":"i-lucide-radar","description":["1) Kontrol kemudi (teknik memutar setir saat berbelok patah)","2) Berjalan mundur (mengendalikan mobil berjalan mundur hanya menggunakan spion)","3) Berbelok di persimpangan (teknik mengambil sudut belokan yang tepat ke kiri maupun ke kanan)"]},{"title":"Teknik Mengemudi di Tanjakan & Turunan","icon":"i-lucide-car","description":["1) Teknik start-stop di tanjakan","2) Teknik start-stop di turunan"]},{"title":"Teknik Parkir","icon":"i-lucide-car","description":["1) Parkir mundur serong atau lurus (masuk ke slot parkir dengan posisi mobil mundur)","2) Parkir paralel (teknik menyisipkan mobil diantara dua kendaraan lain secara sejajar)"]},{"title":"Mengemudi di Jalan Raya","icon":"i-lucide-car","description":["1) Rambu dan marka jalan (mematuhi rambu lalu lintas, tanda di larang parkir dan garis marka jalan)","2) Etika berkendara (penggunaan lampu sein, menjaga jarak aman, dan cara mendahului kendaraan lain dengan benar)","3) Blind spot (teknik memeriksa area yang tidak terlihat oleh spion sebelum berpindah jalur)"]}]}},{"id":"areas-services","type":"service_areas","data":{"headline":"Cakupan","title":"Area Layanan","description":"Rute pelatihan kami mencakup area utama di sekitar Alam Sutera, BSD City, Gading Serpong dan Tangerang Sekitarnya.","footer":"Sesi berangkat dari titik penjemputan. Hubungi kami untuk ketersediaan rute spesifik.","areas":["Alam Sutera","BSD City","Gading Serpong","Tangerang Sekitarnya"]}},{"id":"cta-services","type":"cta","data":{"heading":"Siap untuk Memulai?","description":"Pesan konsultasi gratis atau lihat pilihan paket kami untuk memulai perjalanan mengemudi Anda.","links":[{"label":"Pesan Konsultasi","to":"https://wa.me/628119124848?text=Halo%20Drive%20Master%2C%20saya%20ingin%20bertanya%20tentang%20kursus%20mengemudi","color":"warning","icon":"i-simple-icons-whatsapp"},{"label":"Lihat Paket","to":"/packages","color":"neutral","variant":"outline","icon":"i-lucide-package"}],"buttonText":"Lihat Paket","buttonLink":"/packages","buttonIcon":"i-lucide-package","secondaryButtonText":"Chat WhatsApp","secondaryButtonLink":"https://wa.me/6285286160029?text=Halo%20Drive%20Master%2C%20saya%20ingin%20bertanya%20tentang%20kursus%20mengemudi","secondaryButtonIcon":"i-simple-icons-whatsapp"}}]	2026-07-16 15:59:59.02033+07	2026-08-05 20:25:42.06494+07	\N
00000000-0000-0000-0000-000000000001	Home	/	published	[{"data":{"bgImage":"https://ik.imagekit.io/oy4rsvid5/pages/1d8650b7-abca-4d5b-b193-a1b3db1bdcf2_1ueSTrIv_.png?tr=w-700,h-450,fo-auto","ctaLink":"/auth/register","ctaText":"Pesan Sesi Pertama","features":[{"icon":"i-lucide-battery-charging","title":"Safety-Confidence-Mastery"},{"icon":"i-lucide-shield-check","title":"Kontrol Ganda"}],"heading":"Kursus Mengemudi Di Tangerang Dengan Belajar Aman, Mengemudi Percaya Diri!","secondaryCtaLink":"/packages","secondaryCtaText":"Lihat Paket","subheading":"DriveMaster Driving Course adalah lembaga kursus mengemudi yang memposisikan diri sebagai Driving Education System berkomitmen membantu setiap siswa menjadi pengemudi yang aman (Safety), percaya diri (Confidence), dan terampil (Mastery).\\nDriveMaster dibangun berdasarkan tiga prinsip utama: Keselamatan Sebagai Prioritas Utama, Pembelajaran Berbasis Kompetensi ,dan Pengembangan Kepercayaan Diri.\\n"},"id":"hero-home","type":"hero"},{"data":{"description":"Materi kursus yang Anda pelajari akan membuat Anda lebih percaya diri dalam mengemudi.","headline":"Materi kursus yang akan Anda pelajari","materials":[{"description":["1) Latihan kokpit (posisi duduk yang ergonomis,pengaturan spion tengah dan samping, sertapenggunaan sabuk pengaman)","2) Pengenalan instrumen (pedal gas, rem, tuastransmisi, rem tangan, lampu indikator dalamdashboard)","3) Pemeriksaan keselamatan (mengecek kondisi ban,oli, dan air radiator sebelum berkendara)"],"icon":"i-lucide-book-open","title":"Pengenalan Kendaraan dan Kontrol Dasar"},{"description":["1) Menghidupkan & mematikan mesin (prosedur standar menyalakan mesin dengan aman)","2) Teknik pedal gas dengan aman (seimbang dan halus)","3) Teknik pengereman dan berhenti (pengereman yang halus dan cara berhenti di titik tertentu secara presisi)"],"icon":"i-lucide-shield-check","title":"Pengendalian Awal"},{"description":["1) Kontrol kemudi (teknik memutar setir saat berbelok patah)","2) Berjalan mundur (mengendalikan mobil berjalanmundur hanya menggunakan spion)","3) Berbelok di persimpangan (teknik mengambil sudut belokan yang tepat ke kiri maupun ke kanan)"],"icon":"i-lucide-radar","title":"Teknik Manuver Dasar"},{"description":["1) Teknik start-stop di tanjakan","2) Teknik start-stop di turunan"],"icon":"i-lucide-car","title":"Teknik Mengemudi di Tanjakan & Turunan"},{"description":["1) Parkir mundur serong atau lurus (masuk ke slot parkir dengan posisi mobil mundur)","2) Parkir paralel (teknik menyisipkan mobil diantara dua kendaraan lain secara sejajar)"],"icon":"i-lucide-car","title":"Teknik Parkir"},{"description":["1) Rambu dan marka jalan (mematuhi rambu lalu lintas,tanda di larang parkir dan garis marka jalan)","2)Etika berkendara (penggunaan lampu sein, menjaga jarak aman, dan cara mendahului kendaraan lain dengan benar)","3) Blind spot (teknik memeriksa area yang tidak terlihat oleh spion sebelum berpindah jalur)"],"icon":"i-lucide-car","title":"Mengemudi di Jalan Raya"}],"title":"Materi Kursus"},"id":"material-home","type":"course_material"},{"data":{"description":"Pesan sesi pertama Anda hari ini dan rasakan kesenangan belajar di kendaraan listrik premium ya!","heading":"Siap Mengemudi di Masa Depan?","links":[{"color":"warning","icon":"i-lucide-rocket","label":"Mulai Perjalanan Anda","to":"/auth/register"},{"color":"neutral","label":"Lihat Semua Paket","to":"/packages","trailingIcon":"i-lucide-arrow-right","variant":"ghost"}]},"id":"cta-home","type":"cta"}]	2026-07-10 09:49:21.414192+07	2026-08-05 15:47:10.115879+07	\N
\.


--
-- Data for Name: provinces; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.provinces (id, name) FROM stdin;
11	Aceh (NAD)
12	Sumatera Utara
13	Sumatera Barat
14	Riau
15	Jambi
16	Sumatera Selatan
17	Bengkulu
18	Lampung
19	Kepulauan Bangka Belitung
21	Kepulauan Riau
31	Dki Jakarta
32	Jawa Barat
33	Jawa Tengah
34	Daerah Istimewa Yogyakarta
35	Jawa Timur
36	Banten
51	Bali
52	Nusa Tenggara Barat
53	Nusa Tenggara Timur
61	Kalimantan Barat
62	Kalimantan Tengah
63	Kalimantan Selatan
64	Kalimantan Timur
65	Kalimantan Utara
71	Sulawesi Utara
72	Sulawesi Tengah
73	Sulawesi Selatan
74	Sulawesi Tenggara
75	Gorontalo
76	Sulawesi Barat
81	Maluku
82	Maluku Utara
91	Papua
92	Papua Barat
93	Papua Selatan
94	Papua Tengah
95	Papua Pegunungan
96	Papua Barat Daya
\.


--
-- Data for Name: regencies; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.regencies (id, province_id, name, type) FROM stdin;
1101	11	Simeulue	Kabupaten
1102	11	Aceh Singkil	Kabupaten
1103	11	Aceh Selatan	Kabupaten
1104	11	Aceh Tenggara	Kabupaten
1105	11	Aceh Timur	Kabupaten
1106	11	Aceh Tengah	Kabupaten
1107	11	Aceh Barat	Kabupaten
1108	11	Aceh Besar	Kabupaten
1109	11	Pidie	Kabupaten
1110	11	Bireuen	Kabupaten
1111	11	Aceh Utara	Kabupaten
1112	11	Aceh Barat Daya	Kabupaten
1113	11	Gayo Lues	Kabupaten
1114	11	Aceh Tamiang	Kabupaten
1115	11	Nagan Raya	Kabupaten
1116	11	Aceh Jaya	Kabupaten
1117	11	Bener Meriah	Kabupaten
1118	11	Pidie Jaya	Kabupaten
1171	11	Banda Aceh	Kota
1172	11	Sabang	Kota
1173	11	Lhokseumawe	Kota
1174	11	Langsa	Kota
1175	11	Subulussalam	Kota
1201	12	Nias	Kabupaten
1202	12	Mandailing Natal	Kabupaten
1203	12	Tapanuli Selatan	Kabupaten
1204	12	Tapanuli Tengah	Kabupaten
1205	12	Tapanuli Utara	Kabupaten
1206	12	Toba Samosir	Kabupaten
1207	12	Labuhan Batu	Kabupaten
1208	12	Asahan	Kabupaten
1209	12	Simalungun	Kabupaten
1210	12	Dairi	Kabupaten
1211	12	Karo	Kabupaten
1212	12	Deliserdang	Kabupaten
1213	12	Langkat	Kabupaten
1214	12	Nias Selatan	Kabupaten
1215	12	Humbang Hasundutan	Kabupaten
1216	12	Pakpak Bharat	Kabupaten
1217	12	Samosir	Kabupaten
1218	12	Serdang Bedagai	Kabupaten
1219	12	Batu Bara	Kabupaten
1220	12	Padang Lawas Utara	Kabupaten
1221	12	Padang Lawas	Kabupaten
1222	12	Labuhan Batu Selatan	Kabupaten
1223	12	Labuhan Batu Utara	Kabupaten
1224	12	Nias Utara	Kabupaten
1225	12	Nias Barat	Kabupaten
1271	12	Sibolga	Kota
1272	12	Tanjung Balai	Kota
1273	12	Pematang Siantar	Kota
1274	12	Tebing Tinggi	Kota
1275	12	Medan	Kota
1276	12	Binjai	Kota
1277	12	Padangsidimpuan	Kota
1278	12	Gunungsitoli	Kota
1301	13	Kepulauan Mentawai	Kabupaten
1302	13	Pesisir Selatan	Kabupaten
1303	13	Solok	Kabupaten
1304	13	Sijunjung	Kabupaten
1305	13	Tanah Datar	Kabupaten
1306	13	Padang Pariaman	Kabupaten
1307	13	Agam	Kabupaten
1308	13	Lima Puluh Kota	Kabupaten
1309	13	Pasaman	Kabupaten
1310	13	Solok Selatan	Kabupaten
1311	13	Dharmasraya	Kabupaten
1312	13	Pasaman Barat	Kabupaten
1371	13	Padang	Kota
1372	13	Solok	Kota
1373	13	Sawah Lunto	Kota
1374	13	Padang Panjang	Kota
1375	13	Bukittinggi	Kota
1376	13	Payakumbuh	Kota
1377	13	Pariaman	Kota
1401	14	Kuantan Singingi	Kabupaten
1402	14	Indragiri Hulu	Kabupaten
1403	14	Indragiri Hilir	Kabupaten
1404	14	Pelalawan	Kabupaten
1405	14	S I A K	Kabupaten
1406	14	Kampar	Kabupaten
1407	14	Rokan Hulu	Kabupaten
1408	14	Bengkalis	Kabupaten
1409	14	Rokan Hilir	Kabupaten
1410	14	Kepulauan Meranti	Kabupaten
1471	14	Pekanbaru	Kota  
1473	14	Dumai	Kota
1501	15	Kerinci	Kabupaten
1502	15	Merangin	Kabupaten
1503	15	Sarolangun	Kabupaten
1504	15	Batang Hari	Kabupaten
1505	15	Muaro Jambi	Kabupaten
1506	15	Tanjung Jabung Timur	Kabupaten
1507	15	Tanjung Jabung Barat	Kabupaten
1508	15	Tebo	Kabupaten
1509	15	Bungo	Kabupaten
1571	15	Jambi	Kota
1572	15	Sungai Penuh	Kota
1601	16	Ogan Komering Ulu	Kabupaten
1602	16	Ogan Komering Ilir	Kabupaten
1603	16	Muara Enim	Kabupaten
1604	16	Lahat	Kabupaten
1605	16	Musi Rawas	Kabupaten
1606	16	Musi Banyuasin	Kabupaten
1607	16	Banyuasin	Kabupaten
1608	16	Ogan Komering Ulu Selatan	Kabupaten
1609	16	Ogan Komering Ulu Timur	Kabupaten
1610	16	Ogan Ilir	Kabupaten
1611	16	Empat Lawang	Kabupaten
1612	16	Penukal Abab Lematang Ilir	Kabupaten
1613	16	Musi Rawas Utara	Kabupaten
1671	16	Palembang	Kota
1672	16	Prabumulih	Kota
1673	16	Pagar Alam	Kota
1674	16	Lubuklinggau	Kota
1701	17	Bengkulu Selatan	Kabupaten
1702	17	Rejang Lebong	Kabupaten
1703	17	Bengkulu Utara	Kabupaten
1704	17	Kaur	Kabupaten
1705	17	Seluma	Kabupaten
1706	17	Mukomuko	Kabupaten
1707	17	Lebong	Kabupaten
1708	17	Kepahiang	Kabupaten
1709	17	Bengkulu Tengah	Kabupaten
1771	17	Bengkulu	Kota
1801	18	Lampung Barat	Kabupaten
1802	18	Tanggamus	Kabupaten
1803	18	Lampung Selatan	Kabupaten
1804	18	Lampung Timur	Kabupaten
1805	18	Lampung Tengah	Kabupaten
1806	18	Lampung Utara	Kabupaten
1807	18	Way Kanan	Kabupaten
1808	18	Tulangbawang	Kabupaten
1809	18	Pesawaran	Kabupaten
1810	18	Pringsewu	Kabupaten
1811	18	Mesuji	Kabupaten
1812	18	Tulang Bawang Barat	Kabupaten
1813	18	Pesisir Barat	Kabupaten
1871	18	Bandar Lampung	Kota
1872	18	Metro	Kota
1901	19	Bangka	Kabupaten
1902	19	Belitung	Kabupaten
1903	19	Bangka Barat	Kabupaten
1904	19	Bangka Tengah	Kabupaten
1905	19	Bangka Selatan	Kabupaten
1906	19	Belitung Timur	Kabupaten
1971	19	Pangkal Pinang	Kota
2101	21	Karimun	Kabupaten
2102	21	Bintan	Kabupaten
2103	21	Natuna	Kabupaten
2104	21	Lingga	Kabupaten
2105	21	Kepulauan Anambas	Kabupaten
2171	21	Batam	Kota
2172	21	Tanjung Pinang	Kota
3101	31	Kepulauan Seribu	Kabupaten
3171	31	Jakarta Selatan	Kota
3172	31	Jakarta Timur	Kota
3173	31	Jakarta Pusat	Kota
3174	31	Jakarta Barat	Kota
3175	31	Jakarta Utara	Kota
3201	32	Bogor	Kabupaten
3202	32	Sukabumi	Kabupaten
3203	32	Cianjur	Kabupaten
3204	32	Bandung	Kabupaten
3205	32	Garut	Kabupaten
3206	32	Tasikmalaya	Kabupaten
3207	32	Ciamis	Kabupaten
3208	32	Kuningan	Kabupaten
3209	32	Cirebon	Kabupaten
3210	32	Majalengka	Kabupaten
3211	32	Sumedang	Kabupaten
3212	32	Indramayu	Kabupaten
3213	32	Subang	Kabupaten
3214	32	Purwakarta	Kabupaten  
3215	32	Karawang	Kabupaten
3216	32	Bekasi	Kabupaten
3217	32	Bandung Barat	Kabupaten
3218	32	Pangandaran	Kabupaten
3271	32	Bogor	Kota
3272	32	Sukabumi	Kota
3273	32	Bandung	Kota
3274	32	Cirebon	Kota
3275	32	Bekasi	Kota
3276	32	Depok	Kota
3277	32	Cimahi	Kota
3278	32	Tasikmalaya	Kota
3279	32	Banjar	Kota
3301	33	Cilacap	Kabupaten
3302	33	Banyumas	Kabupaten
3303	33	Purbalingga	Kabupaten
3304	33	Banjarnegara	Kabupaten
3305	33	Kebumen	Kabupaten
3306	33	Purworejo	Kabupaten
3307	33	Wonosobo	Kabupaten
3308	33	Magelang	Kabupaten
3309	33	Boyolali	Kabupaten
3310	33	Klaten	Kabupaten
3311	33	Sukoharjo	Kabupaten
3312	33	Wonogiri	Kabupaten
3313	33	Karanganyar	Kabupaten
3314	33	Sragen	Kabupaten
3315	33	Grobogan	Kabupaten
3316	33	Blora	Kabupaten
3317	33	Rembang	Kabupaten
3318	33	Pati	Kabupaten
3319	33	Kudus	Kabupaten
3320	33	Jepara	Kabupaten
3321	33	Demak	Kabupaten
3322	33	Semarang	Kabupaten
3323	33	Temanggung	Kabupaten
3324	33	Kendal	Kabupaten
3325	33	Batang	Kabupaten
3326	33	Pekalongan	Kabupaten
3327	33	Pemalang	Kabupaten
3328	33	Tegal	Kabupaten
3329	33	Brebes	Kabupaten
3371	33	Magelang	Kota
3372	33	Surakarta	Kota
3373	33	Salatiga	Kota
3374	33	Semarang	Kota
3375	33	Pekalongan	Kota
3376	33	Tegal	Kota  
3401	34	Kulon Progo	Kabupaten
3402	34	Bantul	Kabupaten
3403	34	Gunung Kidul	Kabupaten
3404	34	Sleman	Kabupaten
3471	34	Yogyakarta	Kota
3501	35	Pacitan	Kabupaten
3502	35	Ponorogo	Kabupaten
3503	35	Trenggalek	Kabupaten
3504	35	Tulungagung	Kabupaten
3505	35	Blitar	Kabupaten
3506	35	Kediri	Kabupaten  
3507	35	Malang	Kabupaten
3508	35	Lumajang	Kabupaten
3509	35	Jember	Kabupaten
3510	35	Banyuwangi	Kabupaten
3511	35	Bondowoso	Kabupaten
3512	35	Situbondo	Kabupaten
3513	35	Probolinggo	Kabupaten
3514	35	Pasuruan	Kabupaten
3515	35	Sidoarjo	Kabupaten
3516	35	Mojokerto	Kabupaten
3517	35	Jombang	Kabupaten
3518	35	Nganjuk	Kabupaten
3519	35	Madiun	Kabupaten
3520	35	Magetan	Kabupaten
3521	35	Ngawi	Kabupaten
3522	35	Bojonegoro	Kabupaten
3523	35	Tuban	Kabupaten
3524	35	Lamongan	Kabupaten
3525	35	Gresik	Kabupaten
3526	35	Bangkalan	Kabupaten
3527	35	Sampang	Kabupaten
3528	35	Pamekasan	Kabupaten
3529	35	Sumenep	Kabupaten
3571	35	Kediri	Kota
3572	35	Blitar	Kota
3573	35	Malang	Kota
3574	35	Probolinggo	Kota
3575	35	Pasuruan	Kota
3576	35	Mojokerto	Kota
3577	35	Madiun	Kota
3578	35	Surabaya	Kota
3579	35	Batu	Kota
3601	36	Pandeglang	Kabupaten
3602	36	Lebak	Kabupaten
3603	36	Tangerang	Kabupaten
3604	36	Serang	Kabupaten
3671	36	Tangerang	Kota
3672	36	Cilegon	Kota
3673	36	Serang	Kota
3674	36	Tangerang Selatan	Kota
5101	51	Jembrana	Kabupaten
5102	51	Tabanan	Kabupaten
5103	51	Badung	Kabupaten
5104	51	Gianyar	Kabupaten
5105	51	Klungkung	Kabupaten
5106	51	Bangli	Kabupaten
5107	51	Karang Asem	Kabupaten
5108	51	Buleleng	Kabupaten
5171	51	Denpasar	Kota
5201	52	Lombok Barat	Kabupaten
5202	52	Lombok Tengah	Kabupaten
5203	52	Lombok Timur	Kabupaten
5204	52	Sumbawa	Kabupaten
5205	52	Dompu	Kabupaten
5206	52	Bima	Kabupaten
5207	52	Sumbawa Barat	Kabupaten
5208	52	Lombok Utara	Kabupaten
5271	52	Mataram	Kota
5272	52	Bima	Kota
5301	53	Sumba Barat	Kabupaten
5302	53	Sumba Timur	Kabupaten
5303	53	Kupang	Kabupaten
5304	53	Timor Tengah Selatan	Kabupaten
5305	53	Timor Tengah Utara	Kabupaten
5306	53	Belu	Kabupaten
5307	53	Alor	Kabupaten
5308	53	Flores Timur	Kabupaten
5309	53	Sikka	Kabupaten
5310	53	Ende	Kabupaten
5311	53	Ngada	Kabupaten
5312	53	Manggarai	Kabupaten
5313	53	Rote Ndao	Kabupaten
5314	53	Manggarai Barat	Kabupaten
5315	53	Sumba Tengah	Kabupaten
5316	53	Sumba Barat Daya	Kabupaten
5317	53	Nagekeo	Kabupaten
5318	53	Manggarai Timur	Kabupaten
5319	53	Sabu Raijua	Kabupaten
5320	53	Malaka	Kabupaten
5371	53	Kupang	Kota
6101	61	Sambas	Kabupaten
6102	61	Mempawah	Kabupaten
6103	61	Sanggau	Kabupaten
6104	61	Ketapang	Kabupaten
6105	61	Sintang	Kabupaten
6106	61	Kapuas Hulu	Kabupaten
6107	61	Bengkayang	Kabupaten
6108	61	Landak	Kabupaten
6109	61	Sekadau	Kabupaten
6110	61	Melawi	Kabupaten
6111	61	Kayong Utara	Kabupaten
6112	61	Kubu Raya	Kabupaten
6171	61	Pontianak	Kota
6172	61	Singkawang	Kota
6201	62	Kotawaringin Barat	Kabupaten
6202	62	Kotawaringin Timur	Kabupaten
6203	62	Kapuas	Kabupaten
6204	62	Barito Selatan	Kabupaten
6205	62	Barito Utara	Kabupaten
6206	62	Sukamara	Kabupaten
6207	62	Lamandau	Kabupaten
6208	62	Seruyan	Kabupaten
6209	62	Katingan	Kabupaten
6210	62	Pulang Pisau	Kabupaten
6211	62	Gunung Mas	Kabupaten
6212	62	Barito Timur	Kabupaten
6213	62	Murung Raya	Kabupaten
6271	62	Palangka Raya	Kota
6301	63	Tanah Laut	Kabupaten
6302	63	Kota Baru	Kabupaten
6303	63	Banjar	Kabupaten
6304	63	Barito Kuala	Kabupaten
6305	63	Tapin	Kabupaten
6306	63	Hulu Sungai Selatan	Kabupaten
6307	63	Hulu Sungai Tengah	Kabupaten
6308	63	Hulu Sungai Utara	Kabupaten
6309	63	Tabalong	Kabupaten
6310	63	Tanah Bumbu	Kabupaten
6311	63	Balangan	Kabupaten
6371	63	Banjarmasin	Kota
6372	63	Banjarbaru	Kota
6401	64	Paser	Kabupaten
6402	64	Kutai Barat	Kabupaten
6403	64	Kutai Kartanegara	Kabupaten
6404	64	Kutai Timur	Kabupaten
6405	64	Berau	Kabupaten
6409	64	Penajam Paser Utara	Kabupaten
6471	64	Balikpapan	Kota
6472	64	Samarinda	Kota
6474	64	Bontang	Kota
6501	65	Malinau	Kabupaten
6502	65	Bulungan	Kabupaten
6503	65	Tana Tidung	Kabupaten
6504	65	Nunukan	Kabupaten
6571	65	Tarakan	Kota
7101	71	Bolaang Mongondow	Kabupaten
7102	71	Minahasa	Kabupaten
7103	71	Kepulauan Sangihe	Kabupaten
7104	71	Kepulauan Talaud	Kabupaten
7105	71	Minahasa Selatan	Kabupaten
7106	71	Minahasa Utara	Kabupaten
7107	71	Bolaang Mongondow Utara	Kabupaten
7108	71	Siau Tagulandang Biaro (Sitaro)	Kabupaten
7109	71	Minahasa Tenggara	Kabupaten
7110	71	Bolaang Mongondow Selatan	Kabupaten
7111	71	Bolaang Mongondow Timur	Kabupaten
7171	71	Manado	Kota
7172	71	Bitung	Kota
7173	71	Tomohon	Kota
7174	71	Kotamobagu	Kota
7201	72	Banggai	Kabupaten
7202	72	Poso	Kabupaten
7203	72	Donggala	Kabupaten
7204	72	Toli-Toli	Kabupaten
7205	72	Buol	Kabupaten
7206	72	Morowali	Kabupaten
7207	72	Banggai Kepulauan	Kabupaten
7208	72	Parigi Moutong	Kabupaten
7209	72	Tojo Una-Una	Kabupaten
7210	72	Sigi	Kabupaten
7271	72	Palu	Kota
7301	73	Kepulauan Selayar	Kabupaten
7302	73	Bulukumba	Kabupaten
7303	73	Bantaeng	Kabupaten
7304	73	Jeneponto	Kabupaten
7305	73	Takalar	Kabupaten
7306	73	Gowa	Kabupaten
7307	73	Sinjai	Kabupaten
7308	73	Bone	Kabupaten
7309	73	Maros	Kabupaten
7310	73	Pangkajene Dan Kepulauan	Kabupaten
7311	73	Barru	Kabupaten
7312	73	Soppeng	Kabupaten
7313	73	Wajo	Kabupaten
7314	73	Sidenreng Rappang	Kabupaten
7315	73	Pinrang	Kabupaten
7316	73	Enrekang	Kabupaten
7317	73	Luwu	Kabupaten
7318	73	Tana Toraja	Kabupaten
7322	73	Luwu Utara	Kabupaten
7325	73	Luwu Timur	Kabupaten
7326	73	Toraja Utara	Kabupaten
7371	73	Makassar	Kota
7372	73	Parepare	Kota
7373	73	Palu	Kota
7401	74	Buton	Kabupaten
7402	74	Muna	Kabupaten
7403	74	Konawe	Kabupaten
7404	74	Kolaka	Kabupaten
7405	74	Konawe Selatan	Kabupaten
7406	74	Bombana	Kabupaten
7407	74	Wakatobi	Kabupaten
7408	74	Kolaka Utara	Kabupaten
7409	74	Buton Utara	Kabupaten
7410	74	Konawe Utara	Kabupaten
7411	74	Kolaka Timur	Kabupaten
7412	74	Konawe Kepulauan	Kabupaten
7413	74	Muna Barat	Kabupaten
7414	74	Buton Tengah	Kabupaten
7415	74	Buton Selatan	Kabupaten
7471	74	Kendari	Kota
7472	74	Baubau	Kota
7501	75	Boalemo	Kabupaten
7502	75	Gorontalo	Kabupaten
7503	75	Pohuwato	Kabupaten
7504	75	Bone Bolango	Kabupaten
7505	75	Gorontalo Utara	Kabupaten
7571	75	Gorontalo	Kota
7601	76	Majene	Kabupaten
7602	76	Polewali Mandar	Kabupaten
7603	76	Mamasa	Kabupaten
7604	76	Mamuju	Kabupaten
7605	76	Mamuju Utara	Kabupaten
7606	76	Mamuju Tengah	Kabupaten
8101	81	Maluku Tenggara Barat	Kabupaten
8102	81	Maluku Tenggara	Kabupaten
8103	81	Maluku Tengah	Kabupaten
8104	81	Buru	Kabupaten
8105	81	Kepulauan Aru	Kabupaten
8106	81	Seram Bagian Barat	Kabupaten
8107	81	Seram Bagian Timur	Kabupaten
8108	81	Maluku Barat Daya	Kabupaten
8109	81	Buru Selatan	Kabupaten
8171	81	Ambon	Kota
8172	81	Tual	Kota
8201	82	Halmahera Barat	Kabupaten
8202	82	Halmahera Tengah	Kabupaten
8203	82	Kepulauan Sula	Kabupaten
8204	82	Halmahera Selatan	Kabupaten
8205	82	Halmahera Utara	Kabupaten
8206	82	Halmahera Timur	Kabupaten
8207	82	Pulau Morotai	Kabupaten
8208	82	Pulau Taliabu	Kabupaten
8271	82	Ternate	Kota
8272	82	Tidore Kepulauan	Kota
9116	91	Sarmi	Kabupaten
9117	91	Keerom	Kabupaten
9118	91	Waropen	Kabupaten
9119	91	Supiori	Kabupaten
9120	91	Mamberamo Raya	Kabupaten
9171	91	Jayapura	Kota
9201	92	Fakfak	Kabupaten
9202	92	Kaimana	Kabupaten
9203	92	Teluk Wondama	Kabupaten
9204	92	Teluk Bintuni	Kabupaten
9205	92	Manokwari	Kabupaten
9206	92	Sorong Selatan	Kabupaten
9207	92	Sorong	Kabupaten
9208	92	Raja Ampat	Kabupaten
9209	92	Tambrauw	Kabupaten
9210	92	Maybrat	Kabupaten
9211	92	Manokwari Selatan	Kabupaten
9212	92	Pegunungan Arfak	Kabupaten
9271	92	Sorong	Kota
9301	93	Asmat	Kabupaten
9302	93	Boven Digoel	Kabupaten
9303	93	Mappi	Kabupaten
9304	93	Merauke	Kabupaten
9401	94	Deiyai (Deliyai)	Kabupaten
9402	94	Dogiyai	Kabupaten
9403	94	Intan Jaya	Kabupaten
9404	94	Mimika	Kabupaten
9408	94	Nabire	Kabupaten
9409	94	Paniai	Kabupaten
9410	94	Puncak	Kabupaten
9411	94	Puncak Jaya	Kabupaten
9501	95	Jayawijaya	Kabupaten
9502	95	Lanny Jaya	Kabupaten
9503	95	Mamberamo Tengah	Kabupaten
9504	95	Nduga	Kabupaten
9505	95	Pegunungan Bintang	Kabupaten
9506	95	Tolikara	Kabupaten
9507	95	Yahukimo	Kabupaten
9508	95	Yalimo	Kabupaten
9601	96	Maybrat	Kabupaten
9602	96	Raja Ampat	Kabupaten
9603	96	Sorong	Kabupaten
9604	96	Sorong Selatan	Kabupaten
9605	96	Tambrauw	Kabupaten
9606	96	Sorong	Kota
\.


--
-- Data for Name: related_articles; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.related_articles (id, article_id, related_article_id, relationship_type) FROM stdin;
\.


--
-- Data for Name: sale_items; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.sale_items (id, sale_id, package_id, package_name, quantity, unit_price, discount, subtotal, created_at) FROM stdin;
\.


--
-- Data for Name: sales; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.sales (id, order_number, payment_id, user_id, package_id, package_name, package_type, total_amount, discount_amount, final_amount, status, source, payment_method, currency, paid_at, refunded_at, notes, created_at, updated_at) FROM stdin;
\.


--
-- Data for Name: tags; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.tags (id, name, slug, description, usage_count, created_at) FROM stdin;
b0000001-0000-0000-0000-000000000001	Tips	tips	Helpful driving tips	5	2026-06-05 09:08:59.899996+07
b0000001-0000-0000-0000-000000000002	Safety	safety	Road safety content	4	2026-06-05 09:08:59.904989+07
b0000001-0000-0000-0000-000000000003	Beginners	beginners	Content for new drivers	3	2026-06-05 09:08:59.910211+07
\.


--
-- Name: districts_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.districts_id_seq', 7148, true);


--
-- Name: provinces_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.provinces_id_seq', 1, false);


--
-- Name: regencies_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.regencies_id_seq', 1, false);


--
-- Name: add_ons add_ons_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.add_ons
    ADD CONSTRAINT add_ons_pkey PRIMARY KEY (id);


--
-- Name: article_tags article_tags_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.article_tags
    ADD CONSTRAINT article_tags_pkey PRIMARY KEY (id);


--
-- Name: articles articles_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.articles
    ADD CONSTRAINT articles_pkey PRIMARY KEY (id);


--
-- Name: cars cars_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.cars
    ADD CONSTRAINT cars_pkey PRIMARY KEY (id);


--
-- Name: categories categories_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.categories
    ADD CONSTRAINT categories_pkey PRIMARY KEY (id);


--
-- Name: contact_inquiries contact_inquiries_pkey; Type: CONSTRAINT; Schema: public; Owner: admin_drive
--

ALTER TABLE ONLY public.contact_inquiries
    ADD CONSTRAINT contact_inquiries_pkey PRIMARY KEY (id);


--
-- Name: districts districts_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.districts
    ADD CONSTRAINT districts_pkey PRIMARY KEY (id);


--
-- Name: faqs faqs_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.faqs
    ADD CONSTRAINT faqs_pkey PRIMARY KEY (id);


--
-- Name: general_settings general_settings_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.general_settings
    ADD CONSTRAINT general_settings_pkey PRIMARY KEY (id);


--
-- Name: monthly_sales monthly_sales_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.monthly_sales
    ADD CONSTRAINT monthly_sales_pkey PRIMARY KEY (year, month);


--
-- Name: package_benefits package_benefits_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.package_benefits
    ADD CONSTRAINT package_benefits_pkey PRIMARY KEY (id);


--
-- Name: packages packages_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.packages
    ADD CONSTRAINT packages_pkey PRIMARY KEY (id);


--
-- Name: pages pages_pkey; Type: CONSTRAINT; Schema: public; Owner: admin_drive
--

ALTER TABLE ONLY public.pages
    ADD CONSTRAINT pages_pkey PRIMARY KEY (id);


--
-- Name: provinces provinces_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.provinces
    ADD CONSTRAINT provinces_pkey PRIMARY KEY (id);


--
-- Name: regencies regencies_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.regencies
    ADD CONSTRAINT regencies_pkey PRIMARY KEY (id);


--
-- Name: related_articles related_articles_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.related_articles
    ADD CONSTRAINT related_articles_pkey PRIMARY KEY (id);


--
-- Name: sale_items sale_items_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sale_items
    ADD CONSTRAINT sale_items_pkey PRIMARY KEY (id);


--
-- Name: sales sales_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sales
    ADD CONSTRAINT sales_pkey PRIMARY KEY (id);


--
-- Name: tags tags_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.tags
    ADD CONSTRAINT tags_pkey PRIMARY KEY (id);


--
-- Name: idx_article_tag; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX idx_article_tag ON public.article_tags USING btree (article_id, tag_id);


--
-- Name: idx_articles_deleted_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_articles_deleted_at ON public.articles USING btree (deleted_at);


--
-- Name: idx_articles_slug; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX idx_articles_slug ON public.articles USING btree (slug) WHERE (deleted_at IS NULL);


--
-- Name: idx_cars_license_plate; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX idx_cars_license_plate ON public.cars USING btree (license_plate);


--
-- Name: idx_categories_slug; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX idx_categories_slug ON public.categories USING btree (slug);


--
-- Name: idx_faqs_deleted_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_faqs_deleted_at ON public.faqs USING btree (deleted_at);


--
-- Name: idx_pages_deleted_at; Type: INDEX; Schema: public; Owner: admin_drive
--

CREATE INDEX idx_pages_deleted_at ON public.pages USING btree (deleted_at);


--
-- Name: idx_pages_slug; Type: INDEX; Schema: public; Owner: admin_drive
--

CREATE UNIQUE INDEX idx_pages_slug ON public.pages USING btree (slug) WHERE (deleted_at IS NULL);


--
-- Name: idx_related; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX idx_related ON public.related_articles USING btree (article_id, related_article_id);


--
-- Name: idx_sale_items_sale_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_sale_items_sale_id ON public.sale_items USING btree (sale_id);


--
-- Name: idx_sales_order_number; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX idx_sales_order_number ON public.sales USING btree (order_number);


--
-- Name: idx_sales_package_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_sales_package_id ON public.sales USING btree (package_id);


--
-- Name: idx_sales_payment_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_sales_payment_id ON public.sales USING btree (payment_id);


--
-- Name: idx_sales_status; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_sales_status ON public.sales USING btree (status);


--
-- Name: idx_sales_user_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_sales_user_id ON public.sales USING btree (user_id);


--
-- Name: idx_tags_name; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX idx_tags_name ON public.tags USING btree (name);


--
-- Name: idx_tags_slug; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX idx_tags_slug ON public.tags USING btree (slug);


--
-- Name: article_tags fk_article_tags_tag; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.article_tags
    ADD CONSTRAINT fk_article_tags_tag FOREIGN KEY (tag_id) REFERENCES public.tags(id);


--
-- Name: articles fk_articles_category; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.articles
    ADD CONSTRAINT fk_articles_category FOREIGN KEY (category_id) REFERENCES public.categories(id);


--
-- Name: categories fk_categories_parent; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.categories
    ADD CONSTRAINT fk_categories_parent FOREIGN KEY (parent_id) REFERENCES public.categories(id);


--
-- Name: sale_items fk_sales_items; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sale_items
    ADD CONSTRAINT fk_sales_items FOREIGN KEY (sale_id) REFERENCES public.sales(id);


--
-- PostgreSQL database dump complete
--

\unrestrict EHYbGxzaUf53OF9JeJdInm5Prm97WPQG1fZ6qFH2A4pje7ok5XrVqqxEyyFzHc5

