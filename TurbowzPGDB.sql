-- DB structure
-- PostgreSQL database dump
--

\restrict vZtAHfgjjXuyUTYkHdhBI9mSRcWasE2CdGANuP3AogeWTQTx8fae63g3DhR0Ed3

-- Dumped from database version 17.6
-- Dumped by pg_dump version 18.0

-- Started on 2026-03-28 11:04:45 IST

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- TOC entry 2 (class 3079 OID 16602)
-- Name: vector; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS vector WITH SCHEMA public;


--
-- TOC entry 4250 (class 0 OID 0)
-- Dependencies: 2
-- Name: EXTENSION vector; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION vector IS 'vector data type and ivfflat and hnsw access methods';


--
-- TOC entry 372 (class 1255 OID 16493)
-- Name: update_updated_at_column(); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.update_updated_at_column() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
   NEW.updated_at = NOW();
   RETURN NEW;
END;
$$;


ALTER FUNCTION public.update_updated_at_column() OWNER TO postgres;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- TOC entry 245 (class 1259 OID 25318)
-- Name: banner_ads; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.banner_ads (
    ad_id bigint NOT NULL,
    title text NOT NULL,
    redirect_url text NOT NULL,
    start_date timestamp with time zone NOT NULL,
    end_date timestamp with time zone NOT NULL,
    views bigint DEFAULT 0,
    clicks bigint DEFAULT 0,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    uploader_id bigint NOT NULL
);


ALTER TABLE public.banner_ads OWNER TO postgres;

--
-- TOC entry 244 (class 1259 OID 25317)
-- Name: banner_ads_ad_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.banner_ads_ad_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.banner_ads_ad_id_seq OWNER TO postgres;

--
-- TOC entry 4251 (class 0 OID 0)
-- Dependencies: 244
-- Name: banner_ads_ad_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.banner_ads_ad_id_seq OWNED BY public.banner_ads.ad_id;


--
-- TOC entry 231 (class 1259 OID 16978)
-- Name: channels; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.channels (
    channel_id bigint NOT NULL,
    creator_id bigint NOT NULL,
    channel_name text NOT NULL,
    embedding public.vector(512),
    channel_url text NOT NULL
);


ALTER TABLE public.channels OWNER TO postgres;

--
-- TOC entry 230 (class 1259 OID 16977)
-- Name: channels_channel_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.channels_channel_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.channels_channel_id_seq OWNER TO postgres;

--
-- TOC entry 4252 (class 0 OID 0)
-- Dependencies: 230
-- Name: channels_channel_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.channels_channel_id_seq OWNED BY public.channels.channel_id;


--
-- TOC entry 224 (class 1259 OID 16474)
-- Name: comments_table; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.comments_table (
    comment_id bigint NOT NULL,
    commenter_id bigint NOT NULL,
    parent_video_id bigint NOT NULL,
    comment_text text NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    parent_comment_id bigint
);


ALTER TABLE public.comments_table OWNER TO postgres;

--
-- TOC entry 225 (class 1259 OID 16497)
-- Name: comments_table_comment_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

ALTER TABLE public.comments_table ALTER COLUMN comment_id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.comments_table_comment_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- TOC entry 227 (class 1259 OID 16539)
-- Name: connections; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.connections (
    user1_id bigint NOT NULL,
    user2_id bigint NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    CONSTRAINT connections_check CHECK ((user1_id < user2_id))
);


ALTER TABLE public.connections OWNER TO postgres;

