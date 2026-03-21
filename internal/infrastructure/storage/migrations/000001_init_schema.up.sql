-- Extension for UUIDs
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 1. Events Table 
CREATE TABLE IF NOT EXISTS events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    slug TEXT UNIQUE NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_by UUID NOT NULL,
    updated_by UUID NOT NULL
);

-- 2. Persons Table 
CREATE TABLE IF NOT EXISTS persons (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    email TEXT UNIQUE,
    avatar_url TEXT,
    github_user TEXT,
    linkedin_url TEXT,
    twitter_url TEXT,
    website_url TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_by UUID NOT NULL,
    updated_by UUID NOT NULL
);

-- 3. ROLES (Speakers, Developers, Collaborators, Organizers)
CREATE TABLE IF NOT EXISTS speakers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    person_id UUID NOT NULL REFERENCES persons(id) ON DELETE CASCADE,
    event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    bio TEXT,
    company TEXT,
    UNIQUE(person_id, event_id),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_by UUID NOT NULL,
    updated_by UUID NOT NULL
);

CREATE TABLE IF NOT EXISTS developers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    person_id UUID NOT NULL REFERENCES persons(id) ON DELETE CASCADE,
    event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    role_description TEXT,
    UNIQUE(person_id, event_id),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_by UUID NOT NULL,
    updated_by UUID NOT NULL
);

CREATE TABLE IF NOT EXISTS collaborators (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    person_id UUID NOT NULL REFERENCES persons(id) ON DELETE CASCADE,
    event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    area TEXT,
    UNIQUE(person_id, event_id),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_by UUID NOT NULL,
    updated_by UUID NOT NULL
);

CREATE TABLE IF NOT EXISTS organizers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    person_id UUID NOT NULL REFERENCES persons(id) ON DELETE CASCADE,
    event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    company TEXT,
    role_description TEXT,
    UNIQUE(person_id, event_id),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_by UUID NOT NULL,
    updated_by UUID NOT NULL
)

-- 4. Sponsors Table
CREATE TABLE IF NOT EXISTS sponsors (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    logo_url TEXT,
    website_url TEXT,
    tier TEXT DEFAULT 'bronce',
    order_priority INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_by UUID NOT NULL,
    updated_by UUID NOT NULL
);

-- 5. Talks Table
CREATE TABLE IF NOT EXISTS talks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    description TEXT,
    tags TEXT[],
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_by UUID NOT NULL,
    updated_by UUID NOT NULL
);

CREATE TABLE IF NOT EXISTS talk_speakers (
    talk_id UUID NOT NULL REFERENCES talks(id) ON DELETE CASCADE,
    speaker_id UUID NOT NULL REFERENCES speakers(id) ON DELETE CASCADE,
    PRIMARY KEY (talk_id, speaker_id)
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_by UUID NOT NULL
);

-- 6. Tracks Table
CREATE TABLE IF NOT EXISTS tracks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    event_date DATE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 7. Scheduler Table
CREATE TABLE IF NOT EXISTS scheduler (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    track_id UUID NOT NULL REFERENCES tracks(id) ON DELETE CASCADE,
    talk_id UUID REFERENCES talks(id) ON DELETE SET NULL,
    
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    
    duration INTERVAL GENERATED ALWAYS AS (end_time - start_time) STORED,
    
    room TEXT DEFAULT 'Main Hall',
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

------------- TRIGGERS! -------------

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

DO $$ 
BEGIN
    -- Events
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'trig_update_events') THEN
        CREATE TRIGGER trig_update_events BEFORE UPDATE ON events 
        FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
    END IF;

    -- Persons
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'trig_update_persons') THEN
        CREATE TRIGGER trig_update_persons BEFORE UPDATE ON persons 
        FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
    END IF;

    -- Speakers
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'trig_update_speakers') THEN
        CREATE TRIGGER trig_update_speakers BEFORE UPDATE ON speakers 
        FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
    END IF;

    -- Developers
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'trig_update_developers') THEN
        CREATE TRIGGER trig_update_developers BEFORE UPDATE ON developers 
        FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
    END IF;

    -- Collaborators
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'trig_update_collaborators') THEN
        CREATE TRIGGER trig_update_collaborators BEFORE UPDATE ON collaborators 
        FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
    END IF;

    -- Sponsors
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'trig_update_sponsors') THEN
        CREATE TRIGGER trig_update_sponsors BEFORE UPDATE ON sponsors 
        FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
    END IF;

    -- Organizers
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'trig_update_organizers') THEN
        CREATE TRIGGER trig_update_organizers BEFORE UPDATE ON organizers 
        FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
    END IF;

    -- Talks
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'trig_update_talks') THEN
        CREATE TRIGGER trig_update_talks BEFORE UPDATE ON talks 
        FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
    END IF;

    -- Talks-Speakers
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'trig_update_talk_speakers') THEN
        CREATE TRIGGER trig_update_talk_speakers BEFORE UPDATE ON talk_speakers 
        FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
    END IF;

    -- Tracks
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'trig_update_tracks') THEN
        CREATE TRIGGER trig_update_tracks BEFORE UPDATE ON tracks 
        FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
    END IF;

    -- Scheduler
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'trig_update_scheduler') THEN
        CREATE TRIGGER trig_update_scheduler BEFORE UPDATE ON scheduler 
        FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
    END IF;
END $$;