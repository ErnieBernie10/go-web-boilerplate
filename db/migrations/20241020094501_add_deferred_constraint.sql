-- migrate:up
alter table frame drop constraint frame_file_id_fkey;
alter table frame add constraint frame_file_id_fkey foreign key (file_id) references file(id) deferrable initially deferred;

-- migrate:down
alter table frame drop constraint frame_file_id_fkey;
