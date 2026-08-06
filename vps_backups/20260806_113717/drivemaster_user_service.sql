--
-- PostgreSQL database dump
--

\restrict Ubfhn5snx0teniGilnnys7MoxEz2It4dMom8nOPupXxXgevPtKOIussYMqNNDCS

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

ALTER TABLE ONLY public.testimonial_media DROP CONSTRAINT fk_testimonial_media_testimonial;
DROP INDEX public.idx_users_username;
DROP INDEX public.idx_users_email_address;
DROP INDEX public.idx_testimonials_user_id;
DROP INDEX public.idx_testimonial_media_testimonial_id;
DROP INDEX public.idx_roles_name;
DROP INDEX public.idx_member_profiles_user_id;
DROP INDEX public.idx_instructor_recurring_schedules_instructor_id;
DROP INDEX public.idx_instructor_profiles_user_id;
DROP INDEX public.idx_instructor_areas_instructor_id;
DROP INDEX public.idx_entitlements_member_id;
DROP INDEX public.idx_entitlements_booking_id;
DROP INDEX public.idx_certifications_member_id;
DROP INDEX public.idx_certifications_instructor_id;
DROP INDEX public.idx_certifications_entitlement_id;
ALTER TABLE ONLY public.work_experiences DROP CONSTRAINT work_experiences_pkey;
ALTER TABLE ONLY public.users DROP CONSTRAINT users_pkey;
ALTER TABLE ONLY public.testimonials DROP CONSTRAINT testimonials_pkey;
ALTER TABLE ONLY public.testimonial_media DROP CONSTRAINT testimonial_media_pkey;
ALTER TABLE ONLY public.roles DROP CONSTRAINT roles_pkey;
ALTER TABLE ONLY public.member_profiles DROP CONSTRAINT member_profiles_pkey;
ALTER TABLE ONLY public.instructor_recurring_schedules DROP CONSTRAINT instructor_recurring_schedules_pkey;
ALTER TABLE ONLY public.instructor_profiles DROP CONSTRAINT instructor_profiles_pkey;
ALTER TABLE ONLY public.entitlements DROP CONSTRAINT entitlements_pkey;
ALTER TABLE ONLY public.certifications DROP CONSTRAINT certifications_pkey;
ALTER TABLE public.work_experiences ALTER COLUMN id DROP DEFAULT;
ALTER TABLE public.testimonial_media ALTER COLUMN id DROP DEFAULT;
ALTER TABLE public.roles ALTER COLUMN id DROP DEFAULT;
DROP SEQUENCE public.work_experiences_id_seq;
DROP TABLE public.work_experiences;
DROP TABLE public.users;
DROP TABLE public.testimonials;
DROP SEQUENCE public.testimonial_media_id_seq;
DROP TABLE public.testimonial_media;
DROP SEQUENCE public.roles_id_seq;
DROP TABLE public.roles;
DROP TABLE public.member_profiles;
DROP TABLE public.instructor_recurring_schedules;
DROP TABLE public.instructor_profiles;
DROP TABLE public.instructor_areas;
DROP TABLE public.entitlements;
DROP TABLE public.certifications;
SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: certifications; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.certifications (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    instructor_id uuid,
    cert_type character varying(50) NOT NULL,
    cert_number character varying(100) NOT NULL,
    issued_by character varying(255),
    issued_date timestamp with time zone NOT NULL,
    expiry_date timestamp with time zone,
    status character varying(20) DEFAULT 'pending'::character varying,
    document_url character varying(500),
    notes text,
    verified_at timestamp with time zone,
    verified_by text,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    member_id uuid NOT NULL,
    entitlement_id uuid
);


ALTER TABLE public.certifications OWNER TO postgres;

