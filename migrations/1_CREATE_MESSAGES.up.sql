BEGIN;

    CREATE TABLE MESSAGES (
        id UUID NOT NULL DEFAULT gen_random_uuid(),
        value TEXT NOT NULL,
        processed BOOLEAN NOT NULL DEFAULT false,
        date_create TIMESTAMP WITH TIME ZONE NOT NULL,
        date_process TIMESTAMP WITH TIME ZONE,
        PRIMARY KEY (id)
    );

    COMMIT;
END;