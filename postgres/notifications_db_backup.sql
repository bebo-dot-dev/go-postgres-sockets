--
-- PostgreSQL database dump
--

-- Dumped from database version 13.3 (Ubuntu 13.3-1.pgdg20.04+1)
-- Dumped by pg_dump version 13.3 (Ubuntu 13.3-1.pgdg20.04+1)

-- Started on 2021-07-31 17:06:37 BST

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

DROP DATABASE notifications;
--
-- TOC entry 3014 (class 1262 OID 26499)
-- Name: notifications; Type: DATABASE; Schema: -; Owner: postgres
--

CREATE DATABASE notifications WITH TEMPLATE = template0 ENCODING = 'UTF8' LOCALE = 'en_GB.UTF-8';


ALTER DATABASE notifications OWNER TO postgres;

\connect notifications

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

--
-- TOC entry 204 (class 1255 OID 26501)
-- Name: new_notification(integer, character varying); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.new_notification(p_notification_type_id integer, p_notification_text character varying) RETURNS integer
    LANGUAGE sql
    AS $$
	INSERT INTO public.notifications
	(
		notification_type_id, 
		notification_text
	)
	SELECT
		p_notification_type_id,
		p_notification_text
	RETURNING id;
$$;


ALTER FUNCTION public.new_notification(p_notification_type_id integer, p_notification_text character varying) OWNER TO postgres;

--
-- TOC entry 205 (class 1255 OID 26502)
-- Name: notify_notification_changes(); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.notify_notification_changes() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
DECLARE 
    tr_record RECORD;
    trigger_data json;
    notification_data json;
BEGIN    
    IF (TG_OP = 'DELETE') THEN
        tr_record = OLD;        
    ELSE
        tr_record = NEW;        
    END IF;
    
    trigger_data = row_to_json(r)
        FROM (
            SELECT 
                n.*,
                nt.name notification_type
            FROM
                (SELECT tr_record.*) n
                JOIN public.notification_type nt on n.notification_type_id = nt.id
        ) r;
    
    notification_data = json_build_object(
        'table', TG_TABLE_NAME,
        'operation', TG_OP,
        'data', trigger_data
	  );

    PERFORM pg_notify(
        'notifications_data_changed',
        notification_data::text
    );

  RETURN tr_record;
END;
$$;


ALTER FUNCTION public.notify_notification_changes() OWNER TO postgres;

--
-- TOC entry 200 (class 1259 OID 26503)
-- Name: notification_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.notification_id_seq
    START WITH 0
    INCREMENT BY 1
    MINVALUE 0
    MAXVALUE 2147483647
    CACHE 1;


ALTER TABLE public.notification_id_seq OWNER TO postgres;

--
-- TOC entry 201 (class 1259 OID 26505)
-- Name: notification_type_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.notification_type_id_seq
    START WITH 0
    INCREMENT BY 1
    MINVALUE 0
    MAXVALUE 2147483647
    CACHE 1;


ALTER TABLE public.notification_type_id_seq OWNER TO postgres;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- TOC entry 202 (class 1259 OID 26507)
-- Name: notification_type; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.notification_type (
    id integer DEFAULT nextval('public.notification_type_id_seq'::regclass) NOT NULL,
    name text NOT NULL
);


ALTER TABLE public.notification_type OWNER TO postgres;

--
-- TOC entry 203 (class 1259 OID 26514)
-- Name: notifications; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.notifications (
    id integer DEFAULT nextval('public.notification_id_seq'::regclass) NOT NULL,
    notification_type_id integer NOT NULL,
    created_timestamp timestamp without time zone DEFAULT timezone('utc'::text, now()) NOT NULL,
    notification_text text NOT NULL
);


ALTER TABLE public.notifications OWNER TO postgres;

--
-- TOC entry 3007 (class 0 OID 26507)
-- Dependencies: 202
-- Data for Name: notification_type; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.notification_type (id, name) FROM stdin;
0	None
1	Email
2	SMS
3	Slack
\.


--
-- TOC entry 3008 (class 0 OID 26514)
-- Dependencies: 203
-- Data for Name: notifications; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.notifications (id, notification_type_id, created_timestamp, notification_text) FROM stdin;
\.


--
-- TOC entry 3015 (class 0 OID 0)
-- Dependencies: 200
-- Name: notification_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.notification_id_seq', 5, true);


--
-- TOC entry 3016 (class 0 OID 0)
-- Dependencies: 201
-- Name: notification_type_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.notification_type_id_seq', 0, false);


--
-- TOC entry 2870 (class 2606 OID 26523)
-- Name: notification_type notification_type_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.notification_type
    ADD CONSTRAINT notification_type_pkey PRIMARY KEY (id);


--
-- TOC entry 2872 (class 2606 OID 26525)
-- Name: notifications notifications_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.notifications
    ADD CONSTRAINT notifications_pkey PRIMARY KEY (id);


--
-- TOC entry 2874 (class 2620 OID 26526)
-- Name: notifications tr_notifications; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER tr_notifications AFTER INSERT OR DELETE OR UPDATE ON public.notifications FOR EACH ROW EXECUTE FUNCTION public.notify_notification_changes();


--
-- TOC entry 2873 (class 2606 OID 26527)
-- Name: notifications notification_type_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.notifications
    ADD CONSTRAINT notification_type_id_fkey FOREIGN KEY (notification_type_id) REFERENCES public.notification_type(id);


-- Completed on 2021-07-31 17:06:38 BST

--
-- PostgreSQL database dump complete
--