--
-- TOC entry 239 (class 1259 OID 25259)
-- Name: eco_comments_table; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.eco_comments_table (
    comment_id bigint NOT NULL,
    parent_eco_id bigint NOT NULL,
    commenter_id bigint NOT NULL,
    comment_text text NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.eco_comments_table OWNER TO postgres;

--
-- TOC entry 238 (class 1259 OID 25258)
-- Name: eco_comments_table_comment_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.eco_comments_table_comment_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.eco_comments_table_comment_id_seq OWNER TO postgres;

--
-- TOC entry 4253 (class 0 OID 0)
-- Dependencies: 238
-- Name: eco_comments_table_comment_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.eco_comments_table_comment_id_seq OWNED BY public.eco_comments_table.comment_id;


--
-- TOC entry 234 (class 1259 OID 17040)
-- Name: eco_data; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.eco_data (
    eco_id bigint NOT NULL,
    uploader_id bigint NOT NULL,
    eco_text text NOT NULL,
    eco_url text NOT NULL,
    images_count integer DEFAULT 0 NOT NULL,
    tags text[],
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    luv_count bigint DEFAULT 0 NOT NULL,
    comment_count bigint DEFAULT 0 NOT NULL,
    embeddings public.vector(512),
    share_count bigint DEFAULT 0 NOT NULL,
    view_count bigint DEFAULT 0 NOT NULL,
    saves_count bigint DEFAULT 0 NOT NULL,
    last_trending_score double precision DEFAULT 0,
    trending_delta double precision DEFAULT 0,
    CONSTRAINT post_data_post_text_check CHECK ((length(TRIM(BOTH FROM eco_text)) > 0))
);


ALTER TABLE public.eco_data OWNER TO postgres;

--
-- TOC entry 237 (class 1259 OID 17086)
-- Name: eco_luv_events; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.eco_luv_events (
    eco_id bigint NOT NULL,
    user_id bigint NOT NULL,
    event_time timestamp with time zone DEFAULT now()
);


ALTER TABLE public.eco_luv_events OWNER TO postgres;

--
-- TOC entry 249 (class 1259 OID 33538)
-- Name: eco_votes; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.eco_votes (
    id bigint NOT NULL,
    eco_id bigint NOT NULL,
    user_id bigint NOT NULL,
    quality smallint NOT NULL,
    ai_usage smallint NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT eco_votes_ai_usage_check CHECK (((ai_usage >= 1) AND (ai_usage <= 5))),
    CONSTRAINT eco_votes_quality_check CHECK (((quality >= 1) AND (quality <= 5)))
);


ALTER TABLE public.eco_votes OWNER TO postgres;

--
-- TOC entry 248 (class 1259 OID 33537)
-- Name: eco_votes_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.eco_votes_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.eco_votes_id_seq OWNER TO postgres;

--
-- TOC entry 4254 (class 0 OID 0)
-- Dependencies: 248
-- Name: eco_votes_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.eco_votes_id_seq OWNED BY public.eco_votes.id;


--
-- TOC entry 253 (class 1259 OID 33646)
-- Name: event_data; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.event_data (
    event_id bigint NOT NULL,
    event_url text NOT NULL,
    uploader_id bigint NOT NULL,
    event_description text,
    view_count bigint DEFAULT 0 NOT NULL,
    luv_count bigint DEFAULT 0 NOT NULL,
    comment_count bigint DEFAULT 0 NOT NULL,
    saves_count bigint DEFAULT 0 NOT NULL,
    images_count integer DEFAULT 0 NOT NULL,
    tags text[] DEFAULT '{}'::text[] NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    event_start_time timestamp with time zone NOT NULL,
    event_end_time timestamp with time zone NOT NULL,
    event_title text NOT NULL,
    last_trending_score double precision DEFAULT 0,
    trending_delta double precision DEFAULT 0,
    CONSTRAINT chk_event_time CHECK ((event_end_time > event_start_time)),
    CONSTRAINT events_event_title_not_empty CHECK ((TRIM(BOTH FROM event_title) <> ''::text))
);


ALTER TABLE public.event_data OWNER TO postgres;

--
-- TOC entry 252 (class 1259 OID 33645)
-- Name: event_data_event_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.event_data_event_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.event_data_event_id_seq OWNER TO postgres;

--
-- TOC entry 4255 (class 0 OID 0)
-- Dependencies: 252
-- Name: event_data_event_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.event_data_event_id_seq OWNED BY public.event_data.event_id;


--
-- TOC entry 259 (class 1259 OID 33750)
-- Name: event_luvs; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.event_luvs (
    luv_id bigint NOT NULL,
    event_id bigint NOT NULL,
    user_id bigint NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.event_luvs OWNER TO postgres;

--
-- TOC entry 258 (class 1259 OID 33749)
-- Name: event_luvs_luv_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.event_luvs_luv_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.event_luvs_luv_id_seq OWNER TO postgres;

--
-- TOC entry 4256 (class 0 OID 0)
-- Dependencies: 258
-- Name: event_luvs_luv_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.event_luvs_luv_id_seq OWNED BY public.event_luvs.luv_id;


--
-- TOC entry 255 (class 1259 OID 33685)
-- Name: event_participants; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.event_participants (
    participant_id bigint NOT NULL,
    event_id bigint NOT NULL,
    user_id bigint NOT NULL,
    status smallint DEFAULT 1 NOT NULL,
    registered_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.event_participants OWNER TO postgres;

--
-- TOC entry 254 (class 1259 OID 33684)
-- Name: event_participants_participant_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.event_participants_participant_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.event_participants_participant_id_seq OWNER TO postgres;

--
-- TOC entry 4257 (class 0 OID 0)
-- Dependencies: 254
-- Name: event_participants_participant_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.event_participants_participant_id_seq OWNED BY public.event_participants.participant_id;


--
-- TOC entry 257 (class 1259 OID 33720)
-- Name: event_saves; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.event_saves (
    save_id bigint NOT NULL,
    event_id bigint NOT NULL,
    user_id bigint NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.event_saves OWNER TO postgres;

--
-- TOC entry 256 (class 1259 OID 33719)
-- Name: event_saves_save_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.event_saves_save_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.event_saves_save_id_seq OWNER TO postgres;

--
-- TOC entry 4258 (class 0 OID 0)
-- Dependencies: 256
-- Name: event_saves_save_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.event_saves_save_id_seq OWNED BY public.event_saves.save_id;


--
-- TOC entry 226 (class 1259 OID 16498)
-- Name: follow_table; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.follow_table (
    follower_id bigint NOT NULL,
    followee_id bigint NOT NULL,
    followed_at timestamp without time zone DEFAULT now()
);


ALTER TABLE public.follow_table OWNER TO postgres;

--
-- TOC entry 236 (class 1259 OID 17064)
-- Name: video_history_table; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.video_history_table (
    history_id bigint NOT NULL,
    watcher_id bigint NOT NULL,
    video_id bigint NOT NULL,
    watched_at timestamp with time zone DEFAULT now(),
    watch_time interval DEFAULT '00:00:00'::interval
);


ALTER TABLE public.video_history_table OWNER TO postgres;

--
-- TOC entry 235 (class 1259 OID 17063)
-- Name: history_table_history_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.history_table_history_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.history_table_history_id_seq OWNER TO postgres;

--
-- TOC entry 4259 (class 0 OID 0)
-- Dependencies: 235
-- Name: history_table_history_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.history_table_history_id_seq OWNED BY public.video_history_table.history_id;


--
-- TOC entry 233 (class 1259 OID 17039)
-- Name: post_data_post_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.post_data_post_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.post_data_post_id_seq OWNER TO postgres;

--
-- TOC entry 4260 (class 0 OID 0)
-- Dependencies: 233
-- Name: post_data_post_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.post_data_post_id_seq OWNED BY public.eco_data.eco_id;


--
-- TOC entry 243 (class 1259 OID 25300)
-- Name: saved_ecos; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.saved_ecos (
    id bigint NOT NULL,
    user_id bigint NOT NULL,
    eco_id bigint NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.saved_ecos OWNER TO postgres;

--
-- TOC entry 242 (class 1259 OID 25299)
-- Name: saved_eco_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.saved_eco_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.saved_eco_id_seq OWNER TO postgres;

--
-- TOC entry 4261 (class 0 OID 0)
-- Dependencies: 242
-- Name: saved_eco_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.saved_eco_id_seq OWNED BY public.saved_ecos.id;


--
-- TOC entry 241 (class 1259 OID 25282)
-- Name: saved_videos; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.saved_videos (
    id bigint NOT NULL,
    user_id bigint NOT NULL,
    video_id bigint NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.saved_videos OWNER TO postgres;

--
-- TOC entry 240 (class 1259 OID 25281)
-- Name: saved_videos_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.saved_videos_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.saved_videos_id_seq OWNER TO postgres;

--
-- TOC entry 4262 (class 0 OID 0)
-- Dependencies: 240
-- Name: saved_videos_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.saved_videos_id_seq OWNED BY public.saved_videos.id;


--
-- TOC entry 247 (class 1259 OID 33503)
-- Name: turbomax_status_table; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.turbomax_status_table (
    turbomax_id bigint NOT NULL,
    user_id bigint,
    start_date timestamp with time zone NOT NULL,
    expiry_date timestamp with time zone NOT NULL,
    is_active boolean DEFAULT true NOT NULL,
    payment_amount numeric(10,2) DEFAULT 0.00 NOT NULL,
    verification_complete boolean DEFAULT false NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.turbomax_status_table OWNER TO postgres;

--
-- TOC entry 246 (class 1259 OID 33502)
-- Name: turbomax_status_table_turbo_max_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.turbomax_status_table_turbo_max_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.turbomax_status_table_turbo_max_id_seq OWNER TO postgres;

--
-- TOC entry 4263 (class 0 OID 0)
-- Dependencies: 246
-- Name: turbomax_status_table_turbo_max_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.turbomax_status_table_turbo_max_id_seq OWNED BY public.turbomax_status_table.turbomax_id;


--
-- TOC entry 222 (class 1259 OID 16433)
-- Name: user_authentication; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.user_authentication (
    auth_id bigint NOT NULL,
    user_id bigint NOT NULL,
    user_login_account text NOT NULL,
    user_phone_number text,
    user_hashed_password text NOT NULL,
    account_created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.user_authentication OWNER TO postgres;

--
-- TOC entry 221 (class 1259 OID 16432)
-- Name: user_authentication_auth_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.user_authentication_auth_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.user_authentication_auth_id_seq OWNER TO postgres;

--
-- TOC entry 4264 (class 0 OID 0)
-- Dependencies: 221
-- Name: user_authentication_auth_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.user_authentication_auth_id_seq OWNED BY public.user_authentication.auth_id;


--
-- TOC entry 220 (class 1259 OID 16421)
-- Name: user_data_table; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.user_data_table (
    user_id bigint NOT NULL,
    user_handle text NOT NULL,
    user_profile_name text NOT NULL,
    user_description text,
    from_location text,
    user_date_of_birth date,
    gender character varying(20),
    embeddings public.vector(512),
    eco_embeddings public.vector(512),
    url text NOT NULL
);


ALTER TABLE public.user_data_table OWNER TO postgres;

--
-- TOC entry 219 (class 1259 OID 16420)
-- Name: user_data_table_user_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.user_data_table_user_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.user_data_table_user_id_seq OWNER TO postgres;

--
-- TOC entry 4265 (class 0 OID 0)
-- Dependencies: 219
-- Name: user_data_table_user_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.user_data_table_user_id_seq OWNED BY public.user_data_table.user_id;


--
-- TOC entry 229 (class 1259 OID 16956)
-- Name: user_favourite_topics; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.user_favourite_topics (
    user_id integer NOT NULL,
    topic character varying(255) NOT NULL
);


ALTER TABLE public.user_favourite_topics OWNER TO postgres;

--
-- TOC entry 218 (class 1259 OID 16402)
-- Name: video_data; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.video_data (
    video_id bigint NOT NULL,
    uploader_id bigint NOT NULL,
    title text NOT NULL,
    video_info text,
    upload_time timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    views bigint DEFAULT 0,
    luv bigint DEFAULT 0,
    report bigint DEFAULT 0,
    video_url text NOT NULL,
    transcript text,
    embeddings public.vector(512),
    last_trending_score double precision DEFAULT 0,
    comments bigint DEFAULT 0,
    shares bigint DEFAULT 0,
    trending_delta double precision DEFAULT 0,
    last_trending_updated_at timestamp without time zone DEFAULT now(),
    tags text[],
    saves_count integer DEFAULT 0 NOT NULL
);


ALTER TABLE public.video_data OWNER TO postgres;

--
-- TOC entry 223 (class 1259 OID 16466)
-- Name: video_data_video_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

ALTER TABLE public.video_data ALTER COLUMN video_id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.video_data_video_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- TOC entry 232 (class 1259 OID 17023)
-- Name: video_luv_events; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.video_luv_events (
    video_id bigint NOT NULL,
    user_id bigint NOT NULL,
    luv_time timestamp with time zone DEFAULT now()
);


ALTER TABLE public.video_luv_events OWNER TO postgres;

--
-- TOC entry 228 (class 1259 OID 16946)
-- Name: video_tags; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.video_tags (
    video_id integer NOT NULL,
    tag character varying(255) NOT NULL
);


ALTER TABLE public.video_tags OWNER TO postgres;

--
-- TOC entry 251 (class 1259 OID 33603)
-- Name: video_votes; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.video_votes (
    id bigint NOT NULL,
    video_id bigint NOT NULL,
    user_id bigint NOT NULL,
    quality smallint NOT NULL,
    ai_usage smallint NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT video_votes_ai_usage_check CHECK (((ai_usage >= 1) AND (ai_usage <= 5))),
    CONSTRAINT video_votes_quality_check CHECK (((quality >= 1) AND (quality <= 5)))
);


ALTER TABLE public.video_votes OWNER TO postgres;

--
-- TOC entry 250 (class 1259 OID 33602)
-- Name: video_votes_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.video_votes_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.video_votes_id_seq OWNER TO postgres;

--
-- TOC entry 4266 (class 0 OID 0)
-- Dependencies: 250
-- Name: video_votes_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.video_votes_id_seq OWNED BY public.video_votes.id;


--
-- TOC entry 3945 (class 2604 OID 25321)
-- Name: banner_ads ad_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.banner_ads ALTER COLUMN ad_id SET DEFAULT nextval('public.banner_ads_ad_id_seq'::regclass);


--
-- TOC entry 3922 (class 2604 OID 16981)
-- Name: channels channel_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.channels ALTER COLUMN channel_id SET DEFAULT nextval('public.channels_channel_id_seq'::regclass);


--
-- TOC entry 3938 (class 2604 OID 25262)
-- Name: eco_comments_table comment_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.eco_comments_table ALTER COLUMN comment_id SET DEFAULT nextval('public.eco_comments_table_comment_id_seq'::regclass);


--
-- TOC entry 3924 (class 2604 OID 17043)
-- Name: eco_data eco_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.eco_data ALTER COLUMN eco_id SET DEFAULT nextval('public.post_data_post_id_seq'::regclass);


--
-- TOC entry 3955 (class 2604 OID 33541)
-- Name: eco_votes id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.eco_votes ALTER COLUMN id SET DEFAULT nextval('public.eco_votes_id_seq'::regclass);


--
-- TOC entry 3961 (class 2604 OID 33649)
-- Name: event_data event_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.event_data ALTER COLUMN event_id SET DEFAULT nextval('public.event_data_event_id_seq'::regclass);


--
-- TOC entry 3977 (class 2604 OID 33753)
-- Name: event_luvs luv_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.event_luvs ALTER COLUMN luv_id SET DEFAULT nextval('public.event_luvs_luv_id_seq'::regclass);


--
-- TOC entry 3971 (class 2604 OID 33688)
-- Name: event_participants participant_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.event_participants ALTER COLUMN participant_id SET DEFAULT nextval('public.event_participants_participant_id_seq'::regclass);


--
-- TOC entry 3975 (class 2604 OID 33723)
-- Name: event_saves save_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.event_saves ALTER COLUMN save_id SET DEFAULT nextval('public.event_saves_save_id_seq'::regclass);


--
-- TOC entry 3943 (class 2604 OID 25303)
-- Name: saved_ecos id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.saved_ecos ALTER COLUMN id SET DEFAULT nextval('public.saved_eco_id_seq'::regclass);


--
-- TOC entry 3941 (class 2604 OID 25285)
-- Name: saved_videos id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.saved_videos ALTER COLUMN id SET DEFAULT nextval('public.saved_videos_id_seq'::regclass);


--
-- TOC entry 3949 (class 2604 OID 33506)
-- Name: turbomax_status_table turbomax_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.turbomax_status_table ALTER COLUMN turbomax_id SET DEFAULT nextval('public.turbomax_status_table_turbo_max_id_seq'::regclass);


--
-- TOC entry 3916 (class 2604 OID 16436)
-- Name: user_authentication auth_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_authentication ALTER COLUMN auth_id SET DEFAULT nextval('public.user_authentication_auth_id_seq'::regclass);


--
-- TOC entry 3915 (class 2604 OID 16424)
-- Name: user_data_table user_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_data_table ALTER COLUMN user_id SET DEFAULT nextval('public.user_data_table_user_id_seq'::regclass);


--
-- TOC entry 3934 (class 2604 OID 17067)
-- Name: video_history_table history_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.video_history_table ALTER COLUMN history_id SET DEFAULT nextval('public.history_table_history_id_seq'::regclass);


--
-- TOC entry 3958 (class 2604 OID 33606)
-- Name: video_votes id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.video_votes ALTER COLUMN id SET DEFAULT nextval('public.video_votes_id_seq'::regclass);


--
-- TOC entry 3988 (class 2606 OID 16458)
-- Name: video_data MetaDataTable_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.video_data
    ADD CONSTRAINT "MetaDataTable_pkey" PRIMARY KEY (video_id);


--
-- TOC entry 4030 (class 2606 OID 25328)
-- Name: banner_ads banner_ads_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.banner_ads
    ADD CONSTRAINT banner_ads_pkey PRIMARY KEY (ad_id);


--
-- TOC entry 4014 (class 2606 OID 16985)
-- Name: channels channels_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.channels
    ADD CONSTRAINT channels_pkey PRIMARY KEY (channel_id);


--
-- TOC entry 4004 (class 2606 OID 16496)
-- Name: comments_table comments_table_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.comments_table
    ADD CONSTRAINT comments_table_pkey PRIMARY KEY (comment_id);


--
-- TOC entry 4008 (class 2606 OID 16544)
-- Name: connections connections_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.connections
    ADD CONSTRAINT connections_pkey PRIMARY KEY (user1_id, user2_id);


--
-- TOC entry 4024 (class 2606 OID 25268)
-- Name: eco_comments_table eco_comments_table_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.eco_comments_table
    ADD CONSTRAINT eco_comments_table_pkey PRIMARY KEY (comment_id);


--
-- TOC entry 4022 (class 2606 OID 17091)
-- Name: eco_luv_events eco_luv_events_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.eco_luv_events
    ADD CONSTRAINT eco_luv_events_pkey PRIMARY KEY (eco_id, user_id);


--
-- TOC entry 4034 (class 2606 OID 33548)
-- Name: eco_votes eco_votes_eco_id_user_id_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.eco_votes
    ADD CONSTRAINT eco_votes_eco_id_user_id_key UNIQUE (eco_id, user_id);


--
-- TOC entry 4036 (class 2606 OID 33546)
-- Name: eco_votes eco_votes_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.eco_votes
    ADD CONSTRAINT eco_votes_pkey PRIMARY KEY (id);


--
-- TOC entry 4042 (class 2606 OID 33663)
-- Name: event_data event_data_event_url_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.event_data
    ADD CONSTRAINT event_data_event_url_key UNIQUE (event_url);


--
-- TOC entry 4044 (class 2606 OID 33661)
-- Name: event_data event_data_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.event_data
    ADD CONSTRAINT event_data_pkey PRIMARY KEY (event_id);


--
-- TOC entry 4057 (class 2606 OID 33756)
-- Name: event_luvs event_luvs_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.event_luvs
    ADD CONSTRAINT event_luvs_pkey PRIMARY KEY (luv_id);


--
-- TOC entry 4046 (class 2606 OID 33693)
-- Name: event_participants event_participants_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.event_participants
    ADD CONSTRAINT event_participants_pkey PRIMARY KEY (participant_id);


--
-- TOC entry 4053 (class 2606 OID 33726)
-- Name: event_saves event_saves_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.event_saves
    ADD CONSTRAINT event_saves_pkey PRIMARY KEY (save_id);


--
-- TOC entry 4006 (class 2606 OID 16503)
-- Name: follow_table follow_table_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.follow_table
    ADD CONSTRAINT follow_table_pkey PRIMARY KEY (follower_id, followee_id);


--
-- TOC entry 4020 (class 2606 OID 17071)
-- Name: video_history_table history_table_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.video_history_table
    ADD CONSTRAINT history_table_pkey PRIMARY KEY (history_id);


--
-- TOC entry 4018 (class 2606 OID 17051)
-- Name: eco_data post_data_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.eco_data
    ADD CONSTRAINT post_data_pkey PRIMARY KEY (eco_id);


--
-- TOC entry 4028 (class 2606 OID 25306)
-- Name: saved_ecos saved_eco_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.saved_ecos
    ADD CONSTRAINT saved_eco_pkey PRIMARY KEY (id);


--
-- TOC entry 4026 (class 2606 OID 25288)
-- Name: saved_videos saved_videos_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.saved_videos
    ADD CONSTRAINT saved_videos_pkey PRIMARY KEY (id);


--
-- TOC entry 4032 (class 2606 OID 33513)
-- Name: turbomax_status_table turbomax_status_table_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.turbomax_status_table
    ADD CONSTRAINT turbomax_status_table_pkey PRIMARY KEY (turbomax_id);


--
-- TOC entry 4059 (class 2606 OID 33758)
-- Name: event_luvs uq_event_luv; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.event_luvs
    ADD CONSTRAINT uq_event_luv UNIQUE (event_id, user_id);


--
-- TOC entry 4055 (class 2606 OID 33728)
-- Name: event_saves uq_event_save; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.event_saves
    ADD CONSTRAINT uq_event_save UNIQUE (event_id, user_id);


--
-- TOC entry 4051 (class 2606 OID 33695)
-- Name: event_participants uq_event_user; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.event_participants
    ADD CONSTRAINT uq_event_user UNIQUE (event_id, user_id);


--
-- TOC entry 3998 (class 2606 OID 16441)
-- Name: user_authentication user_authentication_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_authentication
    ADD CONSTRAINT user_authentication_pkey PRIMARY KEY (auth_id);


--
-- TOC entry 4000 (class 2606 OID 16443)
-- Name: user_authentication user_authentication_user_login_account_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_authentication
    ADD CONSTRAINT user_authentication_user_login_account_key UNIQUE (user_login_account);


--
-- TOC entry 4002 (class 2606 OID 16445)
-- Name: user_authentication user_authentication_user_phone_number_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_authentication
    ADD CONSTRAINT user_authentication_user_phone_number_key UNIQUE (user_phone_number);


--
-- TOC entry 3992 (class 2606 OID 16429)
-- Name: user_data_table user_data_table_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_data_table
    ADD CONSTRAINT user_data_table_pkey PRIMARY KEY (user_id);


--
-- TOC entry 3994 (class 2606 OID 25256)
-- Name: user_data_table user_data_table_url_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_data_table
    ADD CONSTRAINT user_data_table_url_key UNIQUE (url);


--
-- TOC entry 3996 (class 2606 OID 16431)
-- Name: user_data_table user_data_table_user_handle_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_data_table
    ADD CONSTRAINT user_data_table_user_handle_key UNIQUE (user_handle);


--
-- TOC entry 4012 (class 2606 OID 16960)
-- Name: user_favourite_topics user_favourite_topics_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_favourite_topics
    ADD CONSTRAINT user_favourite_topics_pkey PRIMARY KEY (user_id, topic);


--
-- TOC entry 3990 (class 2606 OID 16468)
-- Name: video_data video_data_video_url_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.video_data
    ADD CONSTRAINT video_data_video_url_key UNIQUE (video_url);


--
-- TOC entry 4016 (class 2606 OID 17028)
-- Name: video_luv_events video_luvs_events_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.video_luv_events
    ADD CONSTRAINT video_luvs_events_pkey PRIMARY KEY (video_id, user_id);


--
-- TOC entry 4010 (class 2606 OID 16950)
-- Name: video_tags video_tags_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.video_tags
    ADD CONSTRAINT video_tags_pkey PRIMARY KEY (video_id, tag);


--
-- TOC entry 4038 (class 2606 OID 33612)
-- Name: video_votes video_votes_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.video_votes
    ADD CONSTRAINT video_votes_pkey PRIMARY KEY (id);


--
-- TOC entry 4040 (class 2606 OID 33614)
-- Name: video_votes video_votes_video_id_user_id_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.video_votes
    ADD CONSTRAINT video_votes_video_id_user_id_key UNIQUE (video_id, user_id);


--
-- TOC entry 4047 (class 1259 OID 33706)
-- Name: idx_event_participants_event; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_event_participants_event ON public.event_participants USING btree (event_id);


--
-- TOC entry 4048 (class 1259 OID 33708)
-- Name: idx_event_participants_status; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_event_participants_status ON public.event_participants USING btree (status);


--
-- TOC entry 4049 (class 1259 OID 33707)
-- Name: idx_event_participants_user; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_event_participants_user ON public.event_participants USING btree (user_id);


--
-- TOC entry 4099 (class 2620 OID 16494)
-- Name: comments_table set_timestamp; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER set_timestamp BEFORE UPDATE ON public.comments_table FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- TOC entry 4078 (class 2606 OID 17092)
-- Name: eco_luv_events eco_luv_events_eco_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.eco_luv_events
    ADD CONSTRAINT eco_luv_events_eco_id_fkey FOREIGN KEY (eco_id) REFERENCES public.eco_data(eco_id) ON DELETE CASCADE;


--
-- TOC entry 4079 (class 2606 OID 17097)
-- Name: eco_luv_events eco_luv_events_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.eco_luv_events
    ADD CONSTRAINT eco_luv_events_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.user_data_table(user_id) ON DELETE CASCADE;


--
-- TOC entry 4086 (class 2606 OID 25329)
-- Name: banner_ads fk_ads_uploader; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.banner_ads
    ADD CONSTRAINT fk_ads_uploader FOREIGN KEY (uploader_id) REFERENCES public.user_data_table(user_id);


--
-- TOC entry 4063 (class 2606 OID 16483)
-- Name: comments_table fk_commenter; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.comments_table
    ADD CONSTRAINT fk_commenter FOREIGN KEY (commenter_id) REFERENCES public.user_data_table(user_id) ON DELETE CASCADE;


--
-- TOC entry 4080 (class 2606 OID 25274)
-- Name: eco_comments_table fk_commenter; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.eco_comments_table
    ADD CONSTRAINT fk_commenter FOREIGN KEY (commenter_id) REFERENCES public.user_data_table(user_id) ON DELETE CASCADE;


--
-- TOC entry 4072 (class 2606 OID 16986)
-- Name: channels fk_creator; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.channels
    ADD CONSTRAINT fk_creator FOREIGN KEY (creator_id) REFERENCES public.user_data_table(user_id) ON DELETE CASCADE;


--
-- TOC entry 4088 (class 2606 OID 33549)
-- Name: eco_votes fk_eco; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.eco_votes
    ADD CONSTRAINT fk_eco FOREIGN KEY (eco_id) REFERENCES public.eco_data(eco_id) ON DELETE CASCADE;


--
-- TOC entry 4097 (class 2606 OID 33759)
-- Name: event_luvs fk_event_luvs_event; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.event_luvs
    ADD CONSTRAINT fk_event_luvs_event FOREIGN KEY (event_id) REFERENCES public.event_data(event_id) ON DELETE CASCADE;


--
-- TOC entry 4098 (class 2606 OID 33764)
-- Name: event_luvs fk_event_luvs_user; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.event_luvs
    ADD CONSTRAINT fk_event_luvs_user FOREIGN KEY (user_id) REFERENCES public.user_data_table(user_id) ON DELETE CASCADE;


--
-- TOC entry 4093 (class 2606 OID 33696)
-- Name: event_participants fk_event_participants_event; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.event_participants
    ADD CONSTRAINT fk_event_participants_event FOREIGN KEY (event_id) REFERENCES public.event_data(event_id) ON DELETE CASCADE;


--
-- TOC entry 4094 (class 2606 OID 33701)
-- Name: event_participants fk_event_participants_user; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.event_participants
    ADD CONSTRAINT fk_event_participants_user FOREIGN KEY (user_id) REFERENCES public.user_data_table(user_id) ON DELETE CASCADE;


--
-- TOC entry 4095 (class 2606 OID 33729)
-- Name: event_saves fk_event_saves_event; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.event_saves
    ADD CONSTRAINT fk_event_saves_event FOREIGN KEY (event_id) REFERENCES public.event_data(event_id) ON DELETE CASCADE;


--
-- TOC entry 4096 (class 2606 OID 33734)
-- Name: event_saves fk_event_saves_user; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.event_saves
    ADD CONSTRAINT fk_event_saves_user FOREIGN KEY (user_id) REFERENCES public.user_data_table(user_id) ON DELETE CASCADE;


--
-- TOC entry 4092 (class 2606 OID 33664)
-- Name: event_data fk_event_uploader; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.event_data
    ADD CONSTRAINT fk_event_uploader FOREIGN KEY (uploader_id) REFERENCES public.user_data_table(user_id) ON DELETE CASCADE;


--
-- TOC entry 4064 (class 2606 OID 33520)
-- Name: comments_table fk_parent_comment; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.comments_table
    ADD CONSTRAINT fk_parent_comment FOREIGN KEY (parent_comment_id) REFERENCES public.comments_table(comment_id) ON DELETE CASCADE;


--
-- TOC entry 4081 (class 2606 OID 25269)
-- Name: eco_comments_table fk_parent_eco; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.eco_comments_table
    ADD CONSTRAINT fk_parent_eco FOREIGN KEY (parent_eco_id) REFERENCES public.eco_data(eco_id) ON DELETE CASCADE;


--
-- TOC entry 4084 (class 2606 OID 25312)
-- Name: saved_ecos fk_saved_eco_eco; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.saved_ecos
    ADD CONSTRAINT fk_saved_eco_eco FOREIGN KEY (eco_id) REFERENCES public.eco_data(eco_id) ON DELETE CASCADE;


--
-- TOC entry 4085 (class 2606 OID 25307)
-- Name: saved_ecos fk_saved_eco_user; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.saved_ecos
    ADD CONSTRAINT fk_saved_eco_user FOREIGN KEY (user_id) REFERENCES public.user_data_table(user_id) ON DELETE CASCADE;


--
-- TOC entry 4082 (class 2606 OID 25289)
-- Name: saved_videos fk_saved_videos_user; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.saved_videos
    ADD CONSTRAINT fk_saved_videos_user FOREIGN KEY (user_id) REFERENCES public.user_data_table(user_id) ON DELETE CASCADE;


--
-- TOC entry 4083 (class 2606 OID 25294)
-- Name: saved_videos fk_saved_videos_video; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.saved_videos
    ADD CONSTRAINT fk_saved_videos_video FOREIGN KEY (video_id) REFERENCES public.video_data(video_id) ON DELETE CASCADE;


--
-- TOC entry 4087 (class 2606 OID 33514)
-- Name: turbomax_status_table fk_turbo_user; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.turbomax_status_table
    ADD CONSTRAINT fk_turbo_user FOREIGN KEY (user_id) REFERENCES public.user_data_table(user_id) ON DELETE SET NULL;


--
-- TOC entry 4060 (class 2606 OID 16469)
-- Name: video_data fk_uploader; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.video_data
    ADD CONSTRAINT fk_uploader FOREIGN KEY (uploader_id) REFERENCES public.user_data_table(user_id) ON DELETE SET NULL;


--
-- TOC entry 4075 (class 2606 OID 17052)
-- Name: eco_data fk_uploader; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.eco_data
    ADD CONSTRAINT fk_uploader FOREIGN KEY (uploader_id) REFERENCES public.user_data_table(user_id) ON DELETE CASCADE;


--
-- TOC entry 4061 (class 2606 OID 16452)
-- Name: user_authentication fk_user; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_authentication
    ADD CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES public.user_data_table(user_id) ON DELETE CASCADE;


--
-- TOC entry 4071 (class 2606 OID 16961)
-- Name: user_favourite_topics fk_user; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_favourite_topics
    ADD CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES public.user_data_table(user_id) ON DELETE CASCADE;


--
-- TOC entry 4089 (class 2606 OID 33554)
-- Name: eco_votes fk_user; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.eco_votes
    ADD CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES public.user_data_table(user_id) ON DELETE CASCADE;


--
-- TOC entry 4090 (class 2606 OID 33620)
-- Name: video_votes fk_user; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.video_votes
    ADD CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES public.user_data_table(user_id) ON DELETE CASCADE;


--
-- TOC entry 4065 (class 2606 OID 16488)
-- Name: comments_table fk_video; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.comments_table
    ADD CONSTRAINT fk_video FOREIGN KEY (parent_video_id) REFERENCES public.video_data(video_id) ON DELETE CASCADE;


--
-- TOC entry 4070 (class 2606 OID 16951)
-- Name: video_tags fk_video; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.video_tags
    ADD CONSTRAINT fk_video FOREIGN KEY (video_id) REFERENCES public.video_data(video_id) ON DELETE CASCADE;


--
-- TOC entry 4091 (class 2606 OID 33615)
-- Name: video_votes fk_video; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.video_votes
    ADD CONSTRAINT fk_video FOREIGN KEY (video_id) REFERENCES public.video_data(video_id) ON DELETE CASCADE;


--
-- TOC entry 4066 (class 2606 OID 16509)
-- Name: follow_table follow_table_followee_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.follow_table
    ADD CONSTRAINT follow_table_followee_id_fkey FOREIGN KEY (followee_id) REFERENCES public.user_data_table(user_id) ON DELETE CASCADE;


--
-- TOC entry 4067 (class 2606 OID 16504)
-- Name: follow_table follow_table_follower_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.follow_table
    ADD CONSTRAINT follow_table_follower_id_fkey FOREIGN KEY (follower_id) REFERENCES public.user_data_table(user_id) ON DELETE CASCADE;


--
-- TOC entry 4076 (class 2606 OID 17077)
-- Name: video_history_table history_table_video_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.video_history_table
    ADD CONSTRAINT history_table_video_id_fkey FOREIGN KEY (video_id) REFERENCES public.video_data(video_id) ON DELETE CASCADE;


--
-- TOC entry 4077 (class 2606 OID 17072)
-- Name: video_history_table history_table_watcher_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.video_history_table
    ADD CONSTRAINT history_table_watcher_id_fkey FOREIGN KEY (watcher_id) REFERENCES public.user_data_table(user_id) ON DELETE CASCADE;


--
-- TOC entry 4068 (class 2606 OID 16545)
-- Name: connections user1_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.connections
    ADD CONSTRAINT user1_fk FOREIGN KEY (user1_id) REFERENCES public.user_data_table(user_id) ON DELETE CASCADE;


--
-- TOC entry 4069 (class 2606 OID 16550)
-- Name: connections user2_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.connections
    ADD CONSTRAINT user2_fk FOREIGN KEY (user2_id) REFERENCES public.user_data_table(user_id) ON DELETE CASCADE;


--
-- TOC entry 4062 (class 2606 OID 16446)
-- Name: user_authentication user_authentication_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_authentication
    ADD CONSTRAINT user_authentication_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.user_data_table(user_id) ON DELETE CASCADE;


--
-- TOC entry 4073 (class 2606 OID 17034)
-- Name: video_luv_events video_luvs_events_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.video_luv_events
    ADD CONSTRAINT video_luvs_events_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.user_data_table(user_id) ON DELETE CASCADE;


--
-- TOC entry 4074 (class 2606 OID 17029)
-- Name: video_luv_events video_luvs_events_video_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.video_luv_events
    ADD CONSTRAINT video_luvs_events_video_id_fkey FOREIGN KEY (video_id) REFERENCES public.video_data(video_id) ON DELETE CASCADE;


-- Completed on 2026-03-28 11:04:45 IST

--
-- PostgreSQL database dump complete
--

\unrestrict vZtAHfgjjXuyUTYkHdhBI9mSRcWasE2CdGANuP3AogeWTQTx8fae63g3DhR0Ed3

