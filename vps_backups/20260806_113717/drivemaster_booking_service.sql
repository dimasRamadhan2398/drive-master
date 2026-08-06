--
-- PostgreSQL database dump
--

\restrict 1SsoONKiZieIBhMK0TbdHBfiGumUhrTqgFbKlIca2EuIq6KmwfZmaYjU2lB8fWd

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

ALTER TABLE ONLY public.transaction_items DROP CONSTRAINT fk_transactions_items;
ALTER TABLE ONLY public.user_entitlements DROP CONSTRAINT fk_enrollments_entitlements;
DROP INDEX public.idx_user_entitlements_user_id;
DROP INDEX public.idx_user_entitlements_enrollment_id;
DROP INDEX public.idx_user_entitlements_anonymized_at;
DROP INDEX public.idx_transactions_user_id;
DROP INDEX public.idx_transactions_enrollment_id;
DROP INDEX public.idx_transaction_items_transaction_id;
DROP INDEX public.idx_schedules_user_id;
DROP INDEX public.idx_schedules_instructor_id;
DROP INDEX public.idx_schedules_enrollment_id;
DROP INDEX public.idx_schedules_date;
DROP INDEX public.idx_schedules_car_id;
DROP INDEX public.idx_payments_user_id;
DROP INDEX public.idx_payments_order_id;
DROP INDEX public.idx_payments_enrollment_id;
DROP INDEX public.idx_enrollments_user_id;
DROP INDEX public.idx_enrollments_package_id;
DROP INDEX public.idx_enrollments_anonymized_at;
DROP INDEX public.idx_driving_sessions_user_id;
DROP INDEX public.idx_driving_sessions_schedule_id;
DROP INDEX public.idx_driving_sessions_instructor_id;
DROP INDEX public.idx_driving_sessions_entitlement_id;
DROP INDEX public.idx_driving_sessions_enrollment_id;
DROP INDEX public.idx_driving_sessions_car_id;
DROP INDEX public.idx_driving_sessions_anonymized_at;
ALTER TABLE ONLY public.user_entitlements DROP CONSTRAINT user_entitlements_pkey;
ALTER TABLE ONLY public.user_certifications DROP CONSTRAINT user_certifications_pkey;
ALTER TABLE ONLY public.transactions DROP CONSTRAINT transactions_pkey;
ALTER TABLE ONLY public.transaction_items DROP CONSTRAINT transaction_items_pkey;
ALTER TABLE ONLY public.schedules DROP CONSTRAINT schedules_pkey;
ALTER TABLE ONLY public.payments DROP CONSTRAINT payments_pkey;
ALTER TABLE ONLY public.enrollments DROP CONSTRAINT enrollments_pkey;
ALTER TABLE ONLY public.driving_sessions DROP CONSTRAINT driving_sessions_pkey;
ALTER TABLE public.schedules ALTER COLUMN id DROP DEFAULT;
ALTER TABLE public.payments ALTER COLUMN id DROP DEFAULT;
ALTER TABLE public.driving_sessions ALTER COLUMN id DROP DEFAULT;
DROP TABLE public.user_entitlements;
DROP TABLE public.user_certifications;
DROP TABLE public.transactions;
DROP TABLE public.transaction_items;
DROP SEQUENCE public.schedules_id_seq;
DROP TABLE public.schedules;
DROP SEQUENCE public.payments_id_seq;
DROP TABLE public.payments;
DROP TABLE public.enrollments;
DROP SEQUENCE public.driving_sessions_id_seq;
DROP TABLE public.driving_sessions;
SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: driving_sessions; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.driving_sessions (
    id bigint NOT NULL,
    enrollment_id uuid NOT NULL,
    entitlement_id uuid NOT NULL,
    user_id uuid NOT NULL,
    instructor_id uuid NOT NULL,
    car_id uuid NOT NULL,
    schedule_id bigint,
    date date NOT NULL,
    "time" character varying(10) NOT NULL,
    duration bigint DEFAULT 60,
    status character varying(20) DEFAULT 'scheduled'::character varying,
    area character varying(150),
    notes text,
    anonymized_at timestamp with time zone,
    started_at timestamp with time zone,
    completed_at timestamp with time zone,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    end_time timestamp with time zone,
    is_ended_by_admin boolean DEFAULT false,
    rating numeric(2,1),
    feedback text
);


ALTER TABLE public.driving_sessions OWNER TO postgres;

--
-- Name: driving_sessions_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.driving_sessions_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.driving_sessions_id_seq OWNER TO postgres;

--
-- Name: driving_sessions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.driving_sessions_id_seq OWNED BY public.driving_sessions.id;