--
-- Name: entitlements; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.entitlements (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    member_id uuid NOT NULL,
    booking_id uuid,
    package_id uuid,
    package_name character varying(255),
    is_night_session boolean DEFAULT false,
    is_weekend_session boolean DEFAULT false,
    total_sessions bigint DEFAULT 0,
    remaining bigint DEFAULT 0,
    used_sessions bigint DEFAULT 0,
    start_date timestamp with time zone,
    end_date timestamp with time zone,
    status character varying(20) DEFAULT 'active'::character varying,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.entitlements OWNER TO postgres;

--
-- Name: instructor_areas; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.instructor_areas (
    instructor_id uuid NOT NULL,
    area_type character varying(20) NOT NULL,
    area_id bigint NOT NULL,
    id uuid DEFAULT gen_random_uuid(),
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.instructor_areas OWNER TO postgres;

--
-- Name: instructor_profiles; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.instructor_profiles (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    license_number character varying(50),
    bnsp_certificate_number character varying(50),
    years_of_experience bigint DEFAULT 0,
    bio text,
    license_expiry timestamp with time zone,
    photo_url character varying(500),
    is_active boolean DEFAULT true,
    number_of_students bigint DEFAULT 0,
    sessions_completed bigint DEFAULT 0,
    average_rating numeric DEFAULT 0,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    description character varying(500),
    specialization character varying(50)
);


ALTER TABLE public.instructor_profiles OWNER TO postgres;

--
-- Name: instructor_recurring_schedules; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.instructor_recurring_schedules (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    instructor_id uuid NOT NULL,
    day_of_week bigint NOT NULL,
    start_time character varying(10) NOT NULL,
    end_time character varying(10) NOT NULL,
    is_active boolean DEFAULT true,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.instructor_recurring_schedules OWNER TO postgres;

--
-- Name: member_profiles; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.member_profiles (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    sessions_completed bigint DEFAULT 0,
    training_time bigint DEFAULT 0,
    average_rating numeric DEFAULT 0,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    identity_full_name character varying(255),
    identity_fullname character varying(255) DEFAULT ''::character varying
);


ALTER TABLE public.member_profiles OWNER TO postgres;

--
-- Name: roles; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.roles (
    id bigint NOT NULL,
    name character varying(50) NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.roles OWNER TO postgres;

--
-- Name: roles_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.roles_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.roles_id_seq OWNER TO postgres;

--
-- Name: roles_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.roles_id_seq OWNED BY public.roles.id;


--
-- Name: testimonial_media; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.testimonial_media (
    id bigint NOT NULL,
    url character varying(255) NOT NULL,
    media_type character varying(20),
    sort_order bigint DEFAULT 0,
    created_at timestamp with time zone,
    testimonial_id uuid
);


ALTER TABLE public.testimonial_media OWNER TO postgres;

--
-- Name: testimonial_media_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.testimonial_media_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.testimonial_media_id_seq OWNER TO postgres;

--
-- Name: testimonial_media_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.testimonial_media_id_seq OWNED BY public.testimonial_media.id;


--
-- Name: testimonials; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.testimonials (
    user_id uuid NOT NULL,
    user_name character varying(150) NOT NULL,
    user_image character varying(255),
    user_role character varying(50),
    content text NOT NULL,
    rating numeric(2,1),
    tags character varying(255),
    status character varying(20) DEFAULT 'draft'::character varying,
    is_featured boolean DEFAULT false,
    added_by uuid NOT NULL,
    added_at timestamp with time zone,
    sort_order bigint DEFAULT 0,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    id uuid NOT NULL,
    CONSTRAINT chk_testimonials_rating CHECK (((rating >= (1)::numeric) AND (rating <= (5)::numeric)))
);


ALTER TABLE public.testimonials OWNER TO postgres;

--
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    first_name text,
    last_name text,
    username character varying(120) NOT NULL,
    password_hash character varying(255) NOT NULL,
    email_address character varying(190) NOT NULL,
    phone_number character varying(20),
    image character varying(500),
    date_of_birth date,
    address character varying(255),
    is_active boolean DEFAULT true,
    is_verified boolean DEFAULT false,
    role_id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.users OWNER TO postgres;

--
-- Name: work_experiences; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.work_experiences (
    id bigint NOT NULL,
    instructor_id uuid NOT NULL,
    company_name character varying(255) NOT NULL,
    role character varying(100) NOT NULL,
    start_date timestamp with time zone NOT NULL,
    end_date timestamp with time zone,
    description text,
    is_verified boolean DEFAULT false,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.work_experiences OWNER TO postgres;

--
-- Name: work_experiences_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.work_experiences_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.work_experiences_id_seq OWNER TO postgres;

--
-- Name: work_experiences_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.work_experiences_id_seq OWNED BY public.work_experiences.id;


--
-- Name: roles id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.roles ALTER COLUMN id SET DEFAULT nextval('public.roles_id_seq'::regclass);


--
-- Name: testimonial_media id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.testimonial_media ALTER COLUMN id SET DEFAULT nextval('public.testimonial_media_id_seq'::regclass);


--
-- Name: work_experiences id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.work_experiences ALTER COLUMN id SET DEFAULT nextval('public.work_experiences_id_seq'::regclass);


--
-- Data for Name: certifications; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.certifications (id, instructor_id, cert_type, cert_number, issued_by, issued_date, expiry_date, status, document_url, notes, verified_at, verified_by, created_at, updated_at, member_id, entitlement_id) FROM stdin;
50d51c2b-a94f-40e1-873f-d9cd52a95b30	\N	package_completion	CERT-31252d3d-99999999	10x Session Package	2026-07-30 10:02:54.95295+07	\N	verified		Package: 10x Session Package | Entitlement: 6a62af4e-b950-4aa5-a294-f6ea58856b15	2026-07-30 10:02:54.95295+07	\N	2026-07-30 10:02:54.95295+07	2026-07-30 10:02:54.95295+07	99999999-9999-9999-9999-999999999999	6a62af4e-b950-4aa5-a294-f6ea58856b15
01d007d9-8acd-4b61-a964-7c5f3eb6bc35	\N	package_completion	CERT-4601ad1b-fbb44c95	Paket Percobaan	2026-07-30 19:07:19.621141+07	\N	verified		Package: Paket Percobaan | Entitlement: cbf419f2-9a86-4729-bf4a-4973aeb4afc2	2026-07-30 19:07:19.621141+07	\N	2026-07-30 19:07:19.621141+07	2026-07-30 19:07:19.621141+07	fbb44c95-79aa-421e-ba3c-b4e7b87668d8	cbf419f2-9a86-4729-bf4a-4973aeb4afc2
e575e8a0-4aa1-4e2a-9bb3-476a38bbbd3b	\N	package_completion	CERT-4601ad1b-fead48e2	Paket Percobaan	2026-08-04 10:38:40.231775+07	\N	verified		Package: Paket Percobaan | Entitlement: ff6577f6-2375-4cc0-abd7-a64cec8b79f2	2026-08-04 10:38:40.231775+07	\N	2026-08-04 10:38:40.231775+07	2026-08-04 10:38:40.231775+07	fead48e2-b1b6-4620-9b18-889e64355698	ff6577f6-2375-4cc0-abd7-a64cec8b79f2
1bc31d06-f5be-4e79-97d0-fd56150a71d3	\N	package_completion	CERT-4601ad1b-3c09cb50	Paket Percobaan	2026-08-04 10:55:31.499985+07	\N	verified		Package: Paket Percobaan | Entitlement: dbaedb67-dc5a-4839-8a9c-b15166650baf	2026-08-04 10:55:31.499985+07	\N	2026-08-04 10:55:31.499985+07	2026-08-04 10:55:31.499985+07	3c09cb50-3d99-4ff8-a57f-44fd4135d070	dbaedb67-dc5a-4839-8a9c-b15166650baf
\.


--
-- Data for Name: entitlements; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.entitlements (id, member_id, booking_id, package_id, package_name, is_night_session, is_weekend_session, total_sessions, remaining, used_sessions, start_date, end_date, status, created_at, updated_at) FROM stdin;
17831666-b5b1-45ee-88ae-f3254eba1c13	4d9c9de3-3e61-448a-b475-51f20aa2f41e	3c3c0ead-befa-475f-819a-574a45e39543	11111111-1111-1111-1111-111111111101	6x	f	f	6	6	0	2026-07-22 16:42:51.641367+07	2027-07-22 16:42:51.641367+07	active	2026-07-22 16:42:51.641367+07	2026-07-22 16:42:51.641367+07
a9de4846-b6ab-4a0e-beb4-efec1a27d1cf	3c09cb50-3d99-4ff8-a57f-44fd4135d070	3c3c0ead-befa-475f-819a-574a45e39543	11111111-1111-1111-1111-111111111101	Pembelian 6x	f	f	6	5	1	2026-07-24 11:05:29.372228+07	2027-07-24 11:05:29.372228+07	active	2026-07-24 11:05:29.372228+07	2026-07-24 11:05:29.372228+07
1de4436b-2a2f-4351-9e0b-2fd78966c2d9	a428042e-e617-48df-8688-cb9ffa7f8c32	7fdb8006-33bd-4116-a104-88076e236d0f	11111111-1111-1111-1111-111111111101	6x	f	f	6	5	1	2026-07-30 07:00:00+07	\N	active	2026-07-30 05:19:38.725819+07	2026-07-30 05:19:38.725819+07
6a62af4e-b950-4aa5-a294-f6ea58856b15	99999999-9999-9999-9999-999999999999	3ff22915-363d-4896-a9d9-b5695a6fdecf	31252d3d-3842-456c-b1b8-34464799cec0	10x Session Package	f	f	10	0	10	2026-06-13 09:11:59.160042+07	\N	used	2026-07-13 09:11:59.160174+07	2026-07-30 10:02:54.944757+07
6a3bfefc-ef05-4d5b-a094-0a7744a30786	e3747d57-fa14-4c52-b4b5-bb9de614cfc2	3abdb746-9bd7-417f-800b-d1524e3cc94c	11111111-1111-1111-1111-111111111101	6x	f	f	6	1	5	2026-07-15 10:10:06.894142+07	\N	active	2026-07-15 10:10:06.894142+07	2026-07-30 10:10:06.894142+07
c2d1e3f4-f5a6-7890-a1b2-c3d4e5f67899	a2b1c3d4-e5f6-7890-a1b2-c3d4e5f67899	d2e1f3a4-f5b6-7890-a1b2-c3d4e5f67899	e2f1a3b4-f5c6-7890-a1b2-c3d4e5f67899	6x Package	f	f	6	1	5	2026-07-15 10:23:56.870932+07	\N	active	2026-07-30 10:23:56.870932+07	2026-07-30 10:23:56.870932+07
cbf419f2-9a86-4729-bf4a-4973aeb4afc2	fbb44c95-79aa-421e-ba3c-b4e7b87668d8	a0c14f0b-9d06-4a9d-83b2-2ed8dff7ceca	4601ad1b-4dc9-4727-b229-43cdf6fddecb	Paket Percobaan	f	f	2	0	2	2026-07-30 07:00:00+07	\N	used	2026-07-30 19:00:14.380197+07	2026-07-30 19:07:19.617401+07
7624468c-8c9b-405f-9fcd-8d2a2b8f17ad	3c09cb50-3d99-4ff8-a57f-44fd4135d070	69b0f270-ff15-4352-ba8d-f7f383260fca	4601ad1b-4dc9-4727-b229-43cdf6fddecb	Paket Percobaan	f	f	2	2	0	2026-08-04 07:00:00+07	\N	active	2026-08-04 10:23:20.739014+07	2026-08-04 10:23:20.739014+07
51b0bb38-3259-4e2b-951e-27e7d5608d70	3c09cb50-3d99-4ff8-a57f-44fd4135d070	a723de49-a8e2-4d78-884f-203dd4522910	4601ad1b-4dc9-4727-b229-43cdf6fddecb	Paket Percobaan	f	f	2	2	0	2026-08-04 07:00:00+07	\N	active	2026-08-04 10:33:04.574022+07	2026-08-04 10:33:04.574022+07
ff6577f6-2375-4cc0-abd7-a64cec8b79f2	fead48e2-b1b6-4620-9b18-889e64355698	9737a24d-0d15-48ed-897f-c0df00d57b31	4601ad1b-4dc9-4727-b229-43cdf6fddecb	Paket Percobaan	f	f	2	0	2	2026-08-04 07:00:00+07	\N	used	2026-08-04 10:36:39.470359+07	2026-08-04 10:38:40.226653+07
dbaedb67-dc5a-4839-8a9c-b15166650baf	3c09cb50-3d99-4ff8-a57f-44fd4135d070	5ed13eae-af47-486d-98a2-0160e99c5290	4601ad1b-4dc9-4727-b229-43cdf6fddecb	Paket Percobaan	f	f	2	0	2	2026-08-04 07:00:00+07	\N	used	2026-08-04 10:51:50.830608+07	2026-08-04 10:55:31.495687+07
a1d8731c-75c5-45cd-bffa-6d7abf1ec1e3	447c7653-04dc-4550-a703-6c81847b1dd6	dd8a6ee9-e2d5-4e2f-820f-256a843d429d	11111111-1111-1111-1111-111111111101	6x	f	f	6	2	4	2026-07-29 07:00:00+07	\N	active	2026-07-29 11:41:11.037071+07	2026-07-29 11:41:11.037071+07
\.


--
-- Data for Name: instructor_areas; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.instructor_areas (instructor_id, area_type, area_id, id, created_at, updated_at) FROM stdin;
81c28015-a5a8-4f66-82ad-e3b204a88550		1	b2cbe81b-ee50-4a2c-bf65-a477e0433d0f	\N	\N
81c28015-a5a8-4f66-82ad-e3b204a88550		2	aaf4e6e3-cca6-4b1e-90a0-a26c5c476892	\N	\N
81c28015-a5a8-4f66-82ad-e3b204a88550		3	c63012da-5ff7-4dc3-b641-63ae1e87149b	\N	\N
81c28015-a5a8-4f66-82ad-e3b204a88550		4	7c58e2ff-888c-40f0-bb93-3c7f7f783f82	\N	\N
81c28015-a5a8-4f66-82ad-e3b204a88550		5	04bdbf07-3fdb-4a17-b087-cf14a3722c7a	\N	\N
4ea68344-c0c9-4a76-bb2f-2356154887b9		6	fe55a7b4-ef26-451c-b37d-a545ebb35bfa	\N	\N
4ea68344-c0c9-4a76-bb2f-2356154887b9		7	cf8a71bb-7be9-4882-b39e-f857b34fc1d3	\N	\N
4ea68344-c0c9-4a76-bb2f-2356154887b9		8	7b5e1103-af6b-49dd-93d8-6daa5b6ba11d	\N	\N
162ea8ed-9505-4cdc-b1ae-b1bf0c5a6fef		1	fce73bae-b2ee-4003-a50f-d00cfa7c0ed7	\N	\N
162ea8ed-9505-4cdc-b1ae-b1bf0c5a6fef		2	ece6614a-6dc1-4f07-a45f-e7a4a7e3fd6d	\N	\N
162ea8ed-9505-4cdc-b1ae-b1bf0c5a6fef		3	3ca0bac3-7a39-4bac-ad06-f416cccd15de	\N	\N
162ea8ed-9505-4cdc-b1ae-b1bf0c5a6fef		4	b0db4204-b2f4-4131-9567-ac860571a463	\N	\N
162ea8ed-9505-4cdc-b1ae-b1bf0c5a6fef		5	5f373c47-dbfb-4c45-bf47-811313af9295	\N	\N
e8a0b2eb-834b-47ff-85e0-35455885ec7c		6	2e486f66-6c22-41c0-92ca-da3d13206910	\N	\N
e8a0b2eb-834b-47ff-85e0-35455885ec7c		7	df633b3b-a151-4d25-a0d2-550c05be22ed	\N	\N
e8a0b2eb-834b-47ff-85e0-35455885ec7c		8	317379c4-2d0d-4af6-82b6-c87edda0744e	\N	\N
5e22f25f-9248-4c1f-a086-faeb657510c9		1	40970712-2e6f-4b1a-b75b-6d8c3b32bf94	\N	\N
5e22f25f-9248-4c1f-a086-faeb657510c9		2	70016db4-6799-4fc3-af84-c2e73e20da8d	\N	\N
5e22f25f-9248-4c1f-a086-faeb657510c9		3	ac9f5c83-ab62-4fc9-9030-ba259a4743d6	\N	\N
5e22f25f-9248-4c1f-a086-faeb657510c9		4	f104c2de-920f-4d4c-9384-34944c094deb	\N	\N
5e22f25f-9248-4c1f-a086-faeb657510c9		5	11532b5d-e864-40a4-82e2-0ef8c3634b3b	\N	\N
3f818e72-6d7a-4552-ba97-10d4540c1257		6	dce3034d-15d0-4904-bd76-14ed167ef9e3	\N	\N
3f818e72-6d7a-4552-ba97-10d4540c1257		7	c8731c48-0a0c-41d0-aec2-0cfda979ee7f	\N	\N
3f818e72-6d7a-4552-ba97-10d4540c1257		8	044fc357-d8e4-47de-94a4-7238c394b407	\N	\N
5e22f25f-9248-4c1f-a086-faeb657510c9	district	1	8b7047d3-a22a-409d-ba10-c37c7f8947f5	2026-06-10 11:05:57.965366+07	2026-06-10 11:05:57.965366+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	2	d460824b-ff83-4aca-8e85-650681ff576f	2026-06-10 11:05:57.968423+07	2026-06-10 11:05:57.968423+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	3	15cec9f1-3c94-4176-9dff-e2dccfcc624e	2026-06-10 11:05:57.971158+07	2026-06-10 11:05:57.971158+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	4	1669a33c-e2e2-41bb-9b03-ae12bd5863e3	2026-06-10 11:05:57.973796+07	2026-06-10 11:05:57.973796+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	5	a674f33f-7331-42e2-bfec-fc2bd78f436f	2026-06-10 11:05:57.978674+07	2026-06-10 11:05:57.978674+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	6	1f9a3cac-1a3b-46f9-8200-ff51ad4412fc	2026-06-10 11:05:57.983895+07	2026-06-10 11:05:57.983895+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	7	517d5cfa-38ba-494b-b971-6a74b0e7d6cf	2026-06-10 11:05:57.986882+07	2026-06-10 11:05:57.986882+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	8	7574f834-a341-492a-80c7-0ca818b29aca	2026-06-10 11:05:57.989971+07	2026-06-10 11:05:57.989971+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	ae3e9c22-860b-43cf-a826-c90eecaf2917	2026-06-10 11:05:57.995808+07	2026-06-10 11:05:57.995808+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	811812c9-a295-40fa-9ea4-c8fb37ba4a0f	2026-06-10 11:05:58.00009+07	2026-06-10 11:05:58.00009+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	015c81c3-965e-4df9-83a6-a0478bfa47a3	2026-06-10 11:05:58.004305+07	2026-06-10 11:05:58.004305+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	f216f62b-948d-4b4a-acc9-b88c8bdbf3bf	2026-06-10 11:05:58.008099+07	2026-06-10 11:05:58.008099+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	acd67383-7d09-4173-9c4b-0b63878e066f	2026-06-10 11:05:58.013532+07	2026-06-10 11:05:58.013532+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	c0c7d84b-07bf-4b59-8bd1-ec1ba2ea6c81	2026-06-10 18:04:30.244416+07	2026-06-10 18:04:30.244416+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	cbfd7a9d-573c-44b3-8bab-76f43e481e1d	2026-06-10 18:04:30.257516+07	2026-06-10 18:04:30.257516+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	bff397b3-5f56-4579-b70e-9a6778f2258d	2026-06-10 18:04:30.263185+07	2026-06-10 18:04:30.263185+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	ff59d79c-73d0-46b5-a144-041d621ab205	2026-06-10 18:04:30.27443+07	2026-06-10 18:04:30.27443+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	2dd8ceeb-7e9f-497d-a0c3-7ec861f67844	2026-06-10 18:04:30.28497+07	2026-06-10 18:04:30.28497+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	f2168ce2-37d9-4ea0-af2c-b3be793c1c4b	2026-06-10 18:15:57.839815+07	2026-06-10 18:15:57.839815+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	e3489924-1c2a-45f3-ae2c-a47847adeda5	2026-06-10 18:15:57.874506+07	2026-06-10 18:15:57.874506+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	5d66bb21-7b3a-49ec-b8a9-78672e237864	2026-06-10 18:15:57.884895+07	2026-06-10 18:15:57.884895+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	3dfe51e3-999c-44e9-a727-a3353f55e231	2026-06-10 18:15:57.89241+07	2026-06-10 18:15:57.89241+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	3791ce83-dcf6-4371-995b-26802fe45366	2026-06-10 18:15:57.904839+07	2026-06-10 18:15:57.904839+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	80f9233b-9311-480b-9514-a7e87af56a2c	2026-06-11 01:49:35.746099+07	2026-06-11 01:49:35.746099+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	96e5a9e6-1baf-4917-bea3-c65b320cf4c1	2026-06-11 01:49:35.752447+07	2026-06-11 01:49:35.752447+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	6dca7133-f2e1-4505-b98c-6b317095da52	2026-06-11 01:49:35.764514+07	2026-06-11 01:49:35.764514+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	9e418f00-2aa9-42af-940b-4cbef6acc2f6	2026-06-11 01:49:35.775281+07	2026-06-11 01:49:35.775281+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	0525bc56-002a-4843-a9a9-56e90c5eb748	2026-06-11 01:49:35.781169+07	2026-06-11 01:49:35.781169+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	3190c117-416e-427e-916f-05f51d1160ec	2026-06-11 10:01:07.285008+07	2026-06-11 10:01:07.285008+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	ccbc2338-c683-41b1-b646-96f5bcf6ceb0	2026-06-11 10:01:07.29141+07	2026-06-11 10:01:07.29141+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	aec007ac-df1f-44c4-ba93-058030171944	2026-06-11 10:01:07.29599+07	2026-06-11 10:01:07.29599+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	2e18f56c-d581-4345-b205-4e139f47bed6	2026-06-11 10:01:07.30062+07	2026-06-11 10:01:07.30062+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	f9a4bb6b-59dc-4a0a-8590-651aa00f6a7d	2026-06-11 10:01:07.306238+07	2026-06-11 10:01:07.306238+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	43190d69-e27f-40eb-a34d-526d94e791bc	2026-06-11 10:09:27.232934+07	2026-06-11 10:09:27.232934+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	74615ab4-214b-4e65-9b79-ff2a017ceeef	2026-06-11 10:09:27.263279+07	2026-06-11 10:09:27.263279+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	66fb86d4-bd2b-4be3-b8dd-d7a2e0a38041	2026-06-11 10:09:27.269001+07	2026-06-11 10:09:27.269001+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	5b79b012-f59e-4b71-8b70-8144174ad50a	2026-06-11 10:09:27.280218+07	2026-06-11 10:09:27.280218+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	47252444-db41-40f5-a26a-bd5db5e9ed6a	2026-06-11 10:09:27.282525+07	2026-06-11 10:09:27.282525+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	af1fbf56-dd1f-4a03-8c55-d051c33b92ca	2026-06-11 10:12:15.205589+07	2026-06-11 10:12:15.205589+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	89b9f77f-303d-4f03-b9d1-cae73301fb2a	2026-06-11 10:12:15.210499+07	2026-06-11 10:12:15.210499+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	31cf86b7-d524-4d50-a0b0-568da522d88d	2026-06-11 10:12:15.213195+07	2026-06-11 10:12:15.213195+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	87f6ca93-1ab1-4ee8-bdfe-bb9d00396b6b	2026-06-11 10:12:15.215783+07	2026-06-11 10:12:15.215783+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	415dcd1c-f3dd-4815-8b55-7b13c5208bbf	2026-06-11 10:12:15.217801+07	2026-06-11 10:12:15.217801+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	516f7db9-b02f-4475-a5f9-a3615b1996b2	2026-06-11 10:13:28.569476+07	2026-06-11 10:13:28.569476+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	20f6d1c5-c73f-4c32-a2f1-b1cab8bf9ecc	2026-06-11 10:13:28.575426+07	2026-06-11 10:13:28.575426+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	d98f4850-4a30-495d-811b-0bb8599295f0	2026-06-11 10:13:28.586323+07	2026-06-11 10:13:28.586323+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	9b76b649-1013-471f-91b7-b1311093e663	2026-06-11 10:13:28.644475+07	2026-06-11 10:13:28.644475+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	7080294b-d478-4016-a24d-f78891c2360b	2026-06-11 10:13:28.648128+07	2026-06-11 10:13:28.648128+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	ac98cb42-3d8c-4ddf-bc1b-12646a640451	2026-06-11 10:19:15.10205+07	2026-06-11 10:19:15.10205+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	7460c558-061d-4836-81ca-59c597733ab8	2026-06-11 10:19:15.106755+07	2026-06-11 10:19:15.106755+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	4929722a-56f2-4486-9a7a-92799a101fbc	2026-06-11 10:19:15.108849+07	2026-06-11 10:19:15.108849+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	6a2b9c4f-e7bd-494f-8d61-c2415f58f12f	2026-06-11 10:19:15.110634+07	2026-06-11 10:19:15.110634+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	7197847f-42bb-40a6-8db7-358924b686c3	2026-06-11 10:19:15.112458+07	2026-06-11 10:19:15.112458+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	fa441eae-de04-4f53-9387-d927c38abbac	2026-06-11 10:22:26.166431+07	2026-06-11 10:22:26.166431+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	d104cf9a-d027-4c1a-a184-c94a7525bc5e	2026-06-11 10:22:26.171334+07	2026-06-11 10:22:26.171334+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	5dae5fc9-94ff-4448-9ba1-488eac076682	2026-06-11 10:22:26.173931+07	2026-06-11 10:22:26.173931+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	0a9baf18-e7e7-4eea-9e66-11eedabbc0a6	2026-06-11 10:22:26.17676+07	2026-06-11 10:22:26.17676+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	a78bb19a-9c3c-4450-b0cc-d0af342754d5	2026-06-11 10:22:26.17985+07	2026-06-11 10:22:26.17985+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	8d80955c-d892-4b60-bfaa-b9f2df7538d8	2026-06-11 11:13:55.208372+07	2026-06-11 11:13:55.208372+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	0669e233-db7c-4bc2-8d2b-5a66fde5fade	2026-06-11 11:13:55.216334+07	2026-06-11 11:13:55.216334+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	02e4441e-d0d0-42aa-97a2-86380fba8349	2026-06-11 11:13:55.219531+07	2026-06-11 11:13:55.219531+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	3d978d3b-e8c1-46f4-b1f2-395dd8ec882c	2026-06-11 11:13:55.222626+07	2026-06-11 11:13:55.222626+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	4db417b7-b01a-469d-9a1a-84d83def3b98	2026-06-11 11:13:55.225137+07	2026-06-11 11:13:55.225137+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	ad832b54-288f-4c71-9ea9-ce64fed03c62	2026-06-11 11:18:54.458737+07	2026-06-11 11:18:54.458737+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	c49a46b7-6b8e-4cc8-b9b7-554b6ac1155c	2026-06-11 11:18:54.471871+07	2026-06-11 11:18:54.471871+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	788f325e-3226-4ae6-b931-526ad8dfda5e	2026-06-11 11:18:54.481255+07	2026-06-11 11:18:54.481255+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	8c69ecde-38e9-4f26-9f89-6ee4a9a7e417	2026-06-11 11:18:54.487553+07	2026-06-11 11:18:54.487553+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	d928f367-e774-4893-b980-0d3a16c65911	2026-06-11 11:18:54.490784+07	2026-06-11 11:18:54.490784+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	3ea72a9c-d92d-440c-8aee-f515246a4d77	2026-06-11 11:26:32.423111+07	2026-06-11 11:26:32.423111+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	09bd08c8-5906-4731-ace9-6f0109ccda44	2026-06-11 11:26:32.433776+07	2026-06-11 11:26:32.433776+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	07c20476-21d6-4add-aa01-6e3ab7a9e3c2	2026-06-11 11:26:32.443311+07	2026-06-11 11:26:32.443311+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	c0df5255-f4a8-43ab-809c-5c10702dd073	2026-06-11 11:26:32.467382+07	2026-06-11 11:26:32.467382+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	d0ac7689-a385-46da-bb95-631335b7ca85	2026-06-11 11:26:32.486528+07	2026-06-11 11:26:32.486528+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	1b6e2d27-da09-4484-a07d-937b47cca50f	2026-06-11 11:35:04.410399+07	2026-06-11 11:35:04.410399+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	62f52cae-a562-4d1a-985a-8c8a988334d8	2026-06-11 11:35:04.437601+07	2026-06-11 11:35:04.437601+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	418060e7-ce5f-4b84-9e26-37969cb1f250	2026-06-11 11:35:04.447437+07	2026-06-11 11:35:04.447437+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	cc21177c-5a9f-40b7-a79c-66c425eb5ad4	2026-06-11 11:35:04.45755+07	2026-06-11 11:35:04.45755+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	5140560e-ed60-4d14-8357-1376d284196e	2026-06-11 11:35:04.463473+07	2026-06-11 11:35:04.463473+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	6b09af2d-ba7c-4c76-8f78-544062d88702	2026-06-11 14:14:23.852304+07	2026-06-11 14:14:23.852304+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	9b7fc3f0-22f1-4a5b-83cb-e29aecfdf48a	2026-06-11 14:14:23.875733+07	2026-06-11 14:14:23.875733+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	58e8121d-668d-4f42-9440-bd0e40ff66b7	2026-06-11 14:14:23.879572+07	2026-06-11 14:14:23.879572+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	6485bc30-400c-490c-896f-f68ce4fdde28	2026-06-11 14:14:23.885945+07	2026-06-11 14:14:23.885945+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	1501478b-9c4e-4917-b82c-6ebca4ecf7af	2026-06-11 14:14:23.90204+07	2026-06-11 14:14:23.90204+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	f792eeea-7ae7-4b0b-a49f-7c32435926d2	2026-06-11 14:18:26.042948+07	2026-06-11 14:18:26.042948+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	c5f3b4d5-fef8-48d0-a833-488e3c662866	2026-06-11 14:18:26.050229+07	2026-06-11 14:18:26.050229+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	7a8b3a88-7ad1-424f-8ce5-704107127443	2026-06-11 14:18:26.055511+07	2026-06-11 14:18:26.055511+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	f78184e4-99c5-4f86-b9ca-33e056a0473d	2026-06-11 14:18:26.062408+07	2026-06-11 14:18:26.062408+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	7a39b89e-a7d6-4258-b03c-040af6c67d6e	2026-06-11 14:18:26.06806+07	2026-06-11 14:18:26.06806+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	83140bad-bff3-4252-9f6a-29208c7546a3	2026-06-11 14:31:03.799767+07	2026-06-11 14:31:03.799767+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	3c33b10f-266c-49c3-9731-2ed27bb9c939	2026-06-11 14:31:03.808695+07	2026-06-11 14:31:03.808695+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	c9dd7808-4630-4058-a591-13b114adaa02	2026-06-11 14:31:03.815575+07	2026-06-11 14:31:03.815575+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	b909b05f-e1d2-4ea4-86dc-b76a9b69f25c	2026-06-11 14:31:03.822194+07	2026-06-11 14:31:03.822194+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	e6d90a8d-e612-45d9-a6c3-9db9cd264715	2026-06-11 14:31:03.831546+07	2026-06-11 14:31:03.831546+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	2fe76b90-b1aa-4502-a368-54f4cd6326d7	2026-06-11 16:17:02.118513+07	2026-06-11 16:17:02.118513+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	99acf44e-45e9-4140-b83d-971c88c1c9e4	2026-06-11 16:17:02.124368+07	2026-06-11 16:17:02.124368+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	9c432069-6550-4d56-9cb7-bf5f22d68051	2026-06-11 16:17:02.126664+07	2026-06-11 16:17:02.126664+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	261cb752-2076-459c-9964-4b78af046611	2026-06-11 16:17:02.131287+07	2026-06-11 16:17:02.131287+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	c29663ac-1da2-487f-b8a2-f1d4681a0836	2026-06-11 16:17:02.138072+07	2026-06-11 16:17:02.138072+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	b5b4453f-3821-43b6-9828-1d26983ae0bf	2026-06-11 16:24:06.767049+07	2026-06-11 16:24:06.767049+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	8f2c4d00-0fd8-4292-9f8c-24c9046e1fbf	2026-06-11 16:24:06.776398+07	2026-06-11 16:24:06.776398+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	60ec0443-feda-4955-a1d2-05383741e1a1	2026-06-11 16:24:06.79233+07	2026-06-11 16:24:06.79233+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	5b199c3a-abda-438e-8c55-6ddb98fd6c75	2026-06-11 16:24:06.79946+07	2026-06-11 16:24:06.79946+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	d80bdb12-296f-4acf-9bb8-4a7b03d12b4f	2026-06-11 16:24:06.808894+07	2026-06-11 16:24:06.808894+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	bcf08482-d2c4-4f6e-a757-a1ea415d031a	2026-06-11 17:08:16.480174+07	2026-06-11 17:08:16.480174+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	0c51c3bc-e9e4-4bd1-8a6c-c2b4a6128a1c	2026-06-11 17:08:16.485165+07	2026-06-11 17:08:16.485165+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	c1d03646-f48b-4075-9418-99e022276b39	2026-06-11 17:08:16.48721+07	2026-06-11 17:08:16.48721+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	df586c67-3535-42dc-b2c0-c6d3956a74aa	2026-06-11 17:08:16.489098+07	2026-06-11 17:08:16.489098+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	9c61ae03-f846-417e-8d93-53f197e3c561	2026-06-11 17:08:16.491198+07	2026-06-11 17:08:16.491198+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	aca43802-93a5-47ae-ad0c-d0738303c2a0	2026-06-11 17:11:36.165427+07	2026-06-11 17:11:36.165427+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	5df6f518-1cf9-4343-a8ae-58d6a3d9d47d	2026-06-11 17:11:36.168924+07	2026-06-11 17:11:36.168924+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	333cc884-e936-4014-8591-08cb462ea197	2026-06-11 17:11:36.17231+07	2026-06-11 17:11:36.17231+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	09567d36-99a1-4f6e-b33c-b1fa81d1b4cc	2026-06-11 17:11:36.174926+07	2026-06-11 17:11:36.174926+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	696bb29c-ff38-4b8b-981b-9da66355828c	2026-06-11 17:11:36.178576+07	2026-06-11 17:11:36.178576+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	9491b293-38d5-4fef-bea7-868feceec8d5	2026-06-11 17:13:03.704496+07	2026-06-11 17:13:03.704496+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	70b307d7-d1c4-432c-86dc-8c47f2f7e69b	2026-06-11 17:13:03.710169+07	2026-06-11 17:13:03.710169+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	e1984308-9d3b-4915-9b7c-bec2cc598cff	2026-06-11 17:13:03.713479+07	2026-06-11 17:13:03.713479+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	10eb188a-f9e0-4f71-87e9-427a17cbb03a	2026-06-11 17:13:03.718439+07	2026-06-11 17:13:03.718439+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	6905e132-3680-4d8e-99b6-472aad62aa61	2026-06-11 17:13:03.721011+07	2026-06-11 17:13:03.721011+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	46958ec5-e2f8-4aaf-8e6c-77b46abb47e6	2026-06-11 17:37:30.365+07	2026-06-11 17:37:30.365+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	e38a9b32-d308-4f46-89a0-41d7297b90d5	2026-06-11 17:37:30.370679+07	2026-06-11 17:37:30.370679+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	d0195a64-2aa9-4730-b1c2-866e5736a039	2026-06-11 17:37:30.37329+07	2026-06-11 17:37:30.37329+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	91f375e5-a32d-48be-b85f-004cdad3138f	2026-06-11 17:37:30.376659+07	2026-06-11 17:37:30.376659+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	01d73664-2201-44b5-9cd7-00f4ac7089db	2026-06-11 17:37:30.379645+07	2026-06-11 17:37:30.379645+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	16c56fd8-c024-4973-b5b6-d155695ac2fa	2026-06-11 17:58:06.999344+07	2026-06-11 17:58:06.999344+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	b393a5b0-6d57-4cf2-b63e-a5c48aaf94ff	2026-06-11 17:58:07.003359+07	2026-06-11 17:58:07.003359+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	4c17b8a2-34f2-4e63-9b1c-a2b20974334f	2026-06-11 17:58:07.007459+07	2026-06-11 17:58:07.007459+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	7c937f82-1219-4b79-a087-a88ef57aa495	2026-06-11 17:58:07.011483+07	2026-06-11 17:58:07.011483+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	c8dd93d7-7b8b-46b8-a604-00998081de6d	2026-06-11 17:58:07.015621+07	2026-06-11 17:58:07.015621+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	8ffbdbfc-87d9-4e5b-9b91-5cb2d25285ec	2026-06-12 05:58:52.518884+07	2026-06-12 05:58:52.518884+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	15141e2c-d5b4-4bc4-a9fe-1c718910916c	2026-06-12 05:58:52.530552+07	2026-06-12 05:58:52.530552+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	7386cc50-a0c8-4a34-bfd0-b2c86c35620f	2026-06-12 05:58:52.53664+07	2026-06-12 05:58:52.53664+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	e84a29be-d17f-4517-9660-06c096605a99	2026-06-12 05:58:52.542328+07	2026-06-12 05:58:52.542328+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	a849a0e0-cad4-4f9c-a14c-f0f360bbfd8f	2026-06-12 05:58:52.55032+07	2026-06-12 05:58:52.55032+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	8ed35c45-6369-4a5c-8ad6-7be466c63577	2026-06-12 14:26:16.412141+07	2026-06-12 14:26:16.412141+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	e537c1af-cd3b-4a88-bbd7-b87424e6aa2c	2026-06-12 14:26:16.420229+07	2026-06-12 14:26:16.420229+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	025b0fa4-630a-4135-942a-5362b3e2cda6	2026-06-12 14:26:16.429714+07	2026-06-12 14:26:16.429714+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	3d0a0863-a26c-4b6a-8126-a4bc172ea09f	2026-06-12 14:26:16.432629+07	2026-06-12 14:26:16.432629+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	591a065c-766b-41e4-bcf8-94f1793fc6d0	2026-06-12 14:26:16.436604+07	2026-06-12 14:26:16.436604+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	33d604c5-405c-46cf-9167-c627275b8a09	2026-06-13 11:30:50.961611+07	2026-06-13 11:30:50.961611+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	7eed52da-e4c2-4b2b-bed8-4341b6cc0217	2026-06-13 11:30:50.966997+07	2026-06-13 11:30:50.966997+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	7129e2de-2192-40a0-950e-6941a985418a	2026-06-13 11:30:50.97519+07	2026-06-13 11:30:50.97519+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	79009ad6-3a44-4719-8d56-918ee5e06f83	2026-06-13 11:30:50.97797+07	2026-06-13 11:30:50.97797+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	6be360eb-1be2-4bb7-b783-64c6305b2204	2026-06-13 11:30:50.981298+07	2026-06-13 11:30:50.981298+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	0f0a0bab-1808-409a-8299-0b1c00ffd8bc	2026-06-13 12:00:32.736933+07	2026-06-13 12:00:32.736933+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	461c4b3a-baf6-4434-b8e2-24ca82891cab	2026-06-13 12:00:32.740929+07	2026-06-13 12:00:32.740929+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	6fbc9544-18e4-4c64-8bef-cd49055bbc0d	2026-06-13 12:00:32.745155+07	2026-06-13 12:00:32.745155+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	6b66e4d4-c275-4f71-b570-bc665a37db44	2026-06-13 12:00:32.747266+07	2026-06-13 12:00:32.747266+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	d85a5461-cffa-4310-acc5-9b888009effb	2026-06-13 12:00:32.748944+07	2026-06-13 12:00:32.748944+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	4984ac71-a615-48fe-a11e-2e210aa45333	2026-06-13 13:14:29.124156+07	2026-06-13 13:14:29.124156+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	436d13b6-6ecb-4f67-b829-9362c08d3475	2026-06-13 13:14:29.131864+07	2026-06-13 13:14:29.131864+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	157be8b8-d75b-4c09-932c-094d0023db9f	2026-06-13 13:14:29.136185+07	2026-06-13 13:14:29.136185+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	4982c056-7ef4-42be-adcb-6dfb279793cb	2026-06-13 13:14:29.140973+07	2026-06-13 13:14:29.140973+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	49ee2c3b-a867-4514-91c9-c0b64c1b7a26	2026-06-13 13:14:29.145074+07	2026-06-13 13:14:29.145074+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	e4bb5a38-6ee0-48f6-a2ad-51c43f63dd1c	2026-06-13 13:18:54.607688+07	2026-06-13 13:18:54.607688+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	a7adda99-6fac-4a3f-97ef-e13663644ee0	2026-06-13 13:18:54.610662+07	2026-06-13 13:18:54.610662+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	641039be-184b-4762-b242-eb844157507d	2026-06-13 13:18:54.615014+07	2026-06-13 13:18:54.615014+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	7d46c285-e84b-46ce-b8a4-942b2d6e8f7f	2026-06-13 13:18:54.617741+07	2026-06-13 13:18:54.617741+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	19287010-c492-442a-a0ba-183ec3e7c5cb	2026-06-13 13:18:54.620023+07	2026-06-13 13:18:54.620023+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	ec688141-8638-409b-9915-0fff2598247d	2026-06-13 13:26:13.740702+07	2026-06-13 13:26:13.740702+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	979e59ca-009d-45cb-a61b-92fc44dd052a	2026-06-13 13:26:13.805678+07	2026-06-13 13:26:13.805678+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	e6e64f05-8574-4d64-b240-3138f9517243	2026-06-13 13:26:13.959672+07	2026-06-13 13:26:13.959672+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	bedd6f3f-5db8-49f3-8c69-55040b35e2cd	2026-06-13 13:26:13.970422+07	2026-06-13 13:26:13.970422+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	3c9d89fb-f2fd-46ab-8964-a4dd0490b5a4	2026-06-13 13:26:14.02725+07	2026-06-13 13:26:14.02725+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	84e7c211-bb05-436e-8e9d-365e0cdc7d82	2026-06-13 13:34:33.723745+07	2026-06-13 13:34:33.723745+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	bd66e278-3080-4079-850d-d2d7cb900d5e	2026-06-13 13:34:33.735665+07	2026-06-13 13:34:33.735665+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	c4947d76-ff45-46c4-9719-3b2d79d2bcbc	2026-06-13 13:34:33.746709+07	2026-06-13 13:34:33.746709+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	d8009cfd-1d56-41e7-823a-7f920fa94044	2026-06-13 13:34:33.75838+07	2026-06-13 13:34:33.75838+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	63d9e7db-0e88-4d6a-b5eb-d6302886f28e	2026-06-13 13:34:33.777455+07	2026-06-13 13:34:33.777455+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	7d408a2c-d7ed-4025-a67e-b259cd9a451e	2026-06-13 13:44:46.033211+07	2026-06-13 13:44:46.033211+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	857b6f60-0ddb-445a-98f7-1b13e5a7881e	2026-06-13 13:44:46.039923+07	2026-06-13 13:44:46.039923+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	bfb7845a-e4ec-4ec0-a5df-7024186a5319	2026-06-13 13:44:46.043506+07	2026-06-13 13:44:46.043506+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	642bd6e3-c5cd-44b7-af20-8af784d22188	2026-06-13 13:44:46.054765+07	2026-06-13 13:44:46.054765+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	3bacb803-9579-4d0d-b986-8d3a6971b1a1	2026-06-13 13:44:46.060574+07	2026-06-13 13:44:46.060574+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	1e501355-0592-489e-b0cb-8076541597e1	2026-06-13 13:51:50.356599+07	2026-06-13 13:51:50.356599+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	df4cacb1-ac00-49a5-bcf3-18caec5eba01	2026-06-13 13:51:50.361506+07	2026-06-13 13:51:50.361506+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	4a30c833-aae1-4184-ad28-4d82be05fb1b	2026-06-13 13:51:50.363392+07	2026-06-13 13:51:50.363392+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	30dbf3c7-c271-430d-a1b8-31f607281f76	2026-06-13 13:51:50.365787+07	2026-06-13 13:51:50.365787+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	193ceabc-d2ac-4286-bb9d-63013b53ff6d	2026-06-13 13:51:50.368523+07	2026-06-13 13:51:50.368523+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	bc687f1f-b80e-4e88-9e91-766bd6bb3c8a	2026-06-13 14:09:38.900715+07	2026-06-13 14:09:38.900715+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	d9a52c9e-9bb3-4edf-8d46-8064592de554	2026-06-13 14:09:38.906333+07	2026-06-13 14:09:38.906333+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	bc71bf0a-12ea-41b4-8d49-d8f4a2de1174	2026-06-13 14:09:38.909274+07	2026-06-13 14:09:38.909274+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	59dbb001-1bcb-47fb-a8f5-e24889daee44	2026-06-13 14:09:38.918997+07	2026-06-13 14:09:38.918997+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	adf97c28-d6c5-4845-b3a7-435e9b95aeaf	2026-06-13 14:09:38.937482+07	2026-06-13 14:09:38.937482+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	0987b1ae-2a0a-4d46-901d-438b598f6784	2026-06-13 15:06:45.953385+07	2026-06-13 15:06:45.953385+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	d20d599f-a067-4851-be84-a3e5d3236487	2026-06-13 15:06:45.95623+07	2026-06-13 15:06:45.95623+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	07dd45fb-4adf-4ed7-bf8b-ab21ccd0b703	2026-06-13 15:06:45.95869+07	2026-06-13 15:06:45.95869+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	ce9a327c-cc25-4809-895d-d64b6ac8e037	2026-06-13 15:06:45.960718+07	2026-06-13 15:06:45.960718+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	c7add18d-ffc8-4bf3-9d03-fcb556c8c3b7	2026-06-13 15:06:45.963911+07	2026-06-13 15:06:45.963911+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	f7ebeca7-64b2-4b03-bf72-d8a2d22637b7	2026-06-13 15:34:41.542624+07	2026-06-13 15:34:41.542624+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	0a3123be-6bb8-4a9e-90fb-ca656727d8b2	2026-06-13 15:34:41.564001+07	2026-06-13 15:34:41.564001+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	195dabe6-1a8c-4a8a-af17-a98ab4f1ebdb	2026-06-13 15:34:41.567484+07	2026-06-13 15:34:41.567484+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	108662b3-7129-475b-96bd-c9cfa21fd4b6	2026-06-13 15:34:41.570573+07	2026-06-13 15:34:41.570573+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	b3b90079-2493-4949-a959-ebc3bdd06514	2026-06-13 15:34:41.575803+07	2026-06-13 15:34:41.575803+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	02a0032a-b81a-449c-b63c-e599e8c0933b	2026-06-13 15:53:55.233158+07	2026-06-13 15:53:55.233158+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	793d9c26-8ce3-4c12-939d-8f176800759d	2026-06-13 15:53:55.238034+07	2026-06-13 15:53:55.238034+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	5b221845-0b9c-4510-b5d5-9aaa1da8f197	2026-06-13 15:53:55.243686+07	2026-06-13 15:53:55.243686+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	80704747-9e74-4d38-8146-d9d0ca28a734	2026-06-13 15:53:55.248581+07	2026-06-13 15:53:55.248581+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	a710b720-f7f0-49c4-b55c-43d60f8670bb	2026-06-13 15:53:55.252145+07	2026-06-13 15:53:55.252145+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	23406a93-85b9-4c32-8b31-dd84ae850695	2026-06-13 16:10:10.025453+07	2026-06-13 16:10:10.025453+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	b56ae5dd-148a-4c17-8a46-c615c99b7c0f	2026-06-13 16:10:10.03648+07	2026-06-13 16:10:10.03648+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	ab7bc4b8-7892-44a2-8e67-67effd0e6d6b	2026-06-13 16:10:10.039425+07	2026-06-13 16:10:10.039425+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	cfc33791-aeef-41ee-af81-c29e1c80e070	2026-06-13 16:10:10.048414+07	2026-06-13 16:10:10.048414+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	9cdacde7-22fc-4e42-973a-1da552ea231a	2026-06-13 16:10:10.052745+07	2026-06-13 16:10:10.052745+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	110f92d6-b5a0-4d62-8bfd-dd125a354133	2026-06-13 16:28:39.313202+07	2026-06-13 16:28:39.313202+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	3fe71999-1991-478b-bab5-bff40422137c	2026-06-13 16:28:39.320471+07	2026-06-13 16:28:39.320471+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	5c1db5a2-b1bd-4ef0-9828-4e3b0b272458	2026-06-13 16:28:39.322845+07	2026-06-13 16:28:39.322845+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	40d737c1-ba49-4746-bdd7-7af88dcc7030	2026-06-13 16:28:39.326361+07	2026-06-13 16:28:39.326361+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	ba2f7fc9-a05c-433b-bb96-c9fa4d61d408	2026-06-13 16:28:39.328485+07	2026-06-13 16:28:39.328485+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	4a948456-cfa7-4230-9cf6-240f51283380	2026-06-13 16:31:12.194093+07	2026-06-13 16:31:12.194093+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	caf1de1e-7dcb-49a9-bc4e-faa2f88f830f	2026-06-13 16:31:12.203591+07	2026-06-13 16:31:12.203591+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	49153d66-112d-4a2a-b3a6-1c2f77d5920b	2026-06-13 16:31:12.209374+07	2026-06-13 16:31:12.209374+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	a9b099ea-eba2-4ba5-a021-275ff9cce6de	2026-06-13 16:31:12.212892+07	2026-06-13 16:31:12.212892+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	91f5feb4-0221-4ce7-897d-ce2045bef6c0	2026-06-13 16:31:12.218683+07	2026-06-13 16:31:12.218683+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	b8002659-9486-4e74-84d9-571235db7d76	2026-06-13 17:17:03.186313+07	2026-06-13 17:17:03.186313+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	cdf9c49a-f0ac-4455-95bd-ae2b4d5daf49	2026-06-13 17:17:03.196972+07	2026-06-13 17:17:03.196972+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	18509dc6-7c8c-4077-a0a2-f45de9c637d2	2026-06-13 17:17:03.200279+07	2026-06-13 17:17:03.200279+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	5ed2d312-e109-472b-a1f3-196d5f1404e2	2026-06-13 17:17:03.203993+07	2026-06-13 17:17:03.203993+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	c91d5843-9507-46ba-a94e-b7ba96295e8d	2026-06-13 17:17:03.210481+07	2026-06-13 17:17:03.210481+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	6dc0e1a9-59a1-4be0-9dda-5ae5352a7e05	2026-06-13 20:42:59.82089+07	2026-06-13 20:42:59.82089+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	9800b45b-e309-4d43-915d-ae68ad74a331	2026-06-13 20:42:59.831929+07	2026-06-13 20:42:59.831929+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	2f53472f-8ae2-44f7-ac70-af5dbb68d0ff	2026-06-13 20:42:59.837349+07	2026-06-13 20:42:59.837349+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	14e2ab7f-8746-4f04-9862-13e91e247f03	2026-06-13 20:42:59.845241+07	2026-06-13 20:42:59.845241+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	1ca1516f-4d6c-4106-a26a-f09026f02f05	2026-06-13 20:42:59.85684+07	2026-06-13 20:42:59.85684+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	bcb77cfa-abac-4d73-b842-785b46e96301	2026-06-14 18:24:33.485181+07	2026-06-14 18:24:33.485181+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	b8ea2d64-a68a-483d-96a1-50dbee59ab81	2026-06-14 18:24:33.489515+07	2026-06-14 18:24:33.489515+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	9b4fee5e-ad4c-4e15-a02b-29c4933628bc	2026-06-14 18:24:33.493339+07	2026-06-14 18:24:33.493339+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	587ce60c-ae8c-4228-b3f0-42f2857849e5	2026-06-14 18:24:33.498258+07	2026-06-14 18:24:33.498258+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	f666225a-92e8-4857-8d4f-fc544420b74f	2026-06-14 18:24:33.501084+07	2026-06-14 18:24:33.501084+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	5008bc8d-aed4-4130-903e-9c2f5dc61269	2026-06-15 05:02:52.145027+07	2026-06-15 05:02:52.145027+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	50dd8d38-8a43-4658-a561-0272a2111d51	2026-06-15 05:02:52.150795+07	2026-06-15 05:02:52.150795+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	3b85fd31-73a3-414c-9fd5-3aecd9539403	2026-06-15 05:02:52.15531+07	2026-06-15 05:02:52.15531+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	c07fa231-fde2-4657-9a3b-95dd7e950ffa	2026-06-15 05:02:52.162752+07	2026-06-15 05:02:52.162752+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	8db9b62c-2c75-46fb-a938-ac850a36b29a	2026-06-15 05:02:52.167338+07	2026-06-15 05:02:52.167338+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	7aefd49b-3dab-4c6a-a980-606fcc925c45	2026-06-15 05:56:24.216769+07	2026-06-15 05:56:24.216769+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	f78f5f08-fa79-4fbc-beba-764611e6ac79	2026-06-15 05:56:24.219905+07	2026-06-15 05:56:24.219905+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	c0ff50a5-6541-474e-b66e-f8581b8b5ab9	2026-06-15 05:56:24.223106+07	2026-06-15 05:56:24.223106+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	4656c10c-ed49-4cb1-9626-0ffcfc5fec76	2026-06-15 05:56:24.226292+07	2026-06-15 05:56:24.226292+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	d30d8661-11f6-49ea-814e-9f580e042c1a	2026-06-15 05:56:24.229608+07	2026-06-15 05:56:24.229608+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	9efdc01f-0fc0-4d97-92ee-f32cf877d1c0	2026-06-15 16:24:01.379484+07	2026-06-15 16:24:01.379484+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	2ee66c87-f118-4f97-b697-96773a419242	2026-06-15 16:24:01.389021+07	2026-06-15 16:24:01.389021+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	b432aea1-3f94-4a22-bdd3-5dd31ef576c4	2026-06-15 16:24:01.392571+07	2026-06-15 16:24:01.392571+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	0e0b90ec-1958-4c11-9143-5734b0e78166	2026-06-15 16:24:01.395404+07	2026-06-15 16:24:01.395404+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	b6d2004e-415f-4fe5-87fe-1d4ac561b4cf	2026-06-15 16:24:01.399802+07	2026-06-15 16:24:01.399802+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	3e275334-cb80-4e87-a3e2-e6f69d6d2a10	2026-06-15 16:34:56.421+07	2026-06-15 16:34:56.421+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	afeef3af-d6ab-4947-aa2b-cc78c45c2386	2026-06-15 16:34:56.429277+07	2026-06-15 16:34:56.429277+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	66540ff6-f0f2-4627-baca-b209cf5e3a77	2026-06-15 16:34:56.432748+07	2026-06-15 16:34:56.432748+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	711bd8e7-9e8f-4972-b8da-5e4ba21bb407	2026-06-15 16:34:56.4392+07	2026-06-15 16:34:56.4392+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	b9ffbfad-d4f2-4fd5-9740-1a98cfa176d9	2026-06-15 16:34:56.445771+07	2026-06-15 16:34:56.445771+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	faa0cd03-2ece-4e47-94a4-5a318b28c789	2026-06-15 17:04:31.711019+07	2026-06-15 17:04:31.711019+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	5c1f27fa-7638-4166-b569-8060c92943f0	2026-06-15 17:04:31.722296+07	2026-06-15 17:04:31.722296+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	11032469-079a-4b8f-8c7e-e9bc133966ac	2026-06-15 17:04:31.731138+07	2026-06-15 17:04:31.731138+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	76e63f69-4063-458a-b7b2-dbafca71e3d8	2026-06-15 17:04:31.734706+07	2026-06-15 17:04:31.734706+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	5c977e84-9f5d-466e-9b45-824a0b85a635	2026-06-15 17:04:31.740412+07	2026-06-15 17:04:31.740412+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	070b48e8-4bcd-420a-bb30-192783610335	2026-06-15 17:08:23.938814+07	2026-06-15 17:08:23.938814+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	0e0e498c-d918-44dc-b318-cba9890eae9d	2026-06-15 17:08:23.945801+07	2026-06-15 17:08:23.945801+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	5c3e815f-3d65-418f-a932-a149b4c301a8	2026-06-15 17:08:23.948943+07	2026-06-15 17:08:23.948943+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	5770ef98-8d3f-47e9-9765-43a9a455e8af	2026-06-15 17:08:23.953253+07	2026-06-15 17:08:23.953253+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	3d7bff8a-ed03-40d0-8ca4-967ce3311741	2026-06-15 17:08:23.959194+07	2026-06-15 17:08:23.959194+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	4b143a4f-63da-4049-99cf-04ea9b3e3a4e	2026-06-16 22:05:56.306862+07	2026-06-16 22:05:56.306862+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	4e051825-7044-4b45-8fe7-5ff2ae2e6253	2026-06-16 22:05:56.351371+07	2026-06-16 22:05:56.351371+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	c0004204-d869-4f99-9e24-3a41ce656203	2026-06-16 22:05:56.365981+07	2026-06-16 22:05:56.365981+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	fcc6184f-a49f-4635-8197-31b2cc2ec599	2026-06-16 22:05:56.374956+07	2026-06-16 22:05:56.374956+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	6fd66f2e-3732-4ca9-bd57-46c459437d95	2026-06-16 22:05:56.383787+07	2026-06-16 22:05:56.383787+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	ea7ed52b-5b61-428f-bb39-f06ebc519a49	2026-06-17 13:18:54.40697+07	2026-06-17 13:18:54.40697+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	9a05bb38-1615-4eb2-96ac-27dce2ac383d	2026-06-17 13:18:54.414979+07	2026-06-17 13:18:54.414979+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	46bff2ec-f808-44d4-9dec-8e9a160e9d8d	2026-06-17 13:18:54.419085+07	2026-06-17 13:18:54.419085+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	456ce43b-7773-4458-b95c-cbd12848ffc3	2026-06-17 13:18:54.423139+07	2026-06-17 13:18:54.423139+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	3454e684-df0e-4bd0-800c-1ce92ba3883f	2026-06-17 13:18:54.426716+07	2026-06-17 13:18:54.426716+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	fbd3730f-5b2d-4b41-af24-fb8bb68ed7ef	2026-06-17 14:39:09.802577+07	2026-06-17 14:39:09.802577+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	73ffddff-632d-407e-9b00-24fe9aad4508	2026-06-17 14:39:09.806275+07	2026-06-17 14:39:09.806275+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	343f3b70-c160-4992-8a6a-a88630040125	2026-06-17 14:39:09.809237+07	2026-06-17 14:39:09.809237+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	6a14240f-c581-49a3-85b8-e3ea90789d3b	2026-06-17 14:39:09.812056+07	2026-06-17 14:39:09.812056+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	3d1c10c8-4004-4fbc-ba01-4fd5e136c637	2026-06-17 14:39:09.814839+07	2026-06-17 14:39:09.814839+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	29adfad4-27ec-4462-ae63-c61c7f398eab	2026-06-17 15:33:43.036839+07	2026-06-17 15:33:43.036839+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	9de19025-fd58-4725-9b49-36c46ff12f11	2026-06-17 15:33:43.045688+07	2026-06-17 15:33:43.045688+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	a305c515-e703-42db-8cfe-b0543e2dd221	2026-06-17 15:33:43.05065+07	2026-06-17 15:33:43.05065+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	e4616b03-beb5-4c76-9749-136660115eb1	2026-06-17 15:33:43.066628+07	2026-06-17 15:33:43.066628+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	d4cd47ec-5fb1-49b8-9d9b-c274a70faff4	2026-06-17 15:33:43.070861+07	2026-06-17 15:33:43.070861+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	175b59f5-08a5-46bc-a1d8-21778634d4df	2026-06-17 16:51:35.290486+07	2026-06-17 16:51:35.290486+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	9df9e4c1-d939-4644-8cd2-e2669563f1e6	2026-06-17 16:51:35.296918+07	2026-06-17 16:51:35.296918+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	8e76a209-8787-4d22-8be5-8f076b4870aa	2026-06-17 16:51:35.300121+07	2026-06-17 16:51:35.300121+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	3b359e9a-b511-4812-9170-cc040396c452	2026-06-17 16:51:35.303399+07	2026-06-17 16:51:35.303399+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	6835b423-7e98-4fec-8522-7e7f175634bb	2026-06-17 16:51:35.307285+07	2026-06-17 16:51:35.307285+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	fb1db2c3-923d-44db-8fcb-6a033145ce51	2026-06-18 12:27:39.914584+07	2026-06-18 12:27:39.914584+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	d193356a-c252-43e4-92f4-25d3e243e985	2026-06-18 12:27:39.925027+07	2026-06-18 12:27:39.925027+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	9d65590f-d018-418d-a84a-be519fabd846	2026-06-18 12:27:39.928595+07	2026-06-18 12:27:39.928595+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	bac970ee-0871-472d-820a-566f322f1ea6	2026-06-18 12:27:39.944977+07	2026-06-18 12:27:39.944977+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	381b940e-91de-4985-9d80-b980bf67f855	2026-06-18 12:27:39.967239+07	2026-06-18 12:27:39.967239+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	c18d54de-ac69-4c7b-a2a4-1720efda3018	2026-06-18 14:18:12.992502+07	2026-06-18 14:18:12.992502+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	ff1cbd9c-6eab-43d8-8417-7c223d126ec6	2026-06-18 14:18:12.996368+07	2026-06-18 14:18:12.996368+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	9d4fdd3a-1e54-4340-b1ed-3d06fa54c678	2026-06-18 14:18:12.999869+07	2026-06-18 14:18:12.999869+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	0c70a9cd-a38f-4f5a-a73b-a45adc5113bf	2026-06-18 14:18:13.001991+07	2026-06-18 14:18:13.001991+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	8fdec684-aa37-49f2-b893-5698e0f747cb	2026-06-18 14:18:13.004504+07	2026-06-18 14:18:13.004504+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	c34019cd-fdc5-4680-825f-563dc9ad3fab	2026-06-18 14:23:53.820349+07	2026-06-18 14:23:53.820349+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	746774b9-b8c5-4d52-9448-f3aa296dd3a5	2026-06-18 14:23:53.824009+07	2026-06-18 14:23:53.824009+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	10a81054-462e-4007-adcd-7d186f8cc162	2026-06-18 14:23:53.827581+07	2026-06-18 14:23:53.827581+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	ce5c659d-270c-403e-9288-17b2ae6ca3e9	2026-06-18 14:23:53.830451+07	2026-06-18 14:23:53.830451+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	54a597d1-e407-439b-98ab-3c4cf49c6410	2026-06-18 14:23:53.833669+07	2026-06-18 14:23:53.833669+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	d38d2535-5cfd-4304-967b-a541ea841834	2026-06-18 22:49:28.999796+07	2026-06-18 22:49:28.999796+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	876de36e-0cf7-438e-8785-f4bf609c6662	2026-06-18 22:49:29.007679+07	2026-06-18 22:49:29.007679+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	a03a83a2-d4a2-44bf-96e0-6f303dd9be18	2026-06-18 22:49:29.01061+07	2026-06-18 22:49:29.01061+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	d7753601-45ea-4962-8c70-31815b9283d1	2026-06-18 22:49:29.014855+07	2026-06-18 22:49:29.014855+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	9d2b0f6d-b63d-49ee-9fce-fae4a0b87f5d	2026-06-18 22:49:29.017665+07	2026-06-18 22:49:29.017665+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	db629eb3-04b5-4bde-948f-7635bc662971	2026-06-18 23:12:57.2874+07	2026-06-18 23:12:57.2874+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	c722ad90-3fef-4d5e-9a6d-36e53400bcc8	2026-06-18 23:12:57.291416+07	2026-06-18 23:12:57.291416+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	f7b33bdb-2f23-4e30-bd60-906dfc220904	2026-06-18 23:12:57.297375+07	2026-06-18 23:12:57.297375+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	8d3119b1-0942-4923-8b8f-02fc82601438	2026-06-18 23:12:57.304605+07	2026-06-18 23:12:57.304605+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	946ee6f1-58c3-4a78-9d22-78793a9f5164	2026-06-18 23:12:57.31017+07	2026-06-18 23:12:57.31017+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	c50e54a0-e719-4d8f-b671-9d446f353c62	2026-06-18 23:20:04.387001+07	2026-06-18 23:20:04.387001+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	5c947aa0-bfb8-47b9-85f0-fa053a67cd3d	2026-06-18 23:20:04.41424+07	2026-06-18 23:20:04.41424+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	aaf3f7db-60d3-40a9-a2d5-bd71fa169cd9	2026-06-18 23:20:04.427439+07	2026-06-18 23:20:04.427439+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	fa56cb99-4dc3-4b7b-9e0b-384c2627f3b1	2026-06-18 23:20:04.43358+07	2026-06-18 23:20:04.43358+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	738022e4-86f8-461b-9797-cc3640c50465	2026-06-18 23:20:04.465173+07	2026-06-18 23:20:04.465173+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	8bfdcb9c-ed2c-4f88-af6e-4732041b4ed2	2026-06-19 09:15:16.090173+07	2026-06-19 09:15:16.090173+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	d370dc5c-c2ba-4a72-96f5-f3a9a585c114	2026-06-19 09:15:16.119425+07	2026-06-19 09:15:16.119425+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	83f3401e-34c4-4b8c-9bc8-4ab5e010ef38	2026-06-19 09:15:16.132299+07	2026-06-19 09:15:16.132299+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	dd958a94-fddf-43bb-93d9-cdf126afe1dc	2026-06-19 09:15:16.185151+07	2026-06-19 09:15:16.185151+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	e2ef3233-b90f-4387-a9a8-09e99ed9734f	2026-06-19 09:15:16.210219+07	2026-06-19 09:15:16.210219+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	43d827b2-456f-4bd8-98c5-365fad9999f5	2026-06-19 10:19:33.38769+07	2026-06-19 10:19:33.38769+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	347be4cb-30ad-470d-8ab1-b714af88062b	2026-06-19 10:19:33.44804+07	2026-06-19 10:19:33.44804+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	3e9bbec3-9e53-4e96-a0f9-92026b8bf14b	2026-06-19 10:19:33.451811+07	2026-06-19 10:19:33.451811+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	0bff78db-32c7-4a91-a6b8-03a5b86959b1	2026-06-19 10:19:33.459936+07	2026-06-19 10:19:33.459936+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	946b3bca-8618-43d2-a439-7887886b6286	2026-06-19 10:19:33.466208+07	2026-06-19 10:19:33.466208+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	d8dc04e4-2035-4570-8057-894010a50536	2026-06-19 11:24:05.424801+07	2026-06-19 11:24:05.424801+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	334a0931-d8bd-4343-b045-7f464f8220c4	2026-06-19 11:24:05.445835+07	2026-06-19 11:24:05.445835+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	3b779608-0919-4501-9283-64dcae28150f	2026-06-19 11:24:05.452502+07	2026-06-19 11:24:05.452502+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	1b1f47e4-caf3-4e28-a6af-7ba1cd4d401b	2026-06-19 11:24:05.473404+07	2026-06-19 11:24:05.473404+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	b5a4086a-f5e0-4566-bd22-0f7e19fb53af	2026-06-19 11:24:05.48966+07	2026-06-19 11:24:05.48966+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	aacd2dd6-017c-4084-800e-81e7ed6c4e10	2026-06-19 12:00:40.309741+07	2026-06-19 12:00:40.309741+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	3cd70bd0-8a27-4a3b-8b8f-5239da675e82	2026-06-19 12:00:40.359562+07	2026-06-19 12:00:40.359562+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	0a6860c6-38de-4e55-821e-82c89ace455f	2026-06-19 12:00:40.372173+07	2026-06-19 12:00:40.372173+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	c3bebcf4-dcc4-4736-8555-d8d33af800dd	2026-06-19 12:00:40.374851+07	2026-06-19 12:00:40.374851+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	d8210a42-4ed3-4d24-b5cf-1017b91992d7	2026-06-19 12:00:40.379849+07	2026-06-19 12:00:40.379849+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	a856edc1-4e6c-494d-8dd1-e97a251383a4	2026-06-19 12:51:13.829706+07	2026-06-19 12:51:13.829706+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	42717398-67bc-4bec-8005-109134772731	2026-06-19 12:51:13.837216+07	2026-06-19 12:51:13.837216+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	a5091565-d2d1-4e4e-a11c-4f0cbc5294a6	2026-06-19 12:51:13.840899+07	2026-06-19 12:51:13.840899+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	1ef62079-bb75-4c3f-9872-b8f2882dca43	2026-06-19 12:51:13.845867+07	2026-06-19 12:51:13.845867+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	f53c566c-6b83-4c4c-b92d-6b4bc524b1c6	2026-06-19 12:51:13.848286+07	2026-06-19 12:51:13.848286+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	d6ee4ec9-ccf7-4f5e-93bf-68a9bc40ee00	2026-06-19 13:57:44.665003+07	2026-06-19 13:57:44.665003+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	277eec3e-f754-4b13-83fd-4813966e3a81	2026-06-19 13:57:44.671101+07	2026-06-19 13:57:44.671101+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	777221fa-9723-41c4-968e-3fbd23b17751	2026-06-19 13:57:44.675861+07	2026-06-19 13:57:44.675861+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	c3055bd1-137b-4fd0-b527-f618a1541426	2026-06-19 13:57:44.681874+07	2026-06-19 13:57:44.681874+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	557b5dbb-7090-43a8-aedd-6b92ea163866	2026-06-19 13:57:44.687195+07	2026-06-19 13:57:44.687195+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	367ce449-a6c8-4d88-93ad-6de0c2f6f539	2026-06-19 14:04:59.565999+07	2026-06-19 14:04:59.565999+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	6f709d25-b0d3-4aa3-8f50-b3a125dc8b10	2026-06-19 14:04:59.600763+07	2026-06-19 14:04:59.600763+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	46571302-2f7e-4cdb-ab49-907a0193abc4	2026-06-19 14:04:59.646204+07	2026-06-19 14:04:59.646204+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	001fd4b2-c2d1-4abe-8f9c-113d3aa6f64a	2026-06-19 14:04:59.656765+07	2026-06-19 14:04:59.656765+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	b2975076-7b8b-4c2f-9e26-d4fd90ef7c7a	2026-06-19 14:04:59.682904+07	2026-06-19 14:04:59.682904+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	50c538c9-5c6c-4e8c-83a8-ef64b655137a	2026-06-19 14:11:46.20655+07	2026-06-19 14:11:46.20655+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	472410b6-620a-423f-b48c-d6c501c07116	2026-06-19 14:11:46.611568+07	2026-06-19 14:11:46.611568+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	7091025d-42b7-4dfd-b991-567954ae819c	2026-06-19 14:11:46.643321+07	2026-06-19 14:11:46.643321+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	dc13affe-4bc2-40aa-88bd-8a3e35796eb3	2026-06-19 14:11:46.676462+07	2026-06-19 14:11:46.676462+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	ef39d83d-8367-4089-9af0-5c67f15fcceb	2026-06-19 14:11:46.690583+07	2026-06-19 14:11:46.690583+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	64f7e2cf-af60-46ed-80bd-008a29876d79	2026-06-19 14:46:20.947394+07	2026-06-19 14:46:20.947394+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	efaa9032-4c4e-4289-a3db-a1b35c2ff2e9	2026-06-19 14:46:21.092019+07	2026-06-19 14:46:21.092019+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	b4551c48-4b37-4860-bddf-17def33ec627	2026-06-19 14:46:21.124857+07	2026-06-19 14:46:21.124857+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	157d995d-50e5-42b1-b254-0b1e0c7fc47f	2026-06-19 14:46:21.1329+07	2026-06-19 14:46:21.1329+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	2d60b511-3331-487e-b28b-a4f5e032e14e	2026-06-19 14:46:21.141585+07	2026-06-19 14:46:21.141585+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	f40dce8f-f69b-41b4-914f-d21ee6072fa5	2026-06-19 14:50:06.45137+07	2026-06-19 14:50:06.45137+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	f1e682ab-9880-4cb5-a43d-8995f14804ac	2026-06-19 14:50:06.455352+07	2026-06-19 14:50:06.455352+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	a2a9d80e-6c23-4552-8f62-5915e431a814	2026-06-19 14:50:06.463775+07	2026-06-19 14:50:06.463775+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	e3f53df9-b343-4846-a749-cad4819f4b87	2026-06-19 14:50:06.468818+07	2026-06-19 14:50:06.468818+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	d7b37c1c-e491-4260-8813-8e82d55e4210	2026-06-19 14:50:06.473293+07	2026-06-19 14:50:06.473293+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	72b7f73b-75eb-4207-a923-92a946ad807e	2026-06-19 14:59:26.506236+07	2026-06-19 14:59:26.506236+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	7b8f2f6a-cdc6-4b2d-af60-c04a0949f7c8	2026-06-19 14:59:26.534994+07	2026-06-19 14:59:26.534994+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	f703d913-44c1-4817-91ae-a21a7ed01e53	2026-06-19 14:59:26.54884+07	2026-06-19 14:59:26.54884+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	9f709557-d2ba-469e-8aa9-8afa48612bd8	2026-06-19 14:59:26.557438+07	2026-06-19 14:59:26.557438+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	84ba3538-b09a-498a-8de8-79683bd8b630	2026-06-19 14:59:26.575227+07	2026-06-19 14:59:26.575227+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	f8803534-0c44-4ffc-b00a-9eb6378878ab	2026-06-19 19:15:51.758189+07	2026-06-19 19:15:51.758189+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	21e6385f-314b-4c94-801b-c224cc93d6bb	2026-06-19 19:15:51.767642+07	2026-06-19 19:15:51.767642+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	06a2e66b-cc3f-429b-ab7c-dce551278cbe	2026-06-19 19:15:51.772108+07	2026-06-19 19:15:51.772108+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	e76a580c-8eb0-4a8d-ac6f-7883dff2ebca	2026-06-19 19:15:51.779491+07	2026-06-19 19:15:51.779491+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	c9d3895c-63b7-48ca-846b-bebd2a9b903a	2026-06-19 19:15:51.785468+07	2026-06-19 19:15:51.785468+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	bde0570d-6a47-445a-b462-4982a0478e42	2026-06-19 23:23:47.281733+07	2026-06-19 23:23:47.281733+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	3da14516-1879-4b86-bd6c-020d42a96128	2026-06-19 23:23:47.284827+07	2026-06-19 23:23:47.284827+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	fbd06939-f3d0-4ccc-9662-9e4ff9ed677d	2026-06-19 23:23:47.287975+07	2026-06-19 23:23:47.287975+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	b817624a-d55f-42f7-a96d-05e8beaad403	2026-06-19 23:23:47.290966+07	2026-06-19 23:23:47.290966+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	089c81f0-1061-47c1-adf2-90a24ea7c82c	2026-06-19 23:23:47.293944+07	2026-06-19 23:23:47.293944+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	e659d5c2-9cdf-442d-80ba-e1cfdc82b44c	2026-06-20 15:58:00.702696+07	2026-06-20 15:58:00.702696+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	442b9045-fe2e-41de-96ab-d04222c0a51e	2026-06-20 15:58:00.716611+07	2026-06-20 15:58:00.716611+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	e29f0f34-d182-4551-9a2a-172fbed1ada8	2026-06-20 15:58:00.726535+07	2026-06-20 15:58:00.726535+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	283b9eae-461a-4406-a2d7-fb84a006ef40	2026-06-20 15:58:00.733417+07	2026-06-20 15:58:00.733417+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	f869e162-f0c3-44e2-a662-a5a3a5bb0750	2026-06-20 15:58:00.742026+07	2026-06-20 15:58:00.742026+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	a7e75a8b-64d8-49ce-8410-1859b6fd441c	2026-06-21 00:12:22.391329+07	2026-06-21 00:12:22.391329+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	a9aba572-8493-4a4a-bf57-1fd813f194b3	2026-06-21 00:12:22.397502+07	2026-06-21 00:12:22.397502+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	72a86900-ded9-4bed-b760-dc3c4eef1b26	2026-06-21 00:12:22.407407+07	2026-06-21 00:12:22.407407+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	25ece34a-d6b6-44b0-abdd-8ffabd7cbcf7	2026-06-21 00:12:22.415865+07	2026-06-21 00:12:22.415865+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	22b99b87-26cb-4c51-9951-4cd066531f73	2026-06-21 00:12:22.424653+07	2026-06-21 00:12:22.424653+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	ee1c1f3f-c097-43de-a31a-e84f42abe7a6	2026-06-21 18:12:28.279995+07	2026-06-21 18:12:28.279995+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	007751ce-2c99-4c77-8aeb-895a5e604113	2026-06-21 18:12:28.289163+07	2026-06-21 18:12:28.289163+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	edb63213-fcc2-47a8-a68d-c3ad5f3e41db	2026-06-21 18:12:28.292674+07	2026-06-21 18:12:28.292674+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	69ffda1d-33c7-4b4a-8add-046952ff25eb	2026-06-21 18:12:28.296044+07	2026-06-21 18:12:28.296044+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	a339c93b-7233-4a9d-80e3-c9fa4ec3d1c2	2026-06-21 18:12:28.299448+07	2026-06-21 18:12:28.299448+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	f9633705-b1af-46bf-a844-5d5fa7b3f4e9	2026-06-22 09:27:11.934851+07	2026-06-22 09:27:11.934851+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	4bdd1882-4468-4845-9e13-36974ef2042b	2026-06-22 09:27:11.943+07	2026-06-22 09:27:11.943+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	66a5b63d-c0d5-46e0-9c41-1fa4bb8afa59	2026-06-22 09:27:11.95004+07	2026-06-22 09:27:11.95004+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	ca64c8fa-c933-4455-b611-d6bdc779acd8	2026-06-22 09:27:11.953676+07	2026-06-22 09:27:11.953676+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	618f6c68-f221-4219-8d36-8728216337ec	2026-06-22 09:27:11.962011+07	2026-06-22 09:27:11.962011+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	aa042413-7547-4ffa-915d-cc15ba7a7c2d	2026-06-22 16:06:12.199+07	2026-06-22 16:06:12.199+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	95c6cc09-3a2d-48aa-b6c1-90dc05e5c9b8	2026-06-22 16:06:12.20501+07	2026-06-22 16:06:12.20501+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	60782b30-3d42-4df7-b049-273413054ef7	2026-06-22 16:06:12.210184+07	2026-06-22 16:06:12.210184+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	896a69dc-961d-4475-8231-9ff28d94d1fe	2026-06-22 16:06:12.215267+07	2026-06-22 16:06:12.215267+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	d642292f-5d15-4b7c-96ce-5f17dcac4dcb	2026-06-22 16:06:12.22196+07	2026-06-22 16:06:12.22196+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	d0aa50b3-1fc1-4d42-97d9-8d51dc7c0ab3	2026-06-23 14:46:32.422662+07	2026-06-23 14:46:32.422662+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	9b6b4860-d132-4846-b931-c0fae249308b	2026-06-23 14:46:32.512693+07	2026-06-23 14:46:32.512693+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	a2139035-7f1d-4299-ad60-b27ac1498482	2026-06-23 14:46:32.520484+07	2026-06-23 14:46:32.520484+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	d69977b8-ef46-4b59-9f9a-7f3d625b2e88	2026-06-23 14:46:32.529829+07	2026-06-23 14:46:32.529829+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	fa78d002-ab2a-4538-962e-2e60946aaa5f	2026-06-23 14:46:32.553172+07	2026-06-23 14:46:32.553172+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	e91840ad-f90f-4467-8157-d5443b5fa272	2026-06-23 15:48:43.927493+07	2026-06-23 15:48:43.927493+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	bc963014-52cb-419b-b7cc-0393c6dfe82c	2026-06-23 15:48:43.93494+07	2026-06-23 15:48:43.93494+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	183f3998-f028-4e6e-bcb0-fb5dd8f5b5a2	2026-06-23 15:48:43.944346+07	2026-06-23 15:48:43.944346+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	4e615b08-48d5-4e21-919a-8589dba0171c	2026-06-23 15:48:43.952086+07	2026-06-23 15:48:43.952086+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	5406293d-9f06-4058-9878-0301dc53c094	2026-06-23 15:48:43.959849+07	2026-06-23 15:48:43.959849+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	eddbe4a3-ad29-4802-bf53-321943896638	2026-06-23 15:52:06.724514+07	2026-06-23 15:52:06.724514+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	a16f490d-53fa-4e7a-b560-fb43751d4d80	2026-06-23 15:52:06.727681+07	2026-06-23 15:52:06.727681+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	90303a95-615b-448f-aaa5-c9401ae89d81	2026-06-23 15:52:06.730204+07	2026-06-23 15:52:06.730204+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	f004b6d4-574f-48a9-8458-a32bd9e21fb0	2026-06-23 15:52:06.732347+07	2026-06-23 15:52:06.732347+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	e9bfb6b2-7330-435b-90d9-9e9435fc9b07	2026-06-23 15:52:06.735512+07	2026-06-23 15:52:06.735512+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	890a5396-7696-4519-9915-2ef9755d5238	2026-06-23 16:16:22.012638+07	2026-06-23 16:16:22.012638+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	d95ba07f-4ee3-4d20-9296-de089923f192	2026-06-23 16:16:22.018856+07	2026-06-23 16:16:22.018856+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	93b3740d-20d2-4b19-9cb7-5a6ca3d6d089	2026-06-23 16:16:22.029832+07	2026-06-23 16:16:22.029832+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	af291ea7-8a4c-44b3-8c14-5f08cede7804	2026-06-23 16:16:22.039108+07	2026-06-23 16:16:22.039108+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	7d9278ec-0598-40ac-b040-58fd3bba26f6	2026-06-23 16:16:22.046445+07	2026-06-23 16:16:22.046445+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	ee3de025-cde0-4d95-9f84-00c5e115f8db	2026-06-23 18:19:32.22419+07	2026-06-23 18:19:32.22419+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	62d4333d-6579-4da4-9a70-413c43abd325	2026-06-23 18:19:32.234177+07	2026-06-23 18:19:32.234177+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	7d1811f3-664e-4258-bdb7-6d8566d22439	2026-06-23 18:19:32.240702+07	2026-06-23 18:19:32.240702+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	0fa23c2f-be57-4802-a594-30ec397ebb7b	2026-06-23 18:19:32.243912+07	2026-06-23 18:19:32.243912+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	4f5339ad-98cc-4ea9-8053-418ca51cb650	2026-06-23 18:19:32.255291+07	2026-06-23 18:19:32.255291+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	bace635b-2597-4ad7-8b44-b667b899d715	2026-06-23 18:26:10.670057+07	2026-06-23 18:26:10.670057+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	5b6c8e44-1953-4903-b2f0-a036d247d3f1	2026-06-23 18:26:10.674395+07	2026-06-23 18:26:10.674395+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	4ff90501-c538-4514-8912-a70761ec6f2b	2026-06-23 18:26:10.678134+07	2026-06-23 18:26:10.678134+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	7eacba73-3bfb-4b8a-b35d-823646453e69	2026-06-23 18:26:10.681201+07	2026-06-23 18:26:10.681201+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	92889257-05ee-4dca-8e72-fb5d7a5d5f05	2026-06-23 18:26:10.682947+07	2026-06-23 18:26:10.682947+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	18b4376b-2208-4982-b427-fff4b7d87a81	2026-06-23 21:50:01.271696+07	2026-06-23 21:50:01.271696+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	f9212361-ee83-42ef-ab54-121db89e8b13	2026-06-23 21:50:01.275977+07	2026-06-23 21:50:01.275977+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	7fe62d9d-3d5b-499d-bd53-c10617606b7e	2026-06-23 21:50:01.280209+07	2026-06-23 21:50:01.280209+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	4b1ad9bc-1156-4a57-9fc9-75f424e433ce	2026-06-23 21:50:01.284614+07	2026-06-23 21:50:01.284614+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	b65cd1db-104c-427f-9434-53b32c1e0ed3	2026-06-23 21:50:01.287554+07	2026-06-23 21:50:01.287554+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	fba662f7-af2b-444f-9f83-4a23fa33a09e	2026-06-24 07:11:38.330922+07	2026-06-24 07:11:38.330922+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	80acaeba-c836-417a-98c1-907c4e18bdec	2026-06-24 07:11:38.336633+07	2026-06-24 07:11:38.336633+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	669d7647-17b1-4b7f-b9c2-1358345d0eba	2026-06-24 07:11:38.348492+07	2026-06-24 07:11:38.348492+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	ab8ea925-4dbb-4f33-a817-768c6f61879f	2026-06-24 07:11:38.353075+07	2026-06-24 07:11:38.353075+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	bb52fda4-e797-4a3c-b543-d11323c99c4b	2026-06-24 07:11:38.357879+07	2026-06-24 07:11:38.357879+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	02a440cb-e906-487f-a35c-b72938c2ec76	2026-06-24 17:01:03.706321+07	2026-06-24 17:01:03.706321+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	c0b190fb-ccd6-463f-a94e-3143b6d301c1	2026-06-24 17:01:03.726595+07	2026-06-24 17:01:03.726595+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	cc15656c-7ed9-4fce-9084-3d7d8f3dabd6	2026-06-24 17:01:03.796588+07	2026-06-24 17:01:03.796588+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	31529a73-7ab9-4a08-9786-fc574784d7df	2026-06-24 17:01:03.801054+07	2026-06-24 17:01:03.801054+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	e35450ba-fff7-4299-a991-90e6b03286f1	2026-06-24 17:01:03.808419+07	2026-06-24 17:01:03.808419+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	f4dcac0c-d644-4a30-87c1-9f41c08bd26b	2026-06-25 00:19:34.866712+07	2026-06-25 00:19:34.866712+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	8e5eed7a-3707-4554-9c6d-31121b374e8b	2026-06-25 00:19:34.872183+07	2026-06-25 00:19:34.872183+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	fd48f864-fe9b-4055-bc1e-5f5ed43c03ea	2026-06-25 00:19:34.875014+07	2026-06-25 00:19:34.875014+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	0ea45a09-9e9b-4013-9a54-ce7f6ab30cad	2026-06-25 00:19:34.878121+07	2026-06-25 00:19:34.878121+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	d65f3d6e-2cae-4649-b31f-86d466f456d8	2026-06-25 00:19:34.883362+07	2026-06-25 00:19:34.883362+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	a9373d14-f7f1-425a-baa8-20a494096ff1	2026-06-25 00:22:12.438745+07	2026-06-25 00:22:12.438745+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	1290a471-0d59-4de0-a3a5-754c2e0ac81f	2026-06-25 00:22:12.44287+07	2026-06-25 00:22:12.44287+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	0e89b5aa-a9f9-4036-8c28-85946c5201ce	2026-06-25 00:22:12.448968+07	2026-06-25 00:22:12.448968+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	711b327f-a46c-4ace-b963-493e21418707	2026-06-25 00:22:12.454118+07	2026-06-25 00:22:12.454118+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	df30641d-a593-4087-a306-a578589edfc0	2026-06-25 00:22:12.457547+07	2026-06-25 00:22:12.457547+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	8e9045ca-2cc0-4093-bc0b-da31052e66c2	2026-06-25 01:49:07.802458+07	2026-06-25 01:49:07.802458+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	46435f5f-57e1-4fb9-8e1d-1bdbc4896017	2026-06-25 01:49:07.805917+07	2026-06-25 01:49:07.805917+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	6c2b93a6-7fc1-4d1f-a61d-067fe52818d5	2026-06-25 01:49:07.810673+07	2026-06-25 01:49:07.810673+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	83cb93e2-17f0-4137-baa9-7a2e919222c9	2026-06-25 01:49:07.815434+07	2026-06-25 01:49:07.815434+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	21dcbfe8-49b1-48b3-a51b-d7b601ef3a07	2026-06-25 01:49:07.820018+07	2026-06-25 01:49:07.820018+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	228b7603-6179-42a3-87ea-fab89c195c9e	2026-06-25 01:58:46.090298+07	2026-06-25 01:58:46.090298+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	0023be66-0682-434d-b504-5eb2c3eee771	2026-06-25 01:58:46.099958+07	2026-06-25 01:58:46.099958+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	e578d6e3-b9d3-40b7-a0dc-52f898eff76e	2026-06-25 01:58:46.103294+07	2026-06-25 01:58:46.103294+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	03b3d41e-8f35-4957-9493-43817d1a9fc0	2026-06-25 01:58:46.108689+07	2026-06-25 01:58:46.108689+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	33dd74d8-728b-4674-be7d-c19f58504330	2026-06-25 01:58:46.116449+07	2026-06-25 01:58:46.116449+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	c5302cf1-eec4-46e4-b517-8275c30d5527	2026-06-25 02:14:48.407042+07	2026-06-25 02:14:48.407042+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	354832ce-88a7-4726-8d9a-0a95c3cd4f12	2026-06-25 02:14:48.414801+07	2026-06-25 02:14:48.414801+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	4038f07b-e277-4848-a88f-3c63e8d9fdd1	2026-06-25 02:14:48.419691+07	2026-06-25 02:14:48.419691+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	e672a4d5-5935-44b5-9697-ea0d43d9d03a	2026-06-25 02:14:48.423396+07	2026-06-25 02:14:48.423396+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	0f116082-40d2-471e-88a9-ddbdc5f5a049	2026-06-25 02:14:48.426873+07	2026-06-25 02:14:48.426873+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	155dacc1-630f-4b5d-a782-b8a0f3ea588c	2026-06-25 02:29:13.278272+07	2026-06-25 02:29:13.278272+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	16df664c-878c-4fc4-bec7-8714707ea6cd	2026-06-25 02:29:13.283448+07	2026-06-25 02:29:13.283448+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	a1d7484e-7ab2-47f5-a36a-c103c2668063	2026-06-25 02:29:13.288098+07	2026-06-25 02:29:13.288098+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	afb0cecd-416a-4492-9758-89931e9aabc0	2026-06-25 02:29:13.291962+07	2026-06-25 02:29:13.291962+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	2f622454-08d1-4f22-984f-a3dab35fb3d6	2026-06-25 02:29:13.294984+07	2026-06-25 02:29:13.294984+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	78876cc5-21d8-4ef6-84ba-76d2fda0bc29	2026-06-25 02:38:08.638689+07	2026-06-25 02:38:08.638689+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	38827f72-9b20-4863-951f-ac6183034896	2026-06-25 02:38:08.661637+07	2026-06-25 02:38:08.661637+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	f43980c1-9497-4c07-a383-0868ce7f7a68	2026-06-25 02:38:08.684908+07	2026-06-25 02:38:08.684908+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	5e3da3b6-0ffe-45c2-9678-9b42ff820212	2026-06-25 02:38:08.723657+07	2026-06-25 02:38:08.723657+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	1f6ef26d-8c62-4df1-92f5-284ee30a3d2b	2026-06-25 02:38:08.730641+07	2026-06-25 02:38:08.730641+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	b9e8a4e6-949b-4ba1-95d2-e825a5ca3a72	2026-06-25 09:16:31.050667+07	2026-06-25 09:16:31.050667+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	62dda25a-4cd8-49cc-ae21-c4727bdaac6c	2026-06-25 09:16:31.080404+07	2026-06-25 09:16:31.080404+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	16d410c1-eba0-4615-bb10-6548d7ef7598	2026-06-25 09:16:31.085852+07	2026-06-25 09:16:31.085852+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	0641a8a9-7a5f-40c2-84b7-ad89258800ff	2026-06-25 09:16:31.147362+07	2026-06-25 09:16:31.147362+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	cd07463c-9285-4c76-9d48-e71c8820d74c	2026-06-25 09:16:31.160731+07	2026-06-25 09:16:31.160731+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	a00d0574-4da8-41f9-a3fc-db0839440c29	2026-06-25 10:21:42.830157+07	2026-06-25 10:21:42.830157+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	55196e08-bed6-4158-a8bd-53bd738c25a7	2026-06-25 10:21:42.83586+07	2026-06-25 10:21:42.83586+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	98812da2-7242-4fe4-9bf0-0f2bac2036fc	2026-06-25 10:21:42.842445+07	2026-06-25 10:21:42.842445+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	5ae1b6cc-776a-4aa9-b842-a04604673bfe	2026-06-25 10:21:42.846852+07	2026-06-25 10:21:42.846852+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	12af450f-eeed-42da-9454-66a8b49ae5d1	2026-06-25 10:21:42.851336+07	2026-06-25 10:21:42.851336+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	5ff71fa5-eded-4c66-bb43-b31841b7b8de	2026-06-25 11:18:36.647352+07	2026-06-25 11:18:36.647352+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	9b67679f-87d0-4b89-acc4-b2d35dfd0b45	2026-06-25 11:18:36.657195+07	2026-06-25 11:18:36.657195+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	54f449e8-d981-4163-b29d-5c47825d9d59	2026-06-25 11:18:36.662999+07	2026-06-25 11:18:36.662999+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	e0782603-f8b1-48de-b137-9635145efdc3	2026-06-25 11:18:36.66846+07	2026-06-25 11:18:36.66846+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	06754414-468d-4aec-b645-05c34bb7240c	2026-06-25 11:18:36.673569+07	2026-06-25 11:18:36.673569+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	52a0fd8e-1e92-4b4b-af61-8f8658b052b3	2026-06-25 11:20:22.840943+07	2026-06-25 11:20:22.840943+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	2fd9e129-d7ed-4100-9494-21418001c9ca	2026-06-25 11:20:22.85178+07	2026-06-25 11:20:22.85178+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	df150a79-65c4-477e-a8c9-0694c3332146	2026-06-25 11:20:22.855624+07	2026-06-25 11:20:22.855624+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	2cdb5a2b-df70-4d30-a5b5-321f24adf992	2026-06-25 11:20:22.86179+07	2026-06-25 11:20:22.86179+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	c1113f7f-0ac9-471e-a358-0289cd774c0c	2026-06-25 11:20:22.866568+07	2026-06-25 11:20:22.866568+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	a8e36749-cd8d-466f-a53c-0679aaefc5e5	2026-06-25 11:58:32.410516+07	2026-06-25 11:58:32.410516+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	267d320a-bab4-4d40-843d-dbaf76c89b75	2026-06-25 11:58:32.420211+07	2026-06-25 11:58:32.420211+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	4549d290-0cd0-486e-b609-1683d1fa1a50	2026-06-25 11:58:32.424628+07	2026-06-25 11:58:32.424628+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	912d013a-438b-4000-b54f-87cb4b4804b6	2026-06-25 11:58:32.430485+07	2026-06-25 11:58:32.430485+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	01666565-054d-43aa-83b8-94ef6f41489e	2026-06-25 11:58:32.437176+07	2026-06-25 11:58:32.437176+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	a345e725-2061-45f4-b49b-dd349b2cf0ef	2026-06-25 12:33:41.711596+07	2026-06-25 12:33:41.711596+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	796e4bd5-8d29-461a-9c93-248ff9283699	2026-06-25 12:33:41.716281+07	2026-06-25 12:33:41.716281+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	5b13d9c7-b718-47ef-8e5c-a76b7f664925	2026-06-25 12:33:41.724277+07	2026-06-25 12:33:41.724277+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	cfe524d0-ad56-4608-8cd6-ed58e8c42ac7	2026-06-25 12:33:41.727283+07	2026-06-25 12:33:41.727283+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	bab6bc7a-9802-431f-bf35-bf2620b7c79e	2026-06-25 12:33:41.730691+07	2026-06-25 12:33:41.730691+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	1a2dd832-fd8e-4293-b9c1-da2c40770880	2026-06-26 09:39:22.546872+07	2026-06-26 09:39:22.546872+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	44be86f9-89e2-4977-ae65-776f7eb7143a	2026-06-26 09:39:22.554907+07	2026-06-26 09:39:22.554907+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	b7ce092d-5b25-4cbc-8788-8ebd1cde3d62	2026-06-26 09:39:22.562099+07	2026-06-26 09:39:22.562099+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	dcb25c5b-7a38-4646-bb07-43b56a8e6e3c	2026-06-26 09:39:22.567162+07	2026-06-26 09:39:22.567162+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	665d2569-f1db-47d6-979e-1978ad153e89	2026-06-26 09:39:22.574714+07	2026-06-26 09:39:22.574714+07
a24c7797-3d0a-4190-8c14-bac2bd870c09	district	108	2f0bee12-9531-4685-bec2-510c3b4a7354	2026-06-26 09:39:22.579558+07	2026-06-26 09:39:22.579558+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	1ae782ad-7aa6-4133-8821-30df37c2da8c	2026-06-26 12:15:35.450539+07	2026-06-26 12:15:35.450539+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	d7ec65d9-77a1-4ea2-ae35-c456674d2a11	2026-06-26 12:15:35.457384+07	2026-06-26 12:15:35.457384+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	19938860-6a06-4cc0-9ff2-129dfa6032a1	2026-06-26 12:15:35.461477+07	2026-06-26 12:15:35.461477+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	a8f025d7-ed4c-44b0-83d2-3bb6e742ffad	2026-06-26 12:15:35.465641+07	2026-06-26 12:15:35.465641+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	98318d82-4a83-4f24-b2b2-4fcd96040a39	2026-06-26 12:15:35.470598+07	2026-06-26 12:15:35.470598+07
a24c7797-3d0a-4190-8c14-bac2bd870c09	district	108	a8a7c797-4a40-412f-9720-5a213b3cde82	2026-06-26 12:15:35.473908+07	2026-06-26 12:15:35.473908+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	accd5ad7-1575-49d1-8b85-8c75f81d55cf	2026-06-26 12:35:26.871588+07	2026-06-26 12:35:26.871588+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	61b2fc46-36e3-4b07-b288-a2efa9cc20cf	2026-06-26 12:35:26.894032+07	2026-06-26 12:35:26.894032+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	fb2ca544-2e79-4227-bc68-08683c5a552b	2026-06-26 12:35:26.900088+07	2026-06-26 12:35:26.900088+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	9213f075-ec4a-483c-9ea9-c092d5472601	2026-06-26 12:35:26.920417+07	2026-06-26 12:35:26.920417+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	3f976b41-0a4a-4e17-b704-cfd12c7e4424	2026-06-26 12:35:26.939691+07	2026-06-26 12:35:26.939691+07
a24c7797-3d0a-4190-8c14-bac2bd870c09	district	108	257dd96f-1802-4598-b885-b60376bc66cc	2026-06-26 12:35:26.949382+07	2026-06-26 12:35:26.949382+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	a8357d47-f1fd-4b6e-8b9a-19b107b90ac3	2026-06-26 16:40:30.022232+07	2026-06-26 16:40:30.022232+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	2694481e-9b54-473a-96b1-fc367a5fd370	2026-06-26 16:40:30.033323+07	2026-06-26 16:40:30.033323+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	dcc1d87a-b401-4ab9-b82c-cf4f98e7a536	2026-06-26 16:40:30.03591+07	2026-06-26 16:40:30.03591+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	29dabfde-f9d6-4a8c-ab7a-15553954bad3	2026-06-26 16:40:30.040048+07	2026-06-26 16:40:30.040048+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	d67131a0-ce09-4fa4-bd42-ea77bb07d4fd	2026-06-26 16:40:30.04685+07	2026-06-26 16:40:30.04685+07
a24c7797-3d0a-4190-8c14-bac2bd870c09	district	108	120de228-cc8d-41b5-bc3c-5b7f1b24a10b	2026-06-26 16:40:30.05067+07	2026-06-26 16:40:30.05067+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	fd8954ae-fc14-4356-bcc5-3bbadb776f64	2026-06-26 18:25:23.711901+07	2026-06-26 18:25:23.711901+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	86b48905-2dd3-49b9-8f1f-1a518768f079	2026-06-26 18:25:23.748446+07	2026-06-26 18:25:23.748446+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	622c9e8c-2d70-447d-b73e-fff33ed49d64	2026-06-26 18:25:23.752687+07	2026-06-26 18:25:23.752687+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	5586c6f6-779d-474f-b5cd-9f00c17065f6	2026-06-26 18:25:23.758694+07	2026-06-26 18:25:23.758694+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	69939984-a8f5-4059-b713-a31410f8f082	2026-06-26 18:25:23.778486+07	2026-06-26 18:25:23.778486+07
a24c7797-3d0a-4190-8c14-bac2bd870c09	district	108	99caed4b-1cda-4330-8909-3a5df21570df	2026-06-26 18:25:23.800319+07	2026-06-26 18:25:23.800319+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	20692187-a808-4948-88ae-2f6e9a732603	2026-06-26 20:05:24.856924+07	2026-06-26 20:05:24.856924+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	12b07524-e7f5-4bbb-86cf-4396d3d8cd87	2026-06-26 20:05:24.863163+07	2026-06-26 20:05:24.863163+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	b20aafdf-cbbb-47af-8722-29176d2b9b32	2026-06-26 20:05:24.867181+07	2026-06-26 20:05:24.867181+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	04a929e9-b01b-4088-afc2-84906b2a6fbd	2026-06-26 20:05:24.871457+07	2026-06-26 20:05:24.871457+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	8adcec83-69fb-41e8-98b2-25519aa24184	2026-06-26 20:05:24.874284+07	2026-06-26 20:05:24.874284+07
a24c7797-3d0a-4190-8c14-bac2bd870c09	district	108	58c111e4-e5a4-474c-85ed-71a11f62b6ef	2026-06-26 20:05:24.881717+07	2026-06-26 20:05:24.881717+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	08c79fdb-272d-432c-9ba8-b71e423c120b	2026-06-27 10:36:37.249248+07	2026-06-27 10:36:37.249248+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	d76e7341-bed9-4806-b19e-8bcc12573b9d	2026-06-27 10:36:37.255924+07	2026-06-27 10:36:37.255924+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	cc3fb834-3131-4db7-91ec-e53954cd38da	2026-06-27 10:36:37.259128+07	2026-06-27 10:36:37.259128+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	ce819cb5-4d4b-4d8b-a314-fc6c91b90cb6	2026-06-27 10:36:37.26313+07	2026-06-27 10:36:37.26313+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	aafc5a93-9866-491c-9d8e-b7ac2f2a3970	2026-06-27 10:36:37.265857+07	2026-06-27 10:36:37.265857+07
a24c7797-3d0a-4190-8c14-bac2bd870c09	district	108	2a35122c-1735-4452-b4a1-7fb31705c986	2026-06-27 10:36:37.269517+07	2026-06-27 10:36:37.269517+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	1ce17164-a257-4805-b39b-9217625ed2a6	2026-06-27 11:36:26.958216+07	2026-06-27 11:36:26.958216+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	db6141d2-bb01-4839-8a66-4eb7ef48432a	2026-06-27 11:36:26.967828+07	2026-06-27 11:36:26.967828+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	dff29091-3032-41eb-93da-0c00265e0e84	2026-06-27 11:36:26.975724+07	2026-06-27 11:36:26.975724+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	b816db7f-06d5-4b11-84b8-9f0506c9881f	2026-06-27 11:36:26.98081+07	2026-06-27 11:36:26.98081+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	4a632e1a-3826-4ae9-ad1d-c69ef46633d5	2026-06-27 11:36:26.983826+07	2026-06-27 11:36:26.983826+07
a24c7797-3d0a-4190-8c14-bac2bd870c09	district	108	2859248e-d79f-4e8c-94be-66fad345f805	2026-06-27 11:36:26.988411+07	2026-06-27 11:36:26.988411+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	dcba70cb-c7cd-4bd8-9edd-bff16ecca4a9	2026-06-27 11:49:02.790646+07	2026-06-27 11:49:02.790646+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	791bda3e-03fb-40e1-91f6-0afc85501d71	2026-06-27 11:49:02.796148+07	2026-06-27 11:49:02.796148+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	c704172c-ec8e-454f-a590-670a40207099	2026-06-27 11:49:02.801497+07	2026-06-27 11:49:02.801497+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	4e265a50-8c98-4389-af64-f6d6ce9cdd22	2026-06-27 11:49:02.806096+07	2026-06-27 11:49:02.806096+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	b268f5c0-05a7-4330-9d66-24f9d5542810	2026-06-27 11:49:02.810732+07	2026-06-27 11:49:02.810732+07
a24c7797-3d0a-4190-8c14-bac2bd870c09	district	108	3177c3b0-4510-4e33-910f-a3f60e091a9f	2026-06-27 11:49:02.814761+07	2026-06-27 11:49:02.814761+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	1b5a5ea1-3cba-4c5d-bc01-633f86fc5baa	2026-06-27 13:37:36.202807+07	2026-06-27 13:37:36.202807+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	cc65878f-a68e-4ee3-b141-d5c20ce75fb0	2026-06-27 13:37:36.220016+07	2026-06-27 13:37:36.220016+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	c9cd5b64-2e71-4085-b993-9b166c855307	2026-06-27 13:37:36.228673+07	2026-06-27 13:37:36.228673+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	9c875640-2b23-481a-9abf-80b8e4322e6f	2026-06-27 13:37:36.232674+07	2026-06-27 13:37:36.232674+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	356335ad-be03-4f02-9c52-39b796a58b97	2026-06-27 13:37:36.235643+07	2026-06-27 13:37:36.235643+07
a24c7797-3d0a-4190-8c14-bac2bd870c09	district	108	d5ccd2f9-bdea-495e-9e80-cb6f3edaec1a	2026-06-27 13:37:36.23784+07	2026-06-27 13:37:36.23784+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	9890d8a7-f482-409e-9c49-bda652ff2380	2026-06-27 14:17:53.267903+07	2026-06-27 14:17:53.267903+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	c21b8979-e164-4c93-9b34-4e857ba296c5	2026-06-27 14:17:53.273411+07	2026-06-27 14:17:53.273411+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	4ed8099f-07d1-4491-8f1b-d0a273034c36	2026-06-27 14:17:53.278882+07	2026-06-27 14:17:53.278882+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	ac45048f-6f3e-4914-81b9-b34dfad39beb	2026-06-27 14:17:53.28448+07	2026-06-27 14:17:53.28448+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	452c181b-f72b-42f7-b6a1-a67e1f77e01f	2026-06-27 14:17:53.290713+07	2026-06-27 14:17:53.290713+07
a24c7797-3d0a-4190-8c14-bac2bd870c09	district	108	127e6b95-8f9c-4ddd-b3d2-ffed11a6347f	2026-06-27 14:17:53.295892+07	2026-06-27 14:17:53.295892+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	69e3172a-222d-40b1-b950-85f863528e64	2026-06-27 14:20:45.955548+07	2026-06-27 14:20:45.955548+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	449c3db6-f58e-45f6-8a8d-644017189e70	2026-06-27 14:20:45.95819+07	2026-06-27 14:20:45.95819+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	6ee1bf00-0460-4d89-933e-d7dc30bcbced	2026-06-27 14:20:45.961589+07	2026-06-27 14:20:45.961589+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	94e43987-5fd3-42f1-877e-e2d5db8ada01	2026-06-27 14:20:45.964805+07	2026-06-27 14:20:45.964805+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	9e053082-0397-455e-b09c-8e3e11ba1560	2026-06-27 14:20:45.968335+07	2026-06-27 14:20:45.968335+07
a24c7797-3d0a-4190-8c14-bac2bd870c09	district	108	84c6507b-1f01-4389-893a-ed28c2cc6604	2026-06-27 14:20:45.97083+07	2026-06-27 14:20:45.97083+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	9bc845af-eb93-4015-9958-7c23a88326ab	2026-06-27 18:21:54.060947+07	2026-06-27 18:21:54.060947+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	2b48b96d-d10d-4bd5-bdd5-3964371b13c4	2026-06-27 18:21:54.07415+07	2026-06-27 18:21:54.07415+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	13c19050-bfdf-4e48-b05d-e82398b01230	2026-06-27 18:21:54.085621+07	2026-06-27 18:21:54.085621+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	915ff54d-d945-4109-afd0-b97ab5df8dbb	2026-06-27 18:21:54.096171+07	2026-06-27 18:21:54.096171+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	344704e9-1f39-4fcf-bc52-236e1a715d1c	2026-06-27 18:21:54.107929+07	2026-06-27 18:21:54.107929+07
a24c7797-3d0a-4190-8c14-bac2bd870c09	district	108	32b4e56c-bb1b-43e5-ad2e-31a63f54d26c	2026-06-27 18:21:54.128249+07	2026-06-27 18:21:54.128249+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	6	1a4163bf-03ef-47df-8607-0c9fd415b4d5	2026-07-08 14:08:11.034933+07	2026-07-08 14:08:11.034933+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	7	e66676d7-7c95-40e9-881d-7f0a791f1c5e	2026-07-08 14:08:11.03743+07	2026-07-08 14:08:11.03743+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	8	38d8b397-7bc6-42d8-ad25-6376dd42a935	2026-07-08 14:08:11.038557+07	2026-07-08 14:08:11.038557+07
5e22f25f-9248-4c1f-a086-faeb657510c9	district	108	37638290-93e8-451a-bbd0-93957644d24d	2026-07-08 14:08:11.040693+07	2026-07-08 14:08:11.040693+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	district	108	efea5e90-3ad3-43b2-92f3-7bc954df7ad0	2026-07-08 14:08:11.042435+07	2026-07-08 14:08:11.042435+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	district	108	9ea6638d-7ba3-4451-b775-a14cce8f8e37	2026-07-08 14:08:11.043707+07	2026-07-08 14:08:11.043707+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	6753df4e-a515-40b6-aee3-c405d422b4d1	2026-07-08 14:08:11.04504+07	2026-07-08 14:08:11.04504+07
a24c7797-3d0a-4190-8c14-bac2bd870c09	district	108	c3939483-8179-4bb6-ae9f-6f323bb616e1	2026-07-08 14:08:11.046256+07	2026-07-08 14:08:11.046256+07
32a4a387-a0a7-4659-ac09-0fec6cbe969d	district	108	ba0e924f-4e43-4664-8f6b-4b34a9ddbcfa	2026-07-08 14:08:11.048038+07	2026-07-08 14:08:11.048038+07
39ed7343-f0d0-4b00-a34a-87b2a21a39df	district	108	3e6be84b-42df-4d28-afc1-28a4f9f65a02	2026-07-08 14:08:11.050091+07	2026-07-08 14:08:11.050091+07
193230c9-c583-4441-b648-70e06822b306	district	108	0bf165cb-399e-449b-8e9d-02c1fa612859	2026-07-08 14:08:11.051484+07	2026-07-08 14:08:11.051484+07
44580a41-cc28-4b71-ae62-8a76bb14915b	district	108	9b92cb7f-8129-46bf-9136-f09439bdde52	2026-07-08 14:08:11.057127+07	2026-07-08 14:08:11.057127+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	3dca5b86-ffa4-46cf-bc7b-4d96cf56bb3e	2026-07-08 14:08:11.058715+07	2026-07-08 14:08:11.058715+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	1	377f68ee-c1e7-47bf-8ea4-027a8657f6f4	2026-07-09 14:46:27.486166+07	2026-07-09 14:46:27.486166+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	2	a1ac9140-5737-41cf-a704-5a37e673a151	2026-07-09 14:46:27.488524+07	2026-07-09 14:46:27.488524+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	3	068826e4-e2b1-4f63-bd4d-e9b92389c4c6	2026-07-09 14:46:27.490241+07	2026-07-09 14:46:27.490241+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	4	16afa24e-5392-43cd-b74e-7a15f3147f00	2026-07-09 14:46:27.491776+07	2026-07-09 14:46:27.491776+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	5	fc95de27-345a-4515-8ec4-22d5529a0476	2026-07-09 14:46:27.493297+07	2026-07-09 14:46:27.493297+07
a24c7797-3d0a-4190-8c14-bac2bd870c09	district	6	134e7b53-24ec-4ad3-b40a-e675baaf7e61	2026-07-09 14:46:27.494738+07	2026-07-09 14:46:27.494738+07
a24c7797-3d0a-4190-8c14-bac2bd870c09	district	7	26af681a-35b0-46a9-a440-2ea7a59bf069	2026-07-09 14:46:27.496213+07	2026-07-09 14:46:27.496213+07
a24c7797-3d0a-4190-8c14-bac2bd870c09	district	8	bc81a266-7f7a-49f0-b94a-05d62129e5f9	2026-07-09 14:46:27.497554+07	2026-07-09 14:46:27.497554+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	district	108	a454c71b-f0b2-428c-82fd-edc04bf0468f	2026-07-09 14:46:27.499825+07	2026-07-09 14:46:27.499825+07
a24c7797-3d0a-4190-8c14-bac2bd870c09	district	108	04a8e135-e5a2-48c5-9075-1fbc6e2bdc74	2026-07-09 14:46:27.501369+07	2026-07-09 14:46:27.501369+07
32a4a387-a0a7-4659-ac09-0fec6cbe969d	district	108	6a85512e-701f-46fd-9400-b58770618df4	2026-07-09 14:46:27.502635+07	2026-07-09 14:46:27.502635+07
39ed7343-f0d0-4b00-a34a-87b2a21a39df	district	108	78d3b598-0bbb-4f4c-b30d-559c683d06a6	2026-07-09 14:46:27.504189+07	2026-07-09 14:46:27.504189+07
193230c9-c583-4441-b648-70e06822b306	district	108	51ec92b3-3579-4d32-8826-cf3cd3a54fb4	2026-07-09 14:46:27.505428+07	2026-07-09 14:46:27.505428+07
44580a41-cc28-4b71-ae62-8a76bb14915b	district	108	168ce89b-b255-4ada-a019-b890bfeabccf	2026-07-09 14:46:27.506694+07	2026-07-09 14:46:27.506694+07
3f818e72-6d7a-4552-ba97-10d4540c1257	district	108	bf5d73a0-c9fa-481c-bbb1-43f819c17394	2026-07-09 14:46:27.50853+07	2026-07-09 14:46:27.50853+07
\.


--
-- Data for Name: instructor_profiles; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.instructor_profiles (id, user_id, license_number, bnsp_certificate_number, years_of_experience, bio, license_expiry, photo_url, is_active, number_of_students, sessions_completed, average_rating, created_at, updated_at, description, specialization) FROM stdin;
8fc36431-41bd-46e3-bc5a-bba1f0ea924d	81c28015-a5a8-4f66-82ad-e3b204a88550	DRIVING-12345	BNSP-2024-123456	10	Experienced driving instructor with 10+ years of experience	2026-12-31 07:00:00+07	https://example.com/images/instructor1-profile.jpg	t	150	500	4.8	2026-05-31 21:35:31.717523+07	2026-05-31 21:35:31.717523+07	\N	\N
0871377d-b938-415d-8ab3-a6e9d49db14d	4ea68344-c0c9-4a76-bb2f-2356154887b9	DRIVING-67890	BNSP-2024-789012	5	Professional driving instructor specializing in defensive driving	2027-06-30 07:00:00+07	https://example.com/images/instructor2-profile.jpg	t	85	300	4.6	2026-05-31 21:35:31.765425+07	2026-05-31 21:35:31.765425+07	\N	\N
bd248f22-d6b2-4033-8744-0082d4eaa459	162ea8ed-9505-4cdc-b1ae-b1bf0c5a6fef	DRIVING-12345	BNSP-2024-123456	10	Experienced driving instructor with 10+ years of experience	2026-12-31 07:00:00+07	https://example.com/images/instructor1-profile.jpg	t	150	500	4.8	2026-06-04 20:13:04.98779+07	2026-06-04 20:13:04.98779+07	\N	\N
21d39be9-c907-4121-b7fa-491130cc141e	e8a0b2eb-834b-47ff-85e0-35455885ec7c	DRIVING-67890	BNSP-2024-789012	5	Professional driving instructor specializing in defensive driving	2027-06-30 07:00:00+07	https://example.com/images/instructor2-profile.jpg	t	85	300	4.6	2026-06-04 20:13:05.003006+07	2026-06-04 20:13:05.003006+07	\N	\N
1dc07a02-a0e9-4aef-9ac9-82a52228f480	9312f5d2-3a7a-47cf-aaf9-8b04b9a1d5fc	9876543210987654	BNSP-2026-01-1234567	4		0001-01-01 07:07:12+07:07:12		t	0	0	0	2026-06-05 10:12:49.046377+07	2026-06-05 10:12:49.046408+07		
d5bbd215-389f-4710-9e81-f49cf3909b44	c6db92c2-4891-4bef-9b54-ed6e6d8977b7	192881472875	BNSP-102-4984	5		0001-01-01 07:07:12+07:07:12		t	0	0	0	2026-06-05 11:07:43.888733+07	2026-06-05 11:07:43.888757+07		
6c76d3b5-5abe-49fb-afd7-760545b7861c	fa7fdae3-6bd7-4cec-9760-11c163f9ce3c	19284275	BNSP-102094-202	3	Former formula 1 driver	0001-01-01 07:07:12+07:07:12		t	0	0	0	2026-06-05 11:39:58.387549+07	2026-06-05 11:39:58.38757+07		
1e09db06-7517-4bb3-b82c-2513edf3b880	a86b502c-8ffc-441f-b9a3-4cefc8a265a5	9166494319	BNSP-018272-292926	2	Pandai mengemudi di sirkuit oval	0001-01-01 07:07:12+07:07:12		t	0	0	0	2026-06-05 11:56:36.465768+07	2026-06-05 11:56:36.465899+07		
28f4b025-d07d-4858-b121-922b02eb4efb	304114e1-0a14-43b1-963c-d73ae9e01eb3	SIM A	BNSP-107-2026	10	DriveMaster memposisikan diri\nsebagai Driving Education System yang membantu siswa membangun keterampilan, kepercayaan diri,\nkesiapan mental, dan keselamatan berkendara.	0001-01-01 07:07:12+07:07:12	https://ik.imagekit.io/oy4rsvid5/instructors/profiles/sumarlin1784101552178_1784101554_l0cqJ1Wsa?tr=w-700,h-450,fo-auto	t	0	0	5	0001-01-01 07:07:12+07:07:12	2026-08-04 09:35:43.51575+07	DriveMaster memposisikan diri\nsebagai Driving Education System yang membantu siswa membangun keterampilan, kepercayaan diri,\nkesiapan mental, dan keselamatan berkendara.	
5d10b025-816f-4c12-ac3e-f8f4010f892e	3f818e72-6d7a-4552-ba97-10d4540c1257	SIM A / 201884485	BNSP-2024-789012	5	Saya adalah pemenang	2027-06-30 07:00:00+07	https://ik.imagekit.io/oy4rsvid5/instructors/profiles/3f818e72-6d7a-4552-ba97-10d4540c1257_1783392065_obEmAx3tJl?tr=w-700,h-450,fo-auto	t	85	300	4.607907490115303	0001-01-01 07:07:12+07:07:12	2026-07-30 09:46:18.814147+07	Saya adalah pemenang	
\.


--
-- Data for Name: instructor_recurring_schedules; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.instructor_recurring_schedules (id, instructor_id, day_of_week, start_time, end_time, is_active, created_at, updated_at) FROM stdin;
3ca36a3a-096b-4e78-8437-98252095b673	3f818e72-6d7a-4552-ba97-10d4540c1257	1	10:00	11:00	t	2026-06-13 13:36:59.369539+07	2026-06-13 13:36:59.369539+07
4cd15bcb-bb2d-42f0-822c-4a2b40d5181b	3f818e72-6d7a-4552-ba97-10d4540c1257	1	10:00	11:00	t	2026-06-13 17:18:40.317558+07	2026-06-13 17:18:40.317558+07
535f63be-7255-45ea-94d8-402d2ba9b4ba	3f818e72-6d7a-4552-ba97-10d4540c1257	2	13:00	14:00	t	2026-06-16 22:20:03.902855+07	2026-06-16 22:20:03.902855+07
c920cb39-7de8-4652-bead-c66a32b56ddf	3f818e72-6d7a-4552-ba97-10d4540c1257	3	17:00	18:00	t	2026-06-17 16:50:52.241151+07	2026-06-17 16:50:52.241151+07
d49a2a1e-2916-4e67-a6e6-ef5050e8ca0e	3f818e72-6d7a-4552-ba97-10d4540c1257	5	13:00	14:00	t	2026-06-19 10:12:46.834912+07	2026-06-19 10:12:46.834912+07
\.


--
-- Data for Name: member_profiles; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.member_profiles (id, user_id, sessions_completed, training_time, average_rating, created_at, updated_at, identity_full_name, identity_fullname) FROM stdin;
277c6841-89a7-40fb-be3f-860e1479292d	86f5e950-cedb-495b-9528-48a7dffa6919	12	600	4.5	2026-05-31 21:35:31.798377+07	2026-05-31 21:35:31.798377+07	\N	
e4cea64f-4db9-42f4-acc5-d4946bb41dfa	433c4e17-83cb-4133-a4ab-7abaef8a3afe	8	400	4.2	2026-05-31 21:35:31.82081+07	2026-05-31 21:35:31.82081+07	\N	
fe5dc733-c6ed-4438-9270-b4e287a2d04b	0ccae045-9154-41c6-957c-c42f812b809f	0	0	0	2026-06-05 13:37:47.777131+07	2026-06-05 13:37:47.777131+07	\N	
7cf645bf-4ff2-4bcf-b441-4a87a77a5342	910e2940-b1f9-4ce4-b1cf-3ef7ccc30a6f	0	0	0	2026-06-05 13:52:50.793449+07	2026-06-05 13:52:50.793449+07	\N	
d5ae2f27-b6d5-4659-8eff-aa24e3ddebb2	ea3f6d8e-c49c-4971-8175-74108b7cfdb7	0	0	0	2026-06-05 14:22:40.514142+07	2026-06-05 14:22:40.514142+07	\N	
b33ce195-71d7-4910-9aae-ad66887a73fc	1d8756e8-ef86-4c9c-a45c-5f88c0af06cc	0	0	0	2026-06-16 22:26:05.862847+07	2026-06-16 22:26:05.862847+07	\N	
b6b6d2c1-6d1b-40f0-a8ce-0b397b2a078d	41bf31f3-555d-4ba6-b271-e16ca5d0f79d	0	0	0	2026-06-19 00:11:55.724593+07	2026-06-19 00:11:55.724593+07	\N	
8f6c4916-ce49-4daa-8b5a-f1a38b532417	73e2d7ce-9277-4c05-934c-689fae652b7a	0	0	0	2026-06-19 09:19:37.877809+07	2026-06-19 09:19:37.877809+07	\N	
294ac81f-a5b2-4e00-8d02-dfb8993dcc21	780d7b35-a240-4079-893d-d1d4a4533be4	0	0	0	2026-06-19 10:59:09.055063+07	2026-06-19 10:59:09.055063+07	\N	
4c0819c1-1ef0-4baa-b122-3ceb47ef0ca1	f6cc6a1c-43c0-41f0-bbde-ad779496c0a6	0	0	0	2026-06-19 14:06:45.10353+07	2026-06-19 14:06:45.10353+07	\N	
ee0d31d6-0dbf-4d33-967c-d331623c3cc7	58e96cd2-bbec-432e-8f23-1b38eeef98ac	0	0	0	2026-06-19 23:34:55.921318+07	2026-06-19 23:34:55.921318+07	\N	
dc79f8ce-8e8f-4823-8a05-aa25d3e8d456	16df7f78-2040-4df9-b8c7-27bc11aea48e	0	0	0	2026-06-24 14:39:03.387505+07	2026-06-24 14:39:03.387505+07	\N	
200493aa-4009-4267-a1cc-14a669ffc8f2	ce232507-fc5c-492f-bc8d-b53a5720ea98	0	0	0	2026-06-24 23:57:50.114944+07	2026-06-24 23:57:50.114944+07	\N	
73df6aa6-b919-435d-bc0b-728d1f71fdf0	a71dfa8f-432c-4651-8f61-458f9eb1236a	0	0	0	2026-06-25 00:01:32.238641+07	2026-06-25 00:01:32.238641+07	\N	
ed7d0ea0-11b2-468b-b167-ae34dfbe3622	b093730c-8caa-4238-960d-6eb7db20a21f	0	0	0	2026-06-25 01:38:00.022207+07	2026-06-25 01:38:00.022207+07	\N	
62691b88-51e7-42b0-9479-2978133095fe	79772999-f8d7-4667-8450-089af5e1507c	0	0	0	2026-06-25 01:52:19.000526+07	2026-06-25 01:52:19.000526+07	\N	
65f084d7-5482-4172-9177-272bc33b4086	37024e14-e42f-4e02-b7d5-8488886fa9a9	0	0	0	2026-06-26 14:08:37.720296+07	2026-06-26 14:09:31.59058+07	\N	Muhammad Rizqiko Harliano
defbf690-1235-4205-a99d-a43de4aca9c6	e1bd24fe-3bd5-4485-ae6b-1d8945a44b5d	0	0	0	2026-07-03 13:16:14.804874+07	2026-07-03 13:42:09.854594+07	\N	Sumanadi Rahmanadika
5038544b-edec-4e86-99cc-55fbf4478772	61b9e452-7b55-46d4-bc8f-23a954a7ace8	0	0	0	2026-07-30 18:14:37.265055+07	2026-07-30 18:15:59.108057+07	\N	Sumanadi Rahmanadhika
9f2f96d5-cb1a-45c2-be7d-4b5beb34cd58	54f2b8ea-d33e-401d-a8e2-ee9705e76d9c	0	0	0	2026-06-25 02:08:23.052781+07	2026-07-21 08:59:28.47081+07	\N	Addzkiya Nailah
45140754-65f9-4411-9852-2a48eacdb58a	fbb44c95-79aa-421e-ba3c-b4e7b87668d8	2	0	0	2026-07-30 18:43:27.000772+07	2026-07-30 19:07:20.12986+07	\N	Sumanadi Rahmanadhika
5ac2e497-74c0-498e-a114-793882b7b3ec	fead48e2-b1b6-4620-9b18-889e64355698	2	0	0	2026-08-04 10:35:04.471775+07	2026-08-04 10:38:40.235982+07	\N	Afifah
4d9c9de3-3e61-448a-b475-51f20aa2f41e	3c09cb50-3d99-4ff8-a57f-44fd4135d070	3	0	0	2026-07-15 14:51:33.205487+07	2026-08-04 10:55:31.502661+07	\N	Eki Afifah Rahmawati
3ae4b9f9-e137-4405-b958-81692a5a9225	f366e131-07a8-496f-8bb9-ca5ab5786f4c	0	0	0	2026-08-04 18:37:01.155815+07	2026-08-04 18:39:20.085669+07	\N	Sumanadi Rahmanadhika
17581c6a-94b7-4d48-8ea9-7284dec6dd94	447c7653-04dc-4550-a703-6c81847b1dd6	4	0	0	2026-07-29 11:30:54.13943+07	2026-08-05 18:00:24.822634+07	\N	Sumanadi Ra
20e11bc3-6489-418e-abc5-7a56b710bff1	a428042e-e617-48df-8688-cb9ffa7f8c32	1	0	0	2026-07-29 13:31:44.7522+07	2026-07-30 09:45:50.792325+07	\N	Rizki
08c6183e-40ba-442d-98ac-9e194246cc6e	99999999-9999-9999-9999-999999999999	12	450	4.8	2026-07-06 13:39:33.51076+07	2026-07-30 10:02:55.476585+07	\N	
0d519297-c04c-478c-88a2-b7e7fae243ca	e3747d57-fa14-4c52-b4b5-bb9de614cfc2	5	300	4.8	2026-07-30 10:10:06.894142+07	2026-07-30 10:10:06.894142+07	Gold Member 2	Gold Member 2
b2c1d3e4-f5e6-7890-a1b2-c3d4e5f67899	a2b1c3d4-e5f6-7890-a1b2-c3d4e5f67899	5	250	4.5	2026-07-30 10:23:56.870932+07	2026-07-30 10:23:56.870932+07	\N	Gold Member 2
\.


--
-- Data for Name: roles; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.roles (id, name, created_at, updated_at) FROM stdin;
1	viewer	2026-05-31 21:35:31.341671+07	2026-05-31 21:35:31.341671+07
2	member	2026-05-31 21:35:31.34846+07	2026-05-31 21:35:31.34846+07
3	instructor	2026-05-31 21:35:31.353324+07	2026-05-31 21:35:31.353324+07
4	admin	2026-05-31 21:35:31.359222+07	2026-05-31 21:35:31.359222+07
5	super_admin	2026-05-31 21:35:31.36436+07	2026-05-31 21:35:31.36436+07
\.


--
-- Data for Name: testimonial_media; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.testimonial_media (id, url, media_type, sort_order, created_at, testimonial_id) FROM stdin;
\.


--
-- Data for Name: testimonials; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.testimonials (user_id, user_name, user_image, user_role, content, rating, tags, status, is_featured, added_by, added_at, sort_order, created_at, updated_at, id) FROM stdin;
54f2b8ea-d33e-401d-a8e2-ee9705e76d9c	Addzkiya Nailah		Student	Asik	5.0	Hebat 	published	f	cf475ead-91bf-4a55-a47e-cf93279240b6	0001-01-01 07:07:12+07:07:12	0	2026-06-26 19:51:35.243486+07	2026-08-05 23:44:31.604738+07	f8d23f3c-66c2-4a9a-9b03-7f9834ea3532
433c4e17-83cb-4133-a4ab-7abaef8a3afe	Rizal Kokona		Student	Menakjubkan	4.0	Damai	archived	t	cf475ead-91bf-4a55-a47e-cf93279240b6	0001-01-01 07:07:12+07:07:12	0	2026-06-26 19:47:45.449475+07	2026-08-05 23:44:40.242108+07	64104fbe-ab0c-47b5-aca4-a5e53d58d85f
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.users (id, first_name, last_name, username, password_hash, email_address, phone_number, image, date_of_birth, address, is_active, is_verified, role_id, created_at, updated_at) FROM stdin;
cf475ead-91bf-4a55-a47e-cf93279240b6	Sandy	Arkosh	admin	$2a$10$R25OU7rPerFtwxAPbQeAs.vTaaMG.3nWG0w1MwM9XcMyfg6hGLKaG	admin@example.com	+6281234567890	https://example.com/images/admin.jpg	1990-01-01	Jakarta, Indonesia	t	f	4	2026-05-31 21:35:31.697342+07	2026-05-31 21:35:31.697342+07
70ba6361-a61a-427e-8cdb-9b9cee4f1858	Super	Admin	super_admin	$2a$10$R25OU7rPerFtwxAPbQeAs.vTaaMG.3nWG0w1MwM9XcMyfg6hGLKaG	superadmin@example.com	+6281234567891	https://example.com/images/super_admin.jpg	1985-06-15	Jakarta, Indonesia	t	f	5	2026-05-31 21:35:31.704734+07	2026-05-31 21:35:31.704734+07
86f5e950-cedb-495b-9528-48a7dffa6919	Gabriella	Dianasari	member1	$2a$10$R25OU7rPerFtwxAPbQeAs.vTaaMG.3nWG0w1MwM9XcMyfg6hGLKaG	member1@example.com	+6281234567894	https://example.com/images/member1.jpg	2000-11-25	Jakarta, Indonesia	t	f	2	2026-05-31 21:35:31.79431+07	2026-05-31 21:35:31.79431+07
433c4e17-83cb-4133-a4ab-7abaef8a3afe	Rizal	Kokona	member2	$2a$10$R25OU7rPerFtwxAPbQeAs.vTaaMG.3nWG0w1MwM9XcMyfg6hGLKaG	member2@example.com	+6281234567895	https://example.com/images/member2.jpg	1998-04-08	Bekasi, Indonesia	t	f	2	2026-05-31 21:35:31.81649+07	2026-05-31 21:35:31.81649+07
0ccae045-9154-41c6-957c-c42f812b809f	Zeta		zetarahmadi	$2a$10$hj5TwSNXuSouoVScMvXwGeE6n/SOTH9z/zz7zSzn/emKtRz/MNfr6	zetarahmadi@gmail.com	08172646656		2026-06-05		t	f	2	2026-06-05 13:37:47.766663+07	2026-06-05 13:37:47.766663+07
910e2940-b1f9-4ce4-b1cf-3ef7ccc30a6f	Akira		adrilareza	$2a$10$Yx1t5GjCwrf4lkUzz5K08.AVgaUGQafNpGpHe/XlqFOyiKqs2kIrK	adrilareza@gmail.com	0817266264		2026-06-05		t	f	2	2026-06-05 13:52:50.781332+07	2026-06-05 13:52:50.781332+07
58e96cd2-bbec-432e-8f23-1b38eeef98ac	Syawali	Astra	syawaliastra436	$2a$10$XgMaVeWKyZZETbDmQgFe7Oo9hmfTUZLB4Yak8TaQBjzkBoRF5S8YW	syawaliastra436@gmail.com	081826876463		2026-06-19		t	f	2	2026-06-19 23:34:55.898964+07	2026-06-19 23:34:55.898964+07
54f2b8ea-d33e-401d-a8e2-ee9705e76d9c	Addzkiya	Nailah	addzkiyanailah	$2a$10$h.SFH7v9rPa7lAthzyEV/uaZk9qQuGNxPM24EzuqFlZeKlTnJMsxq	addzkiyanailah@gmail.com	081862764726		2026-06-25		t	f	2	2026-06-25 02:08:23.045336+07	2026-06-25 02:08:23.045336+07
37024e14-e42f-4e02-b7d5-8488886fa9a9	Qiko	1999	qiko1999	$2a$10$.axNCmAli45vs5pq1op7D.UCzShrxnctrGLz68xjz1/lr1IDl9KjS	qiko1999@gmail.com	09197287274		2026-06-26		t	f	2	2026-06-26 14:08:37.542263+07	2026-06-26 14:08:37.542263+07
e1bd24fe-3bd5-4485-ae6b-1d8945a44b5d	Ambrizal	Rais	sumanadir	$2a$10$//v.upS1qJOjszclFgM/e.Fv1xAkZdCVNNAggvyf5xFq48CLkRB.6	sumanadir@gmail.com	0885629294		2026-07-03		t	f	2	2026-07-03 13:16:14.796896+07	2026-07-03 13:16:14.796896+07
99999999-9999-9999-9999-999999999999	Gold	Member	goldmember	$2a$10$yft16uwBpPbXDqca4WYa2OgQkRPQg3NAMYsusBNuCp99nMVrOx1LO	goldmember@example.com	+6281234567899	https://example.com/images/goldmember.jpg	1995-05-05	Jakarta, Indonesia	t	f	2	2026-07-06 13:39:33.50676+07	2026-07-06 13:39:33.50676+07
3f818e72-6d7a-4552-ba97-10d4540c1257	Irena	Diah	instructor2	$2a$10$9wJIhmPNZ0BPggrQyIWXV.tyMkzzLgExnbmQcIvAxqtnqKBDbF5Ly	instructor2@example.com	+6281234567893	https://ik.imagekit.io/oy4rsvid5/instructors/profiles/3f818e72-6d7a-4552-ba97-10d4540c1257_1783392065_obEmAx3tJl?tr=w-700,h-450,fo-auto	1992-07-10	Surabaya, Indonesia	t	f	3	2026-06-04 20:24:05.761501+07	2026-07-07 09:41:12.538556+07
f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	Budi	Santoso	budi1780643919040	$2a$10$Wsk32uHQvna7RJLjYZgQtOUV1.FQzBjCh7KbvSXauunDjpmn71p5K	budi@example.com	081316278284		2026-06-05		t	f	2	2026-06-05 14:18:39.188306+07	2026-07-10 17:04:47.34109+07
cb89c300-dd63-49fc-953a-e53ba79d2f79	Dimas	Saifullah	dimas1780635518088	$2a$10$c9gvIZH5IaZRMyCjv7StCOFS60oTdCL1KltbjolKy3em2oFmd/C1S	dimas@example.com	0812667264		2026-06-05		t	f	2	2026-06-05 11:58:38.267206+07	2026-07-10 17:04:51.721239+07
5e22f25f-9248-4c1f-a086-faeb657510c9	Rahmat	Rahkoda	instructor1	$2a$10$9wJIhmPNZ0BPggrQyIWXV.tyMkzzLgExnbmQcIvAxqtnqKBDbF5Ly	instructor1@example.com	+6281234567892	https://example.com/images/instructor1.jpg	1988-03-20	Bandung, Indonesia	t	f	2	2026-06-04 20:24:05.702563+07	2026-07-10 17:04:56.371686+07
0765a33f-6633-4a4f-8a99-3f8b1495ecd0	Gilang	Pambudi	gilang1780635184551	$2a$10$34jq73ybVvu1UZWWPMq8dO63grKX1t7PvTUrq21rA4HYHvka6TbiC	gilang@example.com	0818266264		2026-06-05		t	f	2	2026-06-05 11:53:04.761292+07	2026-07-10 17:04:59.839387+07
a24c7797-3d0a-4190-8c14-bac2bd870c09	Ikmal	Hamdallah	ikmal1782366905869	$2a$10$YcrHSmHyuvQJbReVhH734eAXiBujUcMWNZ3ndG2iMk9ScnpWtKQ76	ikmal@example.com	082186276273		2026-06-25		t	f	2	2026-06-25 12:55:06.023139+07	2026-07-10 17:05:06.99723+07
a2b1c3d4-e5f6-7890-a1b2-c3d4e5f67899	Gold	Member 2	goldmember2@gmail.com	$2a$10$R25OU7rPerFtwxAPbQeAs.vTaaMG.3nWG0w1MwM9XcMyfg6hGLKaG	goldmember2@gmail.com	+6281234567892	\N	2000-01-01	\N	t	t	2	2026-07-30 10:23:56.870932+07	2026-07-30 10:23:56.870932+07
3c09cb50-3d99-4ff8-a57f-44fd4135d070	Eki Afifah	Rahmawati	afifahrahmawati221	$2a$10$Zk5/NuW.pc/naCES0LMPtumnezojezIQUMls9JOCDUXhWRNBkIpkS	afifahrahmawati221@gmail.com	081298239776		2026-07-15		t	f	2	2026-07-15 14:51:33.180691+07	2026-07-15 14:51:33.180691+07
a428042e-e617-48df-8688-cb9ffa7f8c32	Rizki	H	muhammadrizqiko	$2a$10$pm5DtP2LSBI8KOtLfbUz5OkoNg.aWsarz.0BIehZGFJkPhtx.2Uea	muhammadrizqiko@gmail.com	081312168535		2026-07-29		t	f	2	2026-07-29 13:31:44.749193+07	2026-07-29 13:31:44.749193+07
fbb44c95-79aa-421e-ba3c-b4e7b87668d8	Sumanadi	Rahmanadhika	dika942013	$2a$10$klZYovBI.76pHibj4IGJMOZtONvQwXoSYG3kiarQKHD3ajnY1LiY6	dika942013@gmail.com	081234567890		2026-07-30		t	f	2	2026-07-30 18:43:26.997393+07	2026-07-30 18:43:26.997393+07
304114e1-0a14-43b1-963c-d73ae9e01eb3	Sumarlin	Marlin	sumarlin1784101552178	$2a$10$80RvscEonjpsATzvZy74IOKejHwFIl4ma9mjB6HTxaXsf.BiU3PYe	sumarlin@example.com	089602580111	https://ik.imagekit.io/oy4rsvid5/instructors/profiles/sumarlin1784101552178_1784101554_l0cqJ1Wsa?tr=w-700,h-450,fo-auto	2026-07-15		t	f	3	2026-07-15 14:45:56.131088+07	2026-08-04 09:35:43.509366+07
fead48e2-b1b6-4620-9b18-889e64355698	Afifah	Afifah	administrasi	$2a$10$.wVFkOGT0dg1k0rWhTupheZrWD6C/TDw9uHQFaJ2y31Tbay7qxrui	administrasi@wicgrowth.com	081298239776		2026-08-04		t	f	2	2026-08-04 10:35:04.466041+07	2026-08-04 10:35:04.466041+07
\.


--
-- Data for Name: work_experiences; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.work_experiences (id, instructor_id, company_name, role, start_date, end_date, description, is_verified, created_at, updated_at) FROM stdin;
1	81c28015-a5a8-4f66-82ad-e3b204a88550	Premier Driving School	Senior Driving Instructor	2015-01-15 07:00:00+07	2020-06-30 07:00:00+07	Led driving education programs for new drivers. Specialized in defensive driving techniques.	f	2026-05-31 21:35:31.831103+07	2026-05-31 21:35:31.831103+07
2	81c28015-a5a8-4f66-82ad-e3b204a88550	Safe Wheels Academy	Driving Instructor	2012-03-01 07:00:00+07	2014-12-31 07:00:00+07	Provided driving lessons for students of all ages and skill levels.	f	2026-05-31 21:35:31.838078+07	2026-05-31 21:35:31.838078+07
3	81c28015-a5a8-4f66-82ad-e3b204a88550	National Traffic Safety Council	Road Safety Trainer	2010-06-01 07:00:00+07	2012-02-28 07:00:00+07	Conducted road safety training and awareness programs.	f	2026-05-31 21:35:31.848222+07	2026-05-31 21:35:31.848222+07
4	81c28015-a5a8-4f66-82ad-e3b204a88550	Elite Auto School	Head Instructor	2020-07-01 07:00:00+07	0001-01-01 07:07:12+07:07:12	Currently leading a team of driving instructors. Focus on quality assurance and curriculum development.	f	2026-05-31 21:35:31.852388+07	2026-05-31 21:35:31.852388+07
5	81c28015-a5a8-4f66-82ad-e3b204a88550	Metro Driving Center	Driving Instructor	2018-09-01 07:00:00+07	2023-08-31 07:00:00+07	Taught defensive driving and provided practical driving lessons.	f	2026-05-31 21:35:31.859998+07	2026-05-31 21:35:31.859998+07
6	4ea68344-c0c9-4a76-bb2f-2356154887b9	City Driving Academy	Driving Instructor	2019-02-01 07:00:00+07	0001-01-01 07:07:12+07:07:12	Providing comprehensive driving lessons with focus on city driving and traffic awareness.	f	2026-05-31 21:35:31.86334+07	2026-05-31 21:35:31.86334+07
7	4ea68344-c0c9-4a76-bb2f-2356154887b9	Highway Safety Institute	Defensive Driving Instructor	2017-06-15 07:00:00+07	2018-12-31 07:00:00+07	Specialized in highway driving techniques and emergency handling.	f	2026-05-31 21:35:31.865659+07	2026-05-31 21:35:31.865659+07
8	162ea8ed-9505-4cdc-b1ae-b1bf0c5a6fef	Premier Driving School	Senior Driving Instructor	2015-01-15 07:00:00+07	2020-06-30 07:00:00+07	Led driving education programs for new drivers. Specialized in defensive driving techniques.	f	2026-06-04 20:13:05.025837+07	2026-06-04 20:13:05.025837+07
9	162ea8ed-9505-4cdc-b1ae-b1bf0c5a6fef	Safe Wheels Academy	Driving Instructor	2012-03-01 07:00:00+07	2014-12-31 07:00:00+07	Provided driving lessons for students of all ages and skill levels.	f	2026-06-04 20:13:05.073855+07	2026-06-04 20:13:05.073855+07
10	162ea8ed-9505-4cdc-b1ae-b1bf0c5a6fef	National Traffic Safety Council	Road Safety Trainer	2010-06-01 07:00:00+07	2012-02-28 07:00:00+07	Conducted road safety training and awareness programs.	f	2026-06-04 20:13:05.080225+07	2026-06-04 20:13:05.080225+07
11	162ea8ed-9505-4cdc-b1ae-b1bf0c5a6fef	Elite Auto School	Head Instructor	2020-07-01 07:00:00+07	0001-01-01 07:07:12+07:07:12	Currently leading a team of driving instructors. Focus on quality assurance and curriculum development.	f	2026-06-04 20:13:05.087549+07	2026-06-04 20:13:05.087549+07
12	162ea8ed-9505-4cdc-b1ae-b1bf0c5a6fef	Metro Driving Center	Driving Instructor	2018-09-01 07:00:00+07	2023-08-31 07:00:00+07	Taught defensive driving and provided practical driving lessons.	f	2026-06-04 20:13:05.091625+07	2026-06-04 20:13:05.091625+07
13	e8a0b2eb-834b-47ff-85e0-35455885ec7c	City Driving Academy	Driving Instructor	2019-02-01 07:00:00+07	0001-01-01 07:07:12+07:07:12	Providing comprehensive driving lessons with focus on city driving and traffic awareness.	f	2026-06-04 20:13:05.096251+07	2026-06-04 20:13:05.096251+07
14	e8a0b2eb-834b-47ff-85e0-35455885ec7c	Highway Safety Institute	Defensive Driving Instructor	2017-06-15 07:00:00+07	2018-12-31 07:00:00+07	Specialized in highway driving techniques and emergency handling.	f	2026-06-04 20:13:05.101035+07	2026-06-04 20:13:05.101035+07
15	5e22f25f-9248-4c1f-a086-faeb657510c9	Premier Driving School	Senior Driving Instructor	2015-01-15 07:00:00+07	2020-06-30 07:00:00+07	Led driving education programs for new drivers. Specialized in defensive driving techniques.	f	2026-06-04 20:24:05.785333+07	2026-06-04 20:24:05.785333+07
16	5e22f25f-9248-4c1f-a086-faeb657510c9	Safe Wheels Academy	Driving Instructor	2012-03-01 07:00:00+07	2014-12-31 07:00:00+07	Provided driving lessons for students of all ages and skill levels.	f	2026-06-04 20:24:05.793885+07	2026-06-04 20:24:05.793885+07
17	5e22f25f-9248-4c1f-a086-faeb657510c9	National Traffic Safety Council	Road Safety Trainer	2010-06-01 07:00:00+07	2012-02-28 07:00:00+07	Conducted road safety training and awareness programs.	f	2026-06-04 20:24:05.933925+07	2026-06-04 20:24:05.933925+07
18	5e22f25f-9248-4c1f-a086-faeb657510c9	Elite Auto School	Head Instructor	2020-07-01 07:00:00+07	0001-01-01 07:07:12+07:07:12	Currently leading a team of driving instructors. Focus on quality assurance and curriculum development.	f	2026-06-04 20:24:05.984682+07	2026-06-04 20:24:05.984682+07
19	5e22f25f-9248-4c1f-a086-faeb657510c9	Metro Driving Center	Driving Instructor	2018-09-01 07:00:00+07	2023-08-31 07:00:00+07	Taught defensive driving and provided practical driving lessons.	f	2026-06-04 20:24:05.992159+07	2026-06-04 20:24:05.992159+07
20	3f818e72-6d7a-4552-ba97-10d4540c1257	City Driving Academy	Driving Instructor	2019-02-01 07:00:00+07	0001-01-01 07:07:12+07:07:12	Providing comprehensive driving lessons with focus on city driving and traffic awareness.	f	2026-06-04 20:24:06.045569+07	2026-06-04 20:24:06.045569+07
21	3f818e72-6d7a-4552-ba97-10d4540c1257	Highway Safety Institute	Defensive Driving Instructor	2017-06-15 07:00:00+07	2018-12-31 07:00:00+07	Specialized in highway driving techniques and emergency handling.	f	2026-06-04 20:24:06.050575+07	2026-06-04 20:24:06.050575+07
22	0765a33f-6633-4a4f-8a99-3f8b1495ecd0	City Driving Academy	Driving Instructor	2019-02-01 07:00:00+07	0001-01-01 07:07:12+07:07:12	Providing comprehensive driving lessons with focus on city driving and traffic awareness.	f	2026-07-08 14:08:11.02074+07	2026-07-08 14:08:11.02074+07
23	0765a33f-6633-4a4f-8a99-3f8b1495ecd0	Highway Safety Institute	Defensive Driving Instructor	2017-06-15 07:00:00+07	2018-12-31 07:00:00+07	Specialized in highway driving techniques and emergency handling.	f	2026-07-08 14:08:11.024969+07	2026-07-08 14:08:11.024969+07
24	f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	Premier Driving School	Senior Driving Instructor	2015-01-15 07:00:00+07	2020-06-30 07:00:00+07	Led driving education programs for new drivers. Specialized in defensive driving techniques.	f	2026-07-09 14:46:27.464185+07	2026-07-09 14:46:27.464185+07
25	f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	Safe Wheels Academy	Driving Instructor	2012-03-01 07:00:00+07	2014-12-31 07:00:00+07	Provided driving lessons for students of all ages and skill levels.	f	2026-07-09 14:46:27.469908+07	2026-07-09 14:46:27.469908+07
26	f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	National Traffic Safety Council	Road Safety Trainer	2010-06-01 07:00:00+07	2012-02-28 07:00:00+07	Conducted road safety training and awareness programs.	f	2026-07-09 14:46:27.472196+07	2026-07-09 14:46:27.472196+07
27	f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	Elite Auto School	Head Instructor	2020-07-01 07:00:00+07	0001-01-01 07:07:12+07:07:12	Currently leading a team of driving instructors. Focus on quality assurance and curriculum development.	f	2026-07-09 14:46:27.47492+07	2026-07-09 14:46:27.47492+07
28	f475ab4e-6dfd-4bb3-b7d2-4b6d6790e652	Metro Driving Center	Driving Instructor	2018-09-01 07:00:00+07	2023-08-31 07:00:00+07	Taught defensive driving and provided practical driving lessons.	f	2026-07-09 14:46:27.477114+07	2026-07-09 14:46:27.477114+07
29	a24c7797-3d0a-4190-8c14-bac2bd870c09	City Driving Academy	Driving Instructor	2019-02-01 07:00:00+07	0001-01-01 07:07:12+07:07:12	Providing comprehensive driving lessons with focus on city driving and traffic awareness.	f	2026-07-09 14:46:27.478748+07	2026-07-09 14:46:27.478748+07
30	a24c7797-3d0a-4190-8c14-bac2bd870c09	Highway Safety Institute	Defensive Driving Instructor	2017-06-15 07:00:00+07	2018-12-31 07:00:00+07	Specialized in highway driving techniques and emergency handling.	f	2026-07-09 14:46:27.480132+07	2026-07-09 14:46:27.480132+07
\.


--
-- Name: roles_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.roles_id_seq', 5, true);


--
-- Name: testimonial_media_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.testimonial_media_id_seq', 1, false);


--
-- Name: work_experiences_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.work_experiences_id_seq', 30, true);


--
-- Name: certifications certifications_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.certifications
    ADD CONSTRAINT certifications_pkey PRIMARY KEY (id);


--
-- Name: entitlements entitlements_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.entitlements
    ADD CONSTRAINT entitlements_pkey PRIMARY KEY (id);


--
-- Name: instructor_profiles instructor_profiles_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.instructor_profiles
    ADD CONSTRAINT instructor_profiles_pkey PRIMARY KEY (id);


--
-- Name: instructor_recurring_schedules instructor_recurring_schedules_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.instructor_recurring_schedules
    ADD CONSTRAINT instructor_recurring_schedules_pkey PRIMARY KEY (id);


--
-- Name: member_profiles member_profiles_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.member_profiles
    ADD CONSTRAINT member_profiles_pkey PRIMARY KEY (id);


--
-- Name: roles roles_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.roles
    ADD CONSTRAINT roles_pkey PRIMARY KEY (id);


--
-- Name: testimonial_media testimonial_media_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.testimonial_media
    ADD CONSTRAINT testimonial_media_pkey PRIMARY KEY (id);


--
-- Name: testimonials testimonials_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.testimonials
    ADD CONSTRAINT testimonials_pkey PRIMARY KEY (id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: work_experiences work_experiences_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.work_experiences
    ADD CONSTRAINT work_experiences_pkey PRIMARY KEY (id);


--
-- Name: idx_certifications_entitlement_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_certifications_entitlement_id ON public.certifications USING btree (entitlement_id);


--
-- Name: idx_certifications_instructor_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_certifications_instructor_id ON public.certifications USING btree (instructor_id);


--
-- Name: idx_certifications_member_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_certifications_member_id ON public.certifications USING btree (member_id);


--
-- Name: idx_entitlements_booking_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_entitlements_booking_id ON public.entitlements USING btree (booking_id);


--
-- Name: idx_entitlements_member_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_entitlements_member_id ON public.entitlements USING btree (member_id);


--
-- Name: idx_instructor_areas_instructor_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_instructor_areas_instructor_id ON public.instructor_areas USING btree (instructor_id);


--
-- Name: idx_instructor_profiles_user_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX idx_instructor_profiles_user_id ON public.instructor_profiles USING btree (user_id);


--
-- Name: idx_instructor_recurring_schedules_instructor_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_instructor_recurring_schedules_instructor_id ON public.instructor_recurring_schedules USING btree (instructor_id);


--
-- Name: idx_member_profiles_user_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX idx_member_profiles_user_id ON public.member_profiles USING btree (user_id);


--
-- Name: idx_roles_name; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX idx_roles_name ON public.roles USING btree (name);


--
-- Name: idx_testimonial_media_testimonial_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_testimonial_media_testimonial_id ON public.testimonial_media USING btree (testimonial_id);


--
-- Name: idx_testimonials_user_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_testimonials_user_id ON public.testimonials USING btree (user_id);


--
-- Name: idx_users_email_address; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX idx_users_email_address ON public.users USING btree (email_address);


--
-- Name: idx_users_username; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX idx_users_username ON public.users USING btree (username);


--
-- Name: testimonial_media fk_testimonial_media_testimonial; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.testimonial_media
    ADD CONSTRAINT fk_testimonial_media_testimonial FOREIGN KEY (testimonial_id) REFERENCES public.testimonials(id);


--
-- PostgreSQL database dump complete
--

\unrestrict Ubfhn5snx0teniGilnnys7MoxEz2It4dMom8nOPupXxXgevPtKOIussYMqNNDCS