--
-- Name: enrollments; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.enrollments (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    package_id uuid NOT NULL,
    status character varying(30) DEFAULT 'pending_payment'::character varying,
    total_price numeric,
    paid_at timestamp with time zone,
    expires_at timestamp with time zone,
    anonymized_at timestamp with time zone,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.enrollments OWNER TO postgres;

--
-- Name: payments; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.payments (
    id bigint NOT NULL,
    enrollment_id uuid NOT NULL,
    user_id uuid NOT NULL,
    order_id character varying(100),
    amount numeric NOT NULL,
    payment_method character varying(30),
    status character varying(30) DEFAULT 'pending'::character varying,
    payment_url text,
    transaction_id character varying(100),
    paid_at timestamp with time zone,
    expires_at timestamp with time zone,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.payments OWNER TO postgres;

--
-- Name: payments_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.payments_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.payments_id_seq OWNER TO postgres;

--
-- Name: payments_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.payments_id_seq OWNED BY public.payments.id;


--
-- Name: schedules; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.schedules (
    id bigint NOT NULL,
    date date NOT NULL,
    "time" character varying(10) NOT NULL,
    duration bigint DEFAULT 60,
    instructor_id uuid NOT NULL,
    user_id text,
    enrollment_id text,
    status character varying(20) DEFAULT 'available'::character varying,
    notes text,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    car_id uuid
);


ALTER TABLE public.schedules OWNER TO postgres;

--
-- Name: schedules_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.schedules_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.schedules_id_seq OWNER TO postgres;

--
-- Name: schedules_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.schedules_id_seq OWNED BY public.schedules.id;


--
-- Name: transaction_items; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.transaction_items (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    transaction_id uuid NOT NULL,
    item_type character varying(20) NOT NULL,
    item_id uuid NOT NULL,
    item_name character varying(255) NOT NULL,
    quantity bigint DEFAULT 1,
    unit_price numeric(10,2) NOT NULL,
    subtotal numeric(10,2) NOT NULL,
    sessions bigint DEFAULT 0,
    created_at timestamp with time zone
);


ALTER TABLE public.transaction_items OWNER TO postgres;

--
-- Name: transactions; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.transactions (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    enrollment_id uuid NOT NULL,
    user_id uuid NOT NULL,
    base_price numeric(10,2) DEFAULT 0,
    add_ons_total numeric(10,2) DEFAULT 0,
    total_amount numeric(10,2) NOT NULL,
    status character varying(20) DEFAULT 'pending'::character varying,
    paid_at timestamp with time zone,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.transactions OWNER TO postgres;

--
-- Name: user_certifications; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.user_certifications (
    user_id bigint NOT NULL,
    certification_id bigint NOT NULL,
    issued_at timestamp with time zone
);


ALTER TABLE public.user_certifications OWNER TO postgres;

--
-- Name: user_entitlements; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.user_entitlements (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    enrollment_id uuid NOT NULL,
    user_id uuid NOT NULL,
    source_type character varying(50) NOT NULL,
    source_id character varying(100) NOT NULL,
    total_sessions bigint,
    used_sessions bigint,
    expires_at timestamp with time zone,
    anonymized_at timestamp with time zone,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    CONSTRAINT chk_used_sessions CHECK ((used_sessions <= total_sessions))
);


ALTER TABLE public.user_entitlements OWNER TO postgres;

--
-- Name: driving_sessions id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.driving_sessions ALTER COLUMN id SET DEFAULT nextval('public.driving_sessions_id_seq'::regclass);


--
-- Name: payments id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.payments ALTER COLUMN id SET DEFAULT nextval('public.payments_id_seq'::regclass);


--
-- Name: schedules id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.schedules ALTER COLUMN id SET DEFAULT nextval('public.schedules_id_seq'::regclass);


--
-- Data for Name: driving_sessions; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.driving_sessions (id, enrollment_id, entitlement_id, user_id, instructor_id, car_id, schedule_id, date, "time", duration, status, area, notes, anonymized_at, started_at, completed_at, created_at, updated_at, end_time, is_ended_by_admin, rating, feedback) FROM stdin;
10	3ff22915-363d-4896-a9d9-b5695a6fdecf	6a62af4e-b950-4aa5-a294-f6ea58856b15	99999999-9999-9999-9999-999999999999	3f818e72-6d7a-4552-ba97-10d4540c1257	00000000-0000-0000-0000-000000000000	1940	2026-07-17	13:00	60	completed		Booked via dashboard schedule page	\N	2026-07-28 11:43:54.510851+07	2026-07-28 11:43:54.51671+07	2026-07-14 09:51:30.18952+07	2026-07-28 11:43:54.51671+07	\N	f	\N	\N
6	3ff22915-363d-4896-a9d9-b5695a6fdecf	6a62af4e-b950-4aa5-a294-f6ea58856b15	99999999-9999-9999-9999-999999999999	3f818e72-6d7a-4552-ba97-10d4540c1257	00000000-0000-0000-0000-000000000000	1937	2026-07-13	10:00	60	completed		Booked via dashboard schedule page	\N	2026-07-13 11:04:47.826901+07	2026-07-13 11:05:05.277124+07	2026-07-13 09:34:17.867541+07	2026-07-13 11:05:05.277125+07	\N	f	\N	\N
25	3ff22915-363d-4896-a9d9-b5695a6fdecf	6a62af4e-b950-4aa5-a294-f6ea58856b15	99999999-9999-9999-9999-999999999999	3f818e72-6d7a-4552-ba97-10d4540c1257	898170a5-08db-4467-b33e-049660a4231c	1973	2026-07-30	15:00	60	completed		Booked via dashboard schedule page	\N	2026-07-30 10:02:47.597699+07	2026-07-30 10:02:54.938464+07	2026-07-30 10:02:37.295287+07	2026-07-30 10:02:54.938468+07	2026-07-30 10:02:54.938464+07	t	\N	
8	3ff22915-363d-4896-a9d9-b5695a6fdecf	6a62af4e-b950-4aa5-a294-f6ea58856b15	99999999-9999-9999-9999-999999999999	3f818e72-6d7a-4552-ba97-10d4540c1257	00000000-0000-0000-0000-000000000000	1939	2026-07-15	17:00	60	completed		Booked via dashboard schedule page	\N	2026-07-28 11:43:55.024663+07	2026-07-28 11:43:55.02986+07	2026-07-13 11:13:32.807781+07	2026-07-28 11:43:55.029861+07	\N	f	\N	\N
18	3ff22915-363d-4896-a9d9-b5695a6fdecf	6a62af4e-b950-4aa5-a294-f6ea58856b15	99999999-9999-9999-9999-999999999999	304114e1-0a14-43b1-963c-d73ae9e01eb3	898170a5-08db-4467-b33e-049660a4231c	1951	2026-07-25	08:00	60	cancelled		Booked via dashboard schedule page	\N	\N	\N	2026-07-24 14:49:18.932471+07	2026-07-24 15:05:36.725519+07	\N	f	\N	\N
7	3ff22915-363d-4896-a9d9-b5695a6fdecf	6a62af4e-b950-4aa5-a294-f6ea58856b15	99999999-9999-9999-9999-999999999999	3f818e72-6d7a-4552-ba97-10d4540c1257	00000000-0000-0000-0000-000000000000	1938	2026-07-14	13:00	60	completed		Booked via dashboard schedule page	\N	2026-07-28 11:43:55.537084+07	2026-07-28 11:43:55.55439+07	2026-07-13 11:10:31.847273+07	2026-07-28 11:43:55.554391+07	\N	f	\N	\N
15	3c3c0ead-befa-475f-819a-574a45e39543	a9de4846-b6ab-4a0e-beb4-efec1a27d1cf	3c09cb50-3d99-4ff8-a57f-44fd4135d070	304114e1-0a14-43b1-963c-d73ae9e01eb3	898170a5-08db-4467-b33e-049660a4231c	1954	2026-07-27	08:00	60	completed		Booked via dashboard schedule page	\N	2026-07-27 09:43:28.121849+07	2026-07-27 09:43:36.560784+07	2026-07-24 13:16:43.913819+07	2026-07-30 08:18:31.586045+07	\N	f	5.0	
16	3c3c0ead-befa-475f-819a-574a45e39543	a9de4846-b6ab-4a0e-beb4-efec1a27d1cf	3c09cb50-3d99-4ff8-a57f-44fd4135d070	304114e1-0a14-43b1-963c-d73ae9e01eb3	898170a5-08db-4467-b33e-049660a4231c	1950	2026-07-24	08:00	60	completed		Booked via dashboard schedule page	\N	2026-07-24 13:22:02.225973+07	2026-07-24 13:22:11.044974+07	2026-07-24 13:17:41.208386+07	2026-07-30 08:18:36.202589+07	\N	f	5.0	
27	a0c14f0b-9d06-4a9d-83b2-2ed8dff7ceca	cbf419f2-9a86-4729-bf4a-4973aeb4afc2	fbb44c95-79aa-421e-ba3c-b4e7b87668d8	304114e1-0a14-43b1-963c-d73ae9e01eb3	c0000001-0000-0000-0000-000000000008	1974	2026-08-06	08:00	60	completed		Booked via dashboard schedule page	\N	2026-07-30 19:07:04.070378+07	2026-07-30 19:07:15.000795+07	2026-07-30 19:06:36.75085+07	2026-07-30 19:07:15.000796+07	2026-07-30 19:07:15.000795+07	t	\N	
28	a0c14f0b-9d06-4a9d-83b2-2ed8dff7ceca	cbf419f2-9a86-4729-bf4a-4973aeb4afc2	fbb44c95-79aa-421e-ba3c-b4e7b87668d8	304114e1-0a14-43b1-963c-d73ae9e01eb3	c0000001-0000-0000-0000-000000000008	1975	2026-08-07	08:00	60	completed		Booked via dashboard schedule page	\N	2026-07-30 19:07:09.782237+07	2026-07-30 19:07:19.611427+07	2026-07-30 19:06:43.929138+07	2026-07-30 19:07:19.611428+07	2026-07-30 19:07:19.611427+07	t	\N	
11	3ff22915-363d-4896-a9d9-b5695a6fdecf	6a62af4e-b950-4aa5-a294-f6ea58856b15	99999999-9999-9999-9999-999999999999	3f818e72-6d7a-4552-ba97-10d4540c1257	00000000-0000-0000-0000-000000000000	1943	2026-07-20	10:00	60	completed		Booked via dashboard schedule page	\N	2026-07-28 11:43:53.445868+07	2026-07-28 11:43:53.467503+07	2026-07-14 15:25:57.154304+07	2026-07-28 11:43:53.467504+07	\N	f	\N	\N
9	3ff22915-363d-4896-a9d9-b5695a6fdecf	6a62af4e-b950-4aa5-a294-f6ea58856b15	99999999-9999-9999-9999-999999999999	3f818e72-6d7a-4552-ba97-10d4540c1257	00000000-0000-0000-0000-000000000000	1940	2026-07-17	13:00	60	completed		Booked via dashboard schedule page	\N	2026-07-28 11:43:53.997645+07	2026-07-28 11:43:54.003395+07	2026-07-14 09:51:29.661739+07	2026-07-28 11:43:54.003395+07	\N	f	\N	\N
17	3c3c0ead-befa-475f-819a-574a45e39543	a9de4846-b6ab-4a0e-beb4-efec1a27d1cf	3c09cb50-3d99-4ff8-a57f-44fd4135d070	304114e1-0a14-43b1-963c-d73ae9e01eb3	898170a5-08db-4467-b33e-049660a4231c	1958	2026-07-24	14:00	60	completed		Booked via dashboard schedule page	\N	2026-07-24 13:24:02.823704+07	2026-07-24 13:24:08.761277+07	2026-07-24 13:18:55.685356+07	2026-07-30 08:18:40.558596+07	\N	f	5.0	
20	dd8a6ee9-e2d5-4e2f-820f-256a843d429d	a1d8731c-75c5-45cd-bffa-6d7abf1ec1e3	447c7653-04dc-4550-a703-6c81847b1dd6	3f818e72-6d7a-4552-ba97-10d4540c1257	00000000-0000-0000-0000-000000000000	1966	2026-08-03	10:00	60	cancelled			\N	\N	\N	2026-07-29 12:34:11.352654+07	2026-07-29 12:34:37.616023+07	\N	f	\N	
12	3c3c0ead-befa-475f-819a-574a45e39543	a9de4846-b6ab-4a0e-beb4-efec1a27d1cf	3c09cb50-3d99-4ff8-a57f-44fd4135d070	3f818e72-6d7a-4552-ba97-10d4540c1257	00000000-0000-0000-0000-000000000000	1947	2026-07-27	10:00	60	completed		Booked via dashboard schedule page	\N	2026-07-27 09:43:45.790745+07	2026-07-27 09:43:52.250337+07	2026-07-24 13:14:08.121335+07	2026-07-30 08:19:03.006318+07	\N	f	5.0	
14	3c3c0ead-befa-475f-819a-574a45e39543	a9de4846-b6ab-4a0e-beb4-efec1a27d1cf	3c09cb50-3d99-4ff8-a57f-44fd4135d070	3f818e72-6d7a-4552-ba97-10d4540c1257	00000000-0000-0000-0000-000000000000	1949	2026-07-29	17:00	60	completed		Booked via dashboard schedule page	\N	2026-07-27 13:59:44.624155+07	2026-07-29 18:00:19.109363+07	2026-07-24 13:15:23.006914+07	2026-07-30 08:19:10.874986+07	2026-07-29 18:00:19.109363+07	f	5.0	
21	dd8a6ee9-e2d5-4e2f-820f-256a843d429d	a1d8731c-75c5-45cd-bffa-6d7abf1ec1e3	447c7653-04dc-4550-a703-6c81847b1dd6	3f818e72-6d7a-4552-ba97-10d4540c1257	00000000-0000-0000-0000-000000000000	1966	2026-08-03	10:00	60	completed			\N	2026-07-29 12:35:36.128641+07	2026-07-29 12:35:47.592238+07	2026-07-29 12:35:18.553434+07	2026-07-29 12:36:00.462412+07	2026-07-29 12:35:47.592238+07	t	5.0	
22	dd8a6ee9-e2d5-4e2f-820f-256a843d429d	a1d8731c-75c5-45cd-bffa-6d7abf1ec1e3	447c7653-04dc-4550-a703-6c81847b1dd6	3f818e72-6d7a-4552-ba97-10d4540c1257	00000000-0000-0000-0000-000000000000	1967	2026-08-04	13:00	60	cancelled			\N	\N	\N	2026-07-29 12:38:13.025434+07	2026-07-29 12:39:17.426518+07	\N	f	\N	
13	3c3c0ead-befa-475f-819a-574a45e39543	a9de4846-b6ab-4a0e-beb4-efec1a27d1cf	3c09cb50-3d99-4ff8-a57f-44fd4135d070	3f818e72-6d7a-4552-ba97-10d4540c1257	00000000-0000-0000-0000-000000000000	1948	2026-07-28	13:00	60	completed		Booked via dashboard schedule page	\N	2026-07-27 13:59:25.903094+07	2026-07-28 14:00:07.013668+07	2026-07-24 13:15:10.152573+07	2026-07-30 08:19:22.440359+07	\N	f	5.0	
23	dd8a6ee9-e2d5-4e2f-820f-256a843d429d	a1d8731c-75c5-45cd-bffa-6d7abf1ec1e3	447c7653-04dc-4550-a703-6c81847b1dd6	3f818e72-6d7a-4552-ba97-10d4540c1257	00000000-0000-0000-0000-000000000000	1967	2026-08-04	13:00	60	completed			\N	2026-07-29 12:40:43.952775+07	2026-07-29 12:40:58.482379+07	2026-07-29 12:40:20.146109+07	2026-07-29 12:41:11.094567+07	2026-07-29 12:40:58.482379+07	t	5.0	
29	9737a24d-0d15-48ed-897f-c0df00d57b31	ff6577f6-2375-4cc0-abd7-a64cec8b79f2	fead48e2-b1b6-4620-9b18-889e64355698	3f818e72-6d7a-4552-ba97-10d4540c1257	00000000-0000-0000-0000-000000000000	1976	2026-08-07	13:00	60	completed		Booked via dashboard schedule page	\N	2026-08-04 10:38:13.01195+07	2026-08-04 10:38:20.699942+07	2026-08-04 10:37:00.329633+07	2026-08-04 10:38:20.699942+07	2026-08-04 10:38:20.699942+07	t	\N	
19	dd8a6ee9-e2d5-4e2f-820f-256a843d429d	a1d8731c-75c5-45cd-bffa-6d7abf1ec1e3	447c7653-04dc-4550-a703-6c81847b1dd6	3f818e72-6d7a-4552-ba97-10d4540c1257	00000000-0000-0000-0000-000000000000	1965	2026-07-31	13:00	60	completed		Booked via dashboard schedule page	\N	2026-07-31 13:00:08.525951+07	2026-07-31 14:00:08.549682+07	2026-07-29 12:31:09.762382+07	2026-07-31 14:00:08.549683+07	2026-07-31 14:00:08.549682+07	f	\N	
24	7fdb8006-33bd-4116-a104-88076e236d0f	1de4436b-2a2f-4351-9e0b-2fd78966c2d9	a428042e-e617-48df-8688-cb9ffa7f8c32	3f818e72-6d7a-4552-ba97-10d4540c1257	c0000001-0000-0000-0000-000000000006	1972	2026-07-30	10:00	60	completed		Booked via dashboard schedule page	\N	2026-07-30 09:45:35.677308+07	2026-07-30 09:45:50.785087+07	2026-07-30 09:45:23.359417+07	2026-07-30 09:46:18.810382+07	2026-07-30 09:45:50.785087+07	t	5.0	
30	9737a24d-0d15-48ed-897f-c0df00d57b31	ff6577f6-2375-4cc0-abd7-a64cec8b79f2	fead48e2-b1b6-4620-9b18-889e64355698	3f818e72-6d7a-4552-ba97-10d4540c1257	00000000-0000-0000-0000-000000000000	1977	2026-08-10	10:00	60	completed		Booked via dashboard schedule page	\N	2026-08-04 10:38:31.501234+07	2026-08-04 10:38:40.220773+07	2026-08-04 10:37:11.099708+07	2026-08-04 10:38:40.220774+07	2026-08-04 10:38:40.220773+07	t	\N	
31	a723de49-a8e2-4d78-884f-203dd4522910	51b0bb38-3259-4e2b-951e-27e7d5608d70	3c09cb50-3d99-4ff8-a57f-44fd4135d070	304114e1-0a14-43b1-963c-d73ae9e01eb3	c0000001-0000-0000-0000-000000000005	1978	2026-08-10	12:00	60	scheduled		Booked via dashboard schedule page	\N	\N	\N	2026-08-04 10:40:58.551171+07	2026-08-04 10:40:58.551171+07	\N	f	\N	
32	a723de49-a8e2-4d78-884f-203dd4522910	51b0bb38-3259-4e2b-951e-27e7d5608d70	3c09cb50-3d99-4ff8-a57f-44fd4135d070	304114e1-0a14-43b1-963c-d73ae9e01eb3	c0000001-0000-0000-0000-000000000005	1979	2026-08-10	14:00	60	scheduled		Booked via dashboard schedule page	\N	\N	\N	2026-08-04 10:41:07.350698+07	2026-08-04 10:41:07.350698+07	\N	f	\N	
33	5ed13eae-af47-486d-98a2-0160e99c5290	dbaedb67-dc5a-4839-8a9c-b15166650baf	3c09cb50-3d99-4ff8-a57f-44fd4135d070	3f818e72-6d7a-4552-ba97-10d4540c1257	c0000001-0000-0000-0000-000000000007	1980	2026-08-10	15:00	60	completed		Booked via dashboard schedule page	\N	2026-08-04 10:55:05.527558+07	2026-08-04 10:55:14.139036+07	2026-08-04 10:54:34.052972+07	2026-08-04 10:55:14.139037+07	2026-08-04 10:55:14.139036+07	t	\N	
34	5ed13eae-af47-486d-98a2-0160e99c5290	dbaedb67-dc5a-4839-8a9c-b15166650baf	3c09cb50-3d99-4ff8-a57f-44fd4135d070	3f818e72-6d7a-4552-ba97-10d4540c1257	c0000001-0000-0000-0000-000000000006	1981	2026-08-10	17:00	60	completed		Booked via dashboard schedule page	\N	2026-08-04 10:55:22.465071+07	2026-08-04 10:55:31.490645+07	2026-08-04 10:54:44.031056+07	2026-08-04 10:55:31.490646+07	2026-08-04 10:55:31.490645+07	t	\N	
26	dd8a6ee9-e2d5-4e2f-820f-256a843d429d	a1d8731c-75c5-45cd-bffa-6d7abf1ec1e3	447c7653-04dc-4550-a703-6c81847b1dd6	3f818e72-6d7a-4552-ba97-10d4540c1257	00000000-0000-0000-0000-000000000000	1968	2026-08-05	17:00	60	completed		Booked via dashboard schedule page	\N	2026-08-05 17:00:24.766242+07	2026-08-05 18:00:24.791051+07	2026-07-30 13:28:31.075099+07	2026-08-05 18:00:24.791052+07	2026-08-05 18:00:24.791051+07	f	\N	
\.


--
-- Data for Name: enrollments; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.enrollments (id, user_id, package_id, status, total_price, paid_at, expires_at, anonymized_at, created_at, updated_at) FROM stdin;
b1c2a1da-b04c-4b1b-bc28-eab853854204	54f2b8ea-d33e-401d-a8e2-ee9705e76d9c	11111111-1111-1111-1111-111111111201	pending	2700000	\N	2027-07-15 07:55:54.950665+07	\N	2026-07-15 07:55:54.955195+07	2026-07-15 07:55:54.955195+07
62f454d9-70db-4fc4-8dae-4e582b24c752	3c09cb50-3d99-4ff8-a57f-44fd4135d070	11111111-1111-1111-1111-111111111101	pending	2150000	\N	2027-07-17 10:24:10.542902+07	\N	2026-07-17 10:24:10.543546+07	2026-07-17 10:24:10.543546+07
5d0a2fba-e002-4397-a062-43fd0644205b	54f2b8ea-d33e-401d-a8e2-ee9705e76d9c	11111111-1111-1111-1111-111111111101	pending	2400000	\N	2027-07-17 10:39:29.544488+07	\N	2026-07-17 10:39:29.548012+07	2026-07-17 10:39:29.548012+07
f5f4326a-50b1-494d-af01-7b31bfbbde2d	3c09cb50-3d99-4ff8-a57f-44fd4135d070	11111111-1111-1111-1111-111111111101	pending	2150000	\N	2027-07-20 10:14:17.115418+07	\N	2026-07-20 10:14:17.116844+07	2026-07-20 10:14:17.116844+07
6fe72be4-8563-45e6-9f91-15d026a7a9a7	3c09cb50-3d99-4ff8-a57f-44fd4135d070	11111111-1111-1111-1111-111111111101	pending	2150000	\N	2027-07-21 08:53:17.57362+07	\N	2026-07-21 08:53:17.575076+07	2026-07-21 08:53:17.575076+07
0e927e20-ccf8-4144-8fc4-4191ce83f643	54f2b8ea-d33e-401d-a8e2-ee9705e76d9c	11111111-1111-1111-1111-111111111101	pending	2400000	\N	2027-07-21 08:59:35.055136+07	\N	2026-07-21 08:59:35.057694+07	2026-07-21 08:59:35.057694+07
ae8cb132-ca78-4049-a030-e7f2ab6c0a5a	3c09cb50-3d99-4ff8-a57f-44fd4135d070	11111111-1111-1111-1111-111111111101	pending	2150000	\N	2027-07-23 08:31:16.181457+07	\N	2026-07-23 08:31:16.182292+07	2026-07-23 08:31:16.182292+07
3c3c0ead-befa-475f-819a-574a45e39543	3c09cb50-3d99-4ff8-a57f-44fd4135d070	11111111-1111-1111-1111-111111111101	completed	2150000.00	2026-07-22 16:42:46.981945+07	2027-07-22 12:56:54.709282+07	\N	2026-07-22 12:56:54.710081+07	2026-07-27 13:59:53.453089+07
b310062a-1787-4481-ada4-875f3e9ed5b0	447c7653-04dc-4550-a703-6c81847b1dd6	11111111-1111-1111-1111-111111111101	pending	2150000	\N	2027-07-29 11:50:20.619979+07	\N	2026-07-29 11:50:20.620557+07	2026-07-29 11:50:20.620557+07
63347228-eec3-41e4-b9e0-8ddedcc47b68	447c7653-04dc-4550-a703-6c81847b1dd6	11111111-1111-1111-1111-111111111101	pending	2150000	\N	2027-07-29 11:51:30.615618+07	\N	2026-07-29 11:51:30.616226+07	2026-07-29 11:51:30.616226+07
82df6b26-8a85-452b-a715-835d5c99f24d	447c7653-04dc-4550-a703-6c81847b1dd6	11111111-1111-1111-1111-111111111101	pending	2150000	\N	2027-07-29 11:53:51.876644+07	\N	2026-07-29 11:53:51.877396+07	2026-07-29 11:53:51.877396+07
edec3989-a2d9-4155-bbc5-3fcf11b67123	a428042e-e617-48df-8688-cb9ffa7f8c32	11111111-1111-1111-1111-111111111101	pending	2400000	\N	2027-07-29 13:32:06.128689+07	\N	2026-07-29 13:32:06.130897+07	2026-07-29 13:32:06.130897+07
7fdb8006-33bd-4116-a104-88076e236d0f	a428042e-e617-48df-8688-cb9ffa7f8c32	11111111-1111-1111-1111-111111111101	paid	10000	2026-07-30 05:19:38.108134+07	2027-07-30 05:11:51.476225+07	\N	2026-07-30 05:11:51.476754+07	2026-07-30 05:19:38.108138+07
61778c9a-7ff0-41ef-8651-bad7d2c88de2	447c7653-04dc-4550-a703-6c81847b1dd6	11111111-1111-1111-1111-111111111101	pending	2950000	\N	2027-07-30 13:15:21.056205+07	\N	2026-07-30 13:15:21.06006+07	2026-07-30 13:15:21.06006+07
e61625ec-8bbe-4a44-bf7a-116fdfa82e15	447c7653-04dc-4550-a703-6c81847b1dd6	11111111-1111-1111-1111-111111111101	pending	2950000	\N	2027-07-30 13:15:29.369165+07	\N	2026-07-30 13:15:29.373065+07	2026-07-30 13:15:29.373065+07
26eb4720-50db-4d1d-8ea7-11ab10dc6e6e	447c7653-04dc-4550-a703-6c81847b1dd6	11111111-1111-1111-1111-111111111101	pending	2150000	\N	2027-07-30 13:24:50.581193+07	\N	2026-07-30 13:24:50.581996+07	2026-07-30 13:24:50.581996+07
32c239f2-8c95-467d-ac6e-ff31e19a895e	61b9e452-7b55-46d4-bc8f-23a954a7ace8	11111111-1111-1111-1111-111111111101	pending	2150000	\N	2027-07-30 18:14:56.710661+07	\N	2026-07-30 18:14:56.711328+07	2026-07-30 18:14:56.711328+07
f5113829-32a9-4598-9d50-e994b51f980d	fbb44c95-79aa-421e-ba3c-b4e7b87668d8	4601ad1b-4dc9-4727-b229-43cdf6fddecb	pending	10000	\N	2027-07-30 18:52:49.710317+07	\N	2026-07-30 18:52:49.710837+07	2026-07-30 18:52:49.710837+07
979188ed-90b6-4390-a584-f57eba2d3e16	fbb44c95-79aa-421e-ba3c-b4e7b87668d8	4601ad1b-4dc9-4727-b229-43cdf6fddecb	pending	10000	\N	2027-07-30 18:55:24.103343+07	\N	2026-07-30 18:55:24.104892+07	2026-07-30 18:55:24.104892+07
a0c14f0b-9d06-4a9d-83b2-2ed8dff7ceca	fbb44c95-79aa-421e-ba3c-b4e7b87668d8	4601ad1b-4dc9-4727-b229-43cdf6fddecb	paid	10000	2026-07-30 19:00:13.846356+07	2027-07-30 18:56:44.796661+07	\N	2026-07-30 18:56:44.797205+07	2026-07-30 19:00:13.846357+07
dd8a6ee9-e2d5-4e2f-820f-256a843d429d	447c7653-04dc-4550-a703-6c81847b1dd6	11111111-1111-1111-1111-111111111101	in_progress	2150000	2026-07-29 11:41:11.029757+07	2027-07-29 11:37:06.180691+07	\N	2026-07-29 11:37:06.181219+07	2026-07-31 14:00:08.572189+07
69b0f270-ff15-4352-ba8d-f7f383260fca	3c09cb50-3d99-4ff8-a57f-44fd4135d070	4601ad1b-4dc9-4727-b229-43cdf6fddecb	paid	10000	2026-08-04 10:23:20.729795+07	2027-08-04 10:21:48.473555+07	\N	2026-08-04 10:21:48.474603+07	2026-08-04 10:23:20.729795+07
a723de49-a8e2-4d78-884f-203dd4522910	3c09cb50-3d99-4ff8-a57f-44fd4135d070	4601ad1b-4dc9-4727-b229-43cdf6fddecb	paid	10000	2026-08-04 10:33:04.559852+07	2027-08-04 10:24:36.406334+07	\N	2026-08-04 10:24:36.407267+07	2026-08-04 10:33:04.559856+07
9737a24d-0d15-48ed-897f-c0df00d57b31	fead48e2-b1b6-4620-9b18-889e64355698	4601ad1b-4dc9-4727-b229-43cdf6fddecb	paid	10000	2026-08-04 10:36:39.462694+07	2027-08-04 10:35:53.417612+07	\N	2026-08-04 10:35:53.418272+07	2026-08-04 10:36:39.462694+07
5ed13eae-af47-486d-98a2-0160e99c5290	3c09cb50-3d99-4ff8-a57f-44fd4135d070	4601ad1b-4dc9-4727-b229-43cdf6fddecb	paid	10000	2026-08-04 10:51:50.822169+07	2027-08-04 10:51:38.352484+07	\N	2026-08-04 10:51:38.353297+07	2026-08-04 10:51:50.82217+07
\.


--
-- Data for Name: payments; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.payments (id, enrollment_id, user_id, order_id, amount, payment_method, status, payment_url, transaction_id, paid_at, expires_at, created_at, updated_at) FROM stdin;
\.


--
-- Data for Name: schedules; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.schedules (id, date, "time", duration, instructor_id, user_id, enrollment_id, status, notes, created_at, updated_at, car_id) FROM stdin;
1937	2026-07-13	10:00	60	3f818e72-6d7a-4552-ba97-10d4540c1257	99999999-9999-9999-9999-999999999999	3ff22915-363d-4896-a9d9-b5695a6fdecf	completed		2026-07-10 10:39:37.327156+07	2026-07-13 11:05:05.009293+07	00000000-0000-0000-0000-000000000000
1950	2026-07-24	08:00	60	304114e1-0a14-43b1-963c-d73ae9e01eb3	3c09cb50-3d99-4ff8-a57f-44fd4135d070	3c3c0ead-befa-475f-819a-574a45e39543	completed		2026-07-24 09:12:52.418707+07	2026-07-24 13:22:11.039659+07	898170a5-08db-4467-b33e-049660a4231c
1940	2026-07-17	13:00	60	3f818e72-6d7a-4552-ba97-10d4540c1257	99999999-9999-9999-9999-999999999999	3ff22915-363d-4896-a9d9-b5695a6fdecf	completed		2026-07-11 00:05:00.65798+07	2026-07-28 11:43:54.513345+07	00000000-0000-0000-0000-000000000000
1958	2026-07-24	14:00	60	304114e1-0a14-43b1-963c-d73ae9e01eb3	3c09cb50-3d99-4ff8-a57f-44fd4135d070	3c3c0ead-befa-475f-819a-574a45e39543	completed		2026-07-24 13:17:29.19419+07	2026-07-24 13:24:08.757046+07	898170a5-08db-4467-b33e-049660a4231c
1954	2026-07-27	08:00	60	304114e1-0a14-43b1-963c-d73ae9e01eb3	3c09cb50-3d99-4ff8-a57f-44fd4135d070	3c3c0ead-befa-475f-819a-574a45e39543	completed		2026-07-24 13:13:13.287034+07	2026-07-27 09:43:36.556341+07	898170a5-08db-4467-b33e-049660a4231c
1964	2026-07-27	17:00	60	304114e1-0a14-43b1-963c-d73ae9e01eb3	\N	\N	blocked		2026-07-24 13:35:11.116064+07	2026-07-28 11:43:56.07661+07	898170a5-08db-4467-b33e-049660a4231c
1947	2026-07-27	10:00	60	3f818e72-6d7a-4552-ba97-10d4540c1257	3c09cb50-3d99-4ff8-a57f-44fd4135d070	3c3c0ead-befa-475f-819a-574a45e39543	completed		2026-07-21 00:05:00.022157+07	2026-07-27 09:43:52.246184+07	00000000-0000-0000-0000-000000000000
1939	2026-07-15	17:00	60	3f818e72-6d7a-4552-ba97-10d4540c1257	99999999-9999-9999-9999-999999999999	3ff22915-363d-4896-a9d9-b5695a6fdecf	completed		2026-07-10 10:39:37.498392+07	2026-07-28 11:43:55.027114+07	00000000-0000-0000-0000-000000000000
1956	2026-07-28	08:00	60	304114e1-0a14-43b1-963c-d73ae9e01eb3	\N	\N	blocked		2026-07-24 13:14:45.82457+07	2026-07-28 11:43:56.077171+07	898170a5-08db-4467-b33e-049660a4231c
1938	2026-07-14	13:00	60	3f818e72-6d7a-4552-ba97-10d4540c1257	99999999-9999-9999-9999-999999999999	3ff22915-363d-4896-a9d9-b5695a6fdecf	completed		2026-07-10 10:39:37.407545+07	2026-07-28 11:43:55.549547+07	00000000-0000-0000-0000-000000000000
1943	2026-07-20	10:00	60	3f818e72-6d7a-4552-ba97-10d4540c1257	99999999-9999-9999-9999-999999999999	3ff22915-363d-4896-a9d9-b5695a6fdecf	completed		2026-07-14 10:05:18.85504+07	2026-07-28 11:43:53.459995+07	00000000-0000-0000-0000-000000000000
1936	2026-07-10	13:00	60	3f818e72-6d7a-4552-ba97-10d4540c1257	\N	\N	blocked		2026-07-10 10:39:37.040221+07	2026-07-28 11:43:56.061837+07	00000000-0000-0000-0000-000000000000
1941	2026-07-13	13:00	60	3f818e72-6d7a-4552-ba97-10d4540c1257	\N	\N	blocked		2026-07-13 11:17:53.097652+07	2026-07-28 11:43:56.065205+07	c0000001-0000-0000-0000-000000000005
1953	2026-07-28	09:00	1	304114e1-0a14-43b1-963c-d73ae9e01eb3	\N	\N	blocked		2026-07-24 13:11:56.256496+07	2026-07-28 11:43:56.077874+07	898170a5-08db-4467-b33e-049660a4231c
1944	2026-07-21	13:00	60	3f818e72-6d7a-4552-ba97-10d4540c1257	\N	\N	blocked		2026-07-15 00:05:00.517249+07	2026-07-28 11:43:56.067694+07	00000000-0000-0000-0000-000000000000
1945	2026-07-22	17:00	60	3f818e72-6d7a-4552-ba97-10d4540c1257	\N	\N	blocked		2026-07-16 00:05:00.2744+07	2026-07-28 11:43:56.06903+07	00000000-0000-0000-0000-000000000000
1946	2026-07-24	13:00	60	3f818e72-6d7a-4552-ba97-10d4540c1257	\N	\N	blocked		2026-07-18 00:05:00.072317+07	2026-07-28 11:43:56.069981+07	00000000-0000-0000-0000-000000000000
1951	2026-07-25	08:00	60	304114e1-0a14-43b1-963c-d73ae9e01eb3	\N	\N	blocked		2026-07-24 13:10:50.263319+07	2026-07-28 11:43:56.071001+07	898170a5-08db-4467-b33e-049660a4231c
1952	2026-07-25	09:00	60	304114e1-0a14-43b1-963c-d73ae9e01eb3	\N	\N	blocked		2026-07-24 13:11:36.070707+07	2026-07-28 11:43:56.071571+07	898170a5-08db-4467-b33e-049660a4231c
1961	2026-07-27	08:00	60	304114e1-0a14-43b1-963c-d73ae9e01eb3	\N	\N	blocked		2026-07-24 13:34:01.090894+07	2026-07-28 11:43:56.072059+07	898170a5-08db-4467-b33e-049660a4231c
1957	2026-07-27	08:00	60	304114e1-0a14-43b1-963c-d73ae9e01eb3	\N	\N	blocked		2026-07-24 13:16:08.703188+07	2026-07-28 11:43:56.072589+07	c0000001-0000-0000-0000-000000000002
1960	2026-07-27	08:00	60	304114e1-0a14-43b1-963c-d73ae9e01eb3	\N	\N	blocked		2026-07-24 13:26:10.780537+07	2026-07-28 11:43:56.073525+07	898170a5-08db-4467-b33e-049660a4231c
1955	2026-07-27	10:00	60	304114e1-0a14-43b1-963c-d73ae9e01eb3	\N	\N	blocked		2026-07-24 13:13:26.486903+07	2026-07-28 11:43:56.074132+07	898170a5-08db-4467-b33e-049660a4231c
1962	2026-07-27	10:00	60	304114e1-0a14-43b1-963c-d73ae9e01eb3	\N	\N	blocked		2026-07-24 13:34:21.915441+07	2026-07-28 11:43:56.07489+07	898170a5-08db-4467-b33e-049660a4231c
1963	2026-07-27	15:00	60	304114e1-0a14-43b1-963c-d73ae9e01eb3	\N	\N	blocked		2026-07-24 13:34:53.816108+07	2026-07-28 11:43:56.076105+07	898170a5-08db-4467-b33e-049660a4231c
1959	2026-07-28	11:00	60	304114e1-0a14-43b1-963c-d73ae9e01eb3	\N	\N	blocked		2026-07-24 13:25:52.624919+07	2026-07-28 11:43:56.078357+07	898170a5-08db-4467-b33e-049660a4231c
1972	2026-07-30	10:00	60	3f818e72-6d7a-4552-ba97-10d4540c1257	a428042e-e617-48df-8688-cb9ffa7f8c32	7fdb8006-33bd-4116-a104-88076e236d0f	completed		2026-07-30 09:45:10.624164+07	2026-07-30 09:45:50.771682+07	c0000001-0000-0000-0000-000000000006
1967	2026-08-04	13:00	60	3f818e72-6d7a-4552-ba97-10d4540c1257	447c7653-04dc-4550-a703-6c81847b1dd6	dd8a6ee9-e2d5-4e2f-820f-256a843d429d	completed		2026-07-29 01:37:54.959865+07	2026-07-29 12:40:58.481568+07	00000000-0000-0000-0000-000000000000
1948	2026-07-28	13:00	60	3f818e72-6d7a-4552-ba97-10d4540c1257	3c09cb50-3d99-4ff8-a57f-44fd4135d070	3c3c0ead-befa-475f-819a-574a45e39543	completed		2026-07-22 00:05:00.034145+07	2026-07-28 14:00:07.542056+07	00000000-0000-0000-0000-000000000000
1966	2026-08-03	10:00	60	3f818e72-6d7a-4552-ba97-10d4540c1257	447c7653-04dc-4550-a703-6c81847b1dd6	dd8a6ee9-e2d5-4e2f-820f-256a843d429d	completed		2026-07-28 00:05:00.015329+07	2026-07-29 12:35:47.591456+07	00000000-0000-0000-0000-000000000000
1949	2026-07-29	17:00	60	3f818e72-6d7a-4552-ba97-10d4540c1257	3c09cb50-3d99-4ff8-a57f-44fd4135d070	3c3c0ead-befa-475f-819a-574a45e39543	completed		2026-07-23 00:05:00.019569+07	2026-07-29 18:00:19.769586+07	00000000-0000-0000-0000-000000000000
1970	2026-07-30	09:00	60	304114e1-0a14-43b1-963c-d73ae9e01eb3	\N	\N	blocked		2026-07-30 08:20:25.811178+07	2026-07-30 09:00:17.005338+07	c0000001-0000-0000-0000-000000000007
1973	2026-07-30	15:00	60	3f818e72-6d7a-4552-ba97-10d4540c1257	99999999-9999-9999-9999-999999999999	3ff22915-363d-4896-a9d9-b5695a6fdecf	completed		2026-07-30 09:56:47.078184+07	2026-07-30 10:02:54.937577+07	898170a5-08db-4467-b33e-049660a4231c
1971	2026-07-30	11:00	60	304114e1-0a14-43b1-963c-d73ae9e01eb3	\N	\N	blocked		2026-07-30 08:20:51.681916+07	2026-07-30 11:00:17.005918+07	c0000001-0000-0000-0000-000000000007
1975	2026-08-07	08:00	60	304114e1-0a14-43b1-963c-d73ae9e01eb3	fbb44c95-79aa-421e-ba3c-b4e7b87668d8	a0c14f0b-9d06-4a9d-83b2-2ed8dff7ceca	completed		2026-07-30 19:05:13.175791+07	2026-07-30 19:07:19.610494+07	c0000001-0000-0000-0000-000000000008
1974	2026-08-06	08:00	60	304114e1-0a14-43b1-963c-d73ae9e01eb3	fbb44c95-79aa-421e-ba3c-b4e7b87668d8	a0c14f0b-9d06-4a9d-83b2-2ed8dff7ceca	completed		2026-07-30 19:05:06.158306+07	2026-07-30 19:07:14.999603+07	c0000001-0000-0000-0000-000000000008
1969	2026-07-31	08:00	60	304114e1-0a14-43b1-963c-d73ae9e01eb3	\N	\N	blocked		2026-07-30 08:20:08.854207+07	2026-07-31 08:00:08.531594+07	c0000001-0000-0000-0000-000000000005
1965	2026-07-31	13:00	60	3f818e72-6d7a-4552-ba97-10d4540c1257	447c7653-04dc-4550-a703-6c81847b1dd6	dd8a6ee9-e2d5-4e2f-820f-256a843d429d	completed		2026-07-25 00:05:00.018579+07	2026-07-31 14:00:08.592248+07	00000000-0000-0000-0000-000000000000
1977	2026-08-10	10:00	60	3f818e72-6d7a-4552-ba97-10d4540c1257	fead48e2-b1b6-4620-9b18-889e64355698	9737a24d-0d15-48ed-897f-c0df00d57b31	completed		2026-08-04 00:05:00.025995+07	2026-08-04 10:38:40.219699+07	00000000-0000-0000-0000-000000000000
1976	2026-08-07	13:00	60	3f818e72-6d7a-4552-ba97-10d4540c1257	fead48e2-b1b6-4620-9b18-889e64355698	9737a24d-0d15-48ed-897f-c0df00d57b31	completed		2026-08-01 00:05:00.050917+07	2026-08-04 10:38:20.699074+07	00000000-0000-0000-0000-000000000000
1978	2026-08-10	12:00	60	304114e1-0a14-43b1-963c-d73ae9e01eb3	3c09cb50-3d99-4ff8-a57f-44fd4135d070	a723de49-a8e2-4d78-884f-203dd4522910	booked		2026-08-04 10:40:08.87459+07	2026-08-04 10:40:58.546978+07	c0000001-0000-0000-0000-000000000005
1979	2026-08-10	14:00	60	304114e1-0a14-43b1-963c-d73ae9e01eb3	3c09cb50-3d99-4ff8-a57f-44fd4135d070	a723de49-a8e2-4d78-884f-203dd4522910	booked		2026-08-04 10:40:22.542556+07	2026-08-04 10:41:07.348839+07	c0000001-0000-0000-0000-000000000005
1982	2026-08-11	13:00	60	3f818e72-6d7a-4552-ba97-10d4540c1257	\N	\N	available		2026-08-05 00:05:00.042593+07	2026-08-05 00:05:00.042593+07	00000000-0000-0000-0000-000000000000
1980	2026-08-10	15:00	60	3f818e72-6d7a-4552-ba97-10d4540c1257	3c09cb50-3d99-4ff8-a57f-44fd4135d070	5ed13eae-af47-486d-98a2-0160e99c5290	completed		2026-08-04 10:53:23.616882+07	2026-08-04 10:55:14.138017+07	c0000001-0000-0000-0000-000000000007
1981	2026-08-10	17:00	60	3f818e72-6d7a-4552-ba97-10d4540c1257	3c09cb50-3d99-4ff8-a57f-44fd4135d070	5ed13eae-af47-486d-98a2-0160e99c5290	completed		2026-08-04 10:53:39.077366+07	2026-08-04 10:55:31.489366+07	c0000001-0000-0000-0000-000000000006
1968	2026-08-05	17:00	60	3f818e72-6d7a-4552-ba97-10d4540c1257	447c7653-04dc-4550-a703-6c81847b1dd6	dd8a6ee9-e2d5-4e2f-820f-256a843d429d	completed		2026-07-30 00:05:00.046353+07	2026-08-05 18:00:24.830022+07	00000000-0000-0000-0000-000000000000
1983	2026-08-12	17:00	60	3f818e72-6d7a-4552-ba97-10d4540c1257	\N	\N	available		2026-08-06 00:05:00.03995+07	2026-08-06 00:05:00.03995+07	00000000-0000-0000-0000-000000000000
\.


--
-- Data for Name: transaction_items; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.transaction_items (id, transaction_id, item_type, item_id, item_name, quantity, unit_price, subtotal, sessions, created_at) FROM stdin;
\.


--
-- Data for Name: transactions; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.transactions (id, enrollment_id, user_id, base_price, add_ons_total, total_amount, status, paid_at, created_at, updated_at) FROM stdin;
\.


--
-- Data for Name: user_certifications; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.user_certifications (user_id, certification_id, issued_at) FROM stdin;
\.


--
-- Data for Name: user_entitlements; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.user_entitlements (id, enrollment_id, user_id, source_type, source_id, total_sessions, used_sessions, expires_at, anonymized_at, created_at, updated_at) FROM stdin;
a9de4846-b6ab-4a0e-beb4-efec1a27d1cf	3c3c0ead-befa-475f-819a-574a45e39543	3c09cb50-3d99-4ff8-a57f-44fd4135d070	package	11111111-1111-1111-1111-111111111101	6	0	2027-07-22 16:42:46.981945+07	\N	2026-07-22 16:42:46.981945+07	2026-07-22 16:42:46.981945+07
\.


--
-- Name: driving_sessions_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.driving_sessions_id_seq', 34, true);


--
-- Name: payments_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.payments_id_seq', 1, false);


--
-- Name: schedules_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.schedules_id_seq', 1983, true);


--
-- Name: driving_sessions driving_sessions_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.driving_sessions
    ADD CONSTRAINT driving_sessions_pkey PRIMARY KEY (id);


--
-- Name: enrollments enrollments_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.enrollments
    ADD CONSTRAINT enrollments_pkey PRIMARY KEY (id);


--
-- Name: payments payments_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.payments
    ADD CONSTRAINT payments_pkey PRIMARY KEY (id);


--
-- Name: schedules schedules_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.schedules
    ADD CONSTRAINT schedules_pkey PRIMARY KEY (id);


--
-- Name: transaction_items transaction_items_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.transaction_items
    ADD CONSTRAINT transaction_items_pkey PRIMARY KEY (id);


--
-- Name: transactions transactions_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.transactions
    ADD CONSTRAINT transactions_pkey PRIMARY KEY (id);


--
-- Name: user_certifications user_certifications_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_certifications
    ADD CONSTRAINT user_certifications_pkey PRIMARY KEY (user_id, certification_id);


--
-- Name: user_entitlements user_entitlements_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_entitlements
    ADD CONSTRAINT user_entitlements_pkey PRIMARY KEY (id);


--
-- Name: idx_driving_sessions_anonymized_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_driving_sessions_anonymized_at ON public.driving_sessions USING btree (anonymized_at);


--
-- Name: idx_driving_sessions_car_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_driving_sessions_car_id ON public.driving_sessions USING btree (car_id);


--
-- Name: idx_driving_sessions_enrollment_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_driving_sessions_enrollment_id ON public.driving_sessions USING btree (enrollment_id);


--
-- Name: idx_driving_sessions_entitlement_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_driving_sessions_entitlement_id ON public.driving_sessions USING btree (entitlement_id);


--
-- Name: idx_driving_sessions_instructor_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_driving_sessions_instructor_id ON public.driving_sessions USING btree (instructor_id);


--
-- Name: idx_driving_sessions_schedule_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_driving_sessions_schedule_id ON public.driving_sessions USING btree (schedule_id);


--
-- Name: idx_driving_sessions_user_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_driving_sessions_user_id ON public.driving_sessions USING btree (user_id);


--
-- Name: idx_enrollments_anonymized_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_enrollments_anonymized_at ON public.enrollments USING btree (anonymized_at);


--
-- Name: idx_enrollments_package_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_enrollments_package_id ON public.enrollments USING btree (package_id);


--
-- Name: idx_enrollments_user_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_enrollments_user_id ON public.enrollments USING btree (user_id);


--
-- Name: idx_payments_enrollment_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_payments_enrollment_id ON public.payments USING btree (enrollment_id);


--
-- Name: idx_payments_order_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX idx_payments_order_id ON public.payments USING btree (order_id);


--
-- Name: idx_payments_user_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_payments_user_id ON public.payments USING btree (user_id);


--
-- Name: idx_schedules_car_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_schedules_car_id ON public.schedules USING btree (car_id);


--
-- Name: idx_schedules_date; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_schedules_date ON public.schedules USING btree (date);


--
-- Name: idx_schedules_enrollment_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_schedules_enrollment_id ON public.schedules USING btree (enrollment_id);


--
-- Name: idx_schedules_instructor_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_schedules_instructor_id ON public.schedules USING btree (instructor_id);


--
-- Name: idx_schedules_user_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_schedules_user_id ON public.schedules USING btree (user_id);


--
-- Name: idx_transaction_items_transaction_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_transaction_items_transaction_id ON public.transaction_items USING btree (transaction_id);


--
-- Name: idx_transactions_enrollment_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX idx_transactions_enrollment_id ON public.transactions USING btree (enrollment_id);


--
-- Name: idx_transactions_user_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_transactions_user_id ON public.transactions USING btree (user_id);


--
-- Name: idx_user_entitlements_anonymized_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_user_entitlements_anonymized_at ON public.user_entitlements USING btree (anonymized_at);


--
-- Name: idx_user_entitlements_enrollment_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_user_entitlements_enrollment_id ON public.user_entitlements USING btree (enrollment_id);


--
-- Name: idx_user_entitlements_user_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_user_entitlements_user_id ON public.user_entitlements USING btree (user_id);


--
-- Name: user_entitlements fk_enrollments_entitlements; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_entitlements
    ADD CONSTRAINT fk_enrollments_entitlements FOREIGN KEY (enrollment_id) REFERENCES public.enrollments(id);


--
-- Name: transaction_items fk_transactions_items; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.transaction_items
    ADD CONSTRAINT fk_transactions_items FOREIGN KEY (transaction_id) REFERENCES public.transactions(id);


--
-- PostgreSQL database dump complete
--

\unrestrict 1SsoONKiZieIBhMK0TbdHBfiGumUhrTqgFbKlIca2EuIq6KmwfZmaYjU2lB8fWd

