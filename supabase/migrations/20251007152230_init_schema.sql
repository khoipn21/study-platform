--
-- PostgreSQL database dump
--

-- Dumped from database version 17.6
-- Dumped by pg_dump version 17.5 (Ubuntu 17.5-1.pgdg24.04+1)

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

DROP EVENT TRIGGER IF EXISTS "pgrst_drop_watch";
DROP EVENT TRIGGER IF EXISTS "pgrst_ddl_watch";
DROP EVENT TRIGGER IF EXISTS "issue_pg_net_access";
DROP EVENT TRIGGER IF EXISTS "issue_pg_graphql_access";
DROP EVENT TRIGGER IF EXISTS "issue_pg_cron_access";
DROP EVENT TRIGGER IF EXISTS "issue_graphql_placeholder";
DROP PUBLICATION IF EXISTS "supabase_realtime";
ALTER TABLE IF EXISTS ONLY "storage"."s3_multipart_uploads_parts" DROP CONSTRAINT IF EXISTS "s3_multipart_uploads_parts_upload_id_fkey";
ALTER TABLE IF EXISTS ONLY "storage"."s3_multipart_uploads_parts" DROP CONSTRAINT IF EXISTS "s3_multipart_uploads_parts_bucket_id_fkey";
ALTER TABLE IF EXISTS ONLY "storage"."s3_multipart_uploads" DROP CONSTRAINT IF EXISTS "s3_multipart_uploads_bucket_id_fkey";
ALTER TABLE IF EXISTS ONLY "storage"."prefixes" DROP CONSTRAINT IF EXISTS "prefixes_bucketId_fkey";
ALTER TABLE IF EXISTS ONLY "storage"."objects" DROP CONSTRAINT IF EXISTS "objects_bucketId_fkey";
ALTER TABLE IF EXISTS ONLY "storage"."iceberg_tables" DROP CONSTRAINT IF EXISTS "iceberg_tables_namespace_id_fkey";
ALTER TABLE IF EXISTS ONLY "storage"."iceberg_tables" DROP CONSTRAINT IF EXISTS "iceberg_tables_bucket_id_fkey";
ALTER TABLE IF EXISTS ONLY "storage"."iceberg_namespaces" DROP CONSTRAINT IF EXISTS "iceberg_namespaces_bucket_id_fkey";
ALTER TABLE IF EXISTS ONLY "public"."viewing_sessions" DROP CONSTRAINT IF EXISTS "viewing_sessions_video_id_fkey";
ALTER TABLE IF EXISTS ONLY "public"."video_qualities" DROP CONSTRAINT IF EXISTS "video_qualities_video_id_fkey";
ALTER TABLE IF EXISTS ONLY "public"."video_permissions" DROP CONSTRAINT IF EXISTS "video_permissions_video_id_fkey";
ALTER TABLE IF EXISTS ONLY "public"."video_analytics" DROP CONSTRAINT IF EXISTS "video_analytics_video_id_fkey";
ALTER TABLE IF EXISTS ONLY "public"."user_payment_methods" DROP CONSTRAINT IF EXISTS "user_payment_methods_user_id_fkey";
ALTER TABLE IF EXISTS ONLY "public"."transactions" DROP CONSTRAINT IF EXISTS "transactions_user_id_fkey";
ALTER TABLE IF EXISTS ONLY "public"."transactions" DROP CONSTRAINT IF EXISTS "transactions_payment_method_id_fkey";
ALTER TABLE IF EXISTS ONLY "public"."subscriptions" DROP CONSTRAINT IF EXISTS "subscriptions_user_id_fkey";
ALTER TABLE IF EXISTS ONLY "public"."subscriptions" DROP CONSTRAINT IF EXISTS "subscriptions_payment_method_id_fkey";
ALTER TABLE IF EXISTS ONLY "public"."stripe_products" DROP CONSTRAINT IF EXISTS "stripe_products_course_id_fkey";
ALTER TABLE IF EXISTS ONLY "public"."stripe_customers" DROP CONSTRAINT IF EXISTS "stripe_customers_user_id_fkey";
ALTER TABLE IF EXISTS ONLY "public"."progress" DROP CONSTRAINT IF EXISTS "progress_user_id_fkey";
ALTER TABLE IF EXISTS ONLY "public"."payment_methods" DROP CONSTRAINT IF EXISTS "payment_methods_user_id_fkey";
ALTER TABLE IF EXISTS ONLY "public"."payment_events" DROP CONSTRAINT IF EXISTS "payment_events_user_id_fkey";
ALTER TABLE IF EXISTS ONLY "public"."payment_events" DROP CONSTRAINT IF EXISTS "payment_events_transaction_id_fkey";
ALTER TABLE IF EXISTS ONLY "public"."payment_events" DROP CONSTRAINT IF EXISTS "payment_events_course_id_fkey";
ALTER TABLE IF EXISTS ONLY "public"."oauth_accounts" DROP CONSTRAINT IF EXISTS "oauth_accounts_user_id_fkey";
ALTER TABLE IF EXISTS ONLY "public"."lectures" DROP CONSTRAINT IF EXISTS "lectures_course_id_fkey";
ALTER TABLE IF EXISTS ONLY "public"."lecture_preview_sessions" DROP CONSTRAINT IF EXISTS "lecture_preview_sessions_lecture_id_fkey";
ALTER TABLE IF EXISTS ONLY "public"."forum_votes" DROP CONSTRAINT IF EXISTS "forum_votes_user_id_fkey";
ALTER TABLE IF EXISTS ONLY "public"."forum_votes" DROP CONSTRAINT IF EXISTS "forum_votes_post_id_fkey";
ALTER TABLE IF EXISTS ONLY "public"."forum_topics" DROP CONSTRAINT IF EXISTS "forum_topics_creator_id_fkey";
ALTER TABLE IF EXISTS ONLY "public"."forum_topic_subscriptions" DROP CONSTRAINT IF EXISTS "forum_topic_subscriptions_user_id_fkey";
ALTER TABLE IF EXISTS ONLY "public"."forum_topic_subscriptions" DROP CONSTRAINT IF EXISTS "forum_topic_subscriptions_topic_id_fkey";
ALTER TABLE IF EXISTS ONLY "public"."forum_posts" DROP CONSTRAINT IF EXISTS "forum_posts_topic_id_fkey";
ALTER TABLE IF EXISTS ONLY "public"."forum_posts" DROP CONSTRAINT IF EXISTS "forum_posts_parent_id_fkey";
ALTER TABLE IF EXISTS ONLY "public"."forum_posts" DROP CONSTRAINT IF EXISTS "forum_posts_author_id_fkey";
ALTER TABLE IF EXISTS ONLY "public"."forum_notifications" DROP CONSTRAINT IF EXISTS "forum_notifications_user_id_fkey";
ALTER TABLE IF EXISTS ONLY "public"."forum_mentions" DROP CONSTRAINT IF EXISTS "forum_mentions_post_id_fkey";
ALTER TABLE IF EXISTS ONLY "public"."forum_mentions" DROP CONSTRAINT IF EXISTS "forum_mentions_mentioner_user_id_fkey";
ALTER TABLE IF EXISTS ONLY "public"."forum_mentions" DROP CONSTRAINT IF EXISTS "forum_mentions_mentioned_user_id_fkey";
ALTER TABLE IF EXISTS ONLY "public"."progress" DROP CONSTRAINT IF EXISTS "fk_progress_course_id";
ALTER TABLE IF EXISTS ONLY "public"."notes" DROP CONSTRAINT IF EXISTS "fk_notes_user";
ALTER TABLE IF EXISTS ONLY "public"."notes" DROP CONSTRAINT IF EXISTS "fk_notes_lecture";
ALTER TABLE IF EXISTS ONLY "public"."notes" DROP CONSTRAINT IF EXISTS "fk_notes_course";
ALTER TABLE IF EXISTS ONLY "public"."lemon_squeezy_variants" DROP CONSTRAINT IF EXISTS "fk_lemon_squeezy_variants_product";
ALTER TABLE IF EXISTS ONLY "public"."lecture_resources" DROP CONSTRAINT IF EXISTS "fk_lecture_resources_lecture_id";
ALTER TABLE IF EXISTS ONLY "public"."lecture_resources" DROP CONSTRAINT IF EXISTS "fk_lecture_resources_file_id";
ALTER TABLE IF EXISTS ONLY "public"."enrollments" DROP CONSTRAINT IF EXISTS "fk_enrollments_user_id";
ALTER TABLE IF EXISTS ONLY "public"."courses" DROP CONSTRAINT IF EXISTS "fk_courses_instructor_id";
ALTER TABLE IF EXISTS ONLY "public"."file_permissions" DROP CONSTRAINT IF EXISTS "file_permissions_file_id_fkey";
ALTER TABLE IF EXISTS ONLY "public"."enrollments" DROP CONSTRAINT IF EXISTS "enrollments_transaction_id_fkey";
ALTER TABLE IF EXISTS ONLY "public"."enrollments" DROP CONSTRAINT IF EXISTS "enrollments_course_id_fkey";
ALTER TABLE IF EXISTS ONLY "public"."course_resources" DROP CONSTRAINT IF EXISTS "course_resources_file_id_fkey";
ALTER TABLE IF EXISTS ONLY "public"."course_resources" DROP CONSTRAINT IF EXISTS "course_resources_course_id_fkey";
ALTER TABLE IF EXISTS ONLY "public"."course_access_logs" DROP CONSTRAINT IF EXISTS "course_access_logs_user_id_fkey";
ALTER TABLE IF EXISTS ONLY "public"."course_access_logs" DROP CONSTRAINT IF EXISTS "course_access_logs_transaction_id_fkey";
ALTER TABLE IF EXISTS ONLY "public"."course_access_logs" DROP CONSTRAINT IF EXISTS "course_access_logs_lecture_id_fkey";
ALTER TABLE IF EXISTS ONLY "public"."course_access_logs" DROP CONSTRAINT IF EXISTS "course_access_logs_course_id_fkey";
ALTER TABLE IF EXISTS ONLY "public"."chat_history" DROP CONSTRAINT IF EXISTS "chat_history_user_id_fkey";
ALTER TABLE IF EXISTS ONLY "public"."audit_logs" DROP CONSTRAINT IF EXISTS "audit_logs_user_id_fkey";
ALTER TABLE IF EXISTS ONLY "public"."audit_logs" DROP CONSTRAINT IF EXISTS "audit_logs_transaction_id_fkey";
ALTER TABLE IF EXISTS ONLY "public"."audit_logs" DROP CONSTRAINT IF EXISTS "audit_logs_lecture_id_fkey";
ALTER TABLE IF EXISTS ONLY "public"."audit_logs" DROP CONSTRAINT IF EXISTS "audit_logs_course_id_fkey";
ALTER TABLE IF EXISTS ONLY "auth"."sso_domains" DROP CONSTRAINT IF EXISTS "sso_domains_sso_provider_id_fkey";
ALTER TABLE IF EXISTS ONLY "auth"."sessions" DROP CONSTRAINT IF EXISTS "sessions_user_id_fkey";
ALTER TABLE IF EXISTS ONLY "auth"."sessions" DROP CONSTRAINT IF EXISTS "sessions_oauth_client_id_fkey";
ALTER TABLE IF EXISTS ONLY "auth"."saml_relay_states" DROP CONSTRAINT IF EXISTS "saml_relay_states_sso_provider_id_fkey";
ALTER TABLE IF EXISTS ONLY "auth"."saml_relay_states" DROP CONSTRAINT IF EXISTS "saml_relay_states_flow_state_id_fkey";
ALTER TABLE IF EXISTS ONLY "auth"."saml_providers" DROP CONSTRAINT IF EXISTS "saml_providers_sso_provider_id_fkey";
ALTER TABLE IF EXISTS ONLY "auth"."refresh_tokens" DROP CONSTRAINT IF EXISTS "refresh_tokens_session_id_fkey";
ALTER TABLE IF EXISTS ONLY "auth"."one_time_tokens" DROP CONSTRAINT IF EXISTS "one_time_tokens_user_id_fkey";
ALTER TABLE IF EXISTS ONLY "auth"."oauth_consents" DROP CONSTRAINT IF EXISTS "oauth_consents_user_id_fkey";
ALTER TABLE IF EXISTS ONLY "auth"."oauth_consents" DROP CONSTRAINT IF EXISTS "oauth_consents_client_id_fkey";
ALTER TABLE IF EXISTS ONLY "auth"."oauth_authorizations" DROP CONSTRAINT IF EXISTS "oauth_authorizations_user_id_fkey";
ALTER TABLE IF EXISTS ONLY "auth"."oauth_authorizations" DROP CONSTRAINT IF EXISTS "oauth_authorizations_client_id_fkey";
ALTER TABLE IF EXISTS ONLY "auth"."mfa_factors" DROP CONSTRAINT IF EXISTS "mfa_factors_user_id_fkey";
ALTER TABLE IF EXISTS ONLY "auth"."mfa_challenges" DROP CONSTRAINT IF EXISTS "mfa_challenges_auth_factor_id_fkey";
ALTER TABLE IF EXISTS ONLY "auth"."mfa_amr_claims" DROP CONSTRAINT IF EXISTS "mfa_amr_claims_session_id_fkey";
ALTER TABLE IF EXISTS ONLY "auth"."identities" DROP CONSTRAINT IF EXISTS "identities_user_id_fkey";
ALTER TABLE IF EXISTS ONLY "_realtime"."extensions" DROP CONSTRAINT IF EXISTS "extensions_tenant_external_id_fkey";
DROP TRIGGER IF EXISTS "update_objects_updated_at" ON "storage"."objects";
DROP TRIGGER IF EXISTS "prefixes_delete_hierarchy" ON "storage"."prefixes";
DROP TRIGGER IF EXISTS "prefixes_create_hierarchy" ON "storage"."prefixes";
DROP TRIGGER IF EXISTS "objects_update_create_prefix" ON "storage"."objects";
DROP TRIGGER IF EXISTS "objects_insert_create_prefix" ON "storage"."objects";
DROP TRIGGER IF EXISTS "objects_delete_delete_prefix" ON "storage"."objects";
DROP TRIGGER IF EXISTS "enforce_bucket_name_length_trigger" ON "storage"."buckets";
DROP TRIGGER IF EXISTS "tr_check_filters" ON "realtime"."subscription";
DROP TRIGGER IF EXISTS "update_stripe_products_updated_at" ON "public"."stripe_products";
DROP TRIGGER IF EXISTS "update_stripe_customers_updated_at" ON "public"."stripe_customers";
DROP TRIGGER IF EXISTS "update_progress_updated_at" ON "public"."progress";
DROP TRIGGER IF EXISTS "update_enrollments_updated_at" ON "public"."enrollments";
DROP TRIGGER IF EXISTS "trigger_update_network_analytics" ON "public"."network_metrics";
DROP TRIGGER IF EXISTS "trigger_update_enrollment_payment_status" ON "public"."transactions";
DROP TRIGGER IF EXISTS "trigger_prevent_course_resources_insert" ON "public"."course_resources";
DROP TRIGGER IF EXISTS "trigger_notes_updated_at" ON "public"."notes";
DROP INDEX IF EXISTS "supabase_functions"."supabase_functions_hooks_request_id_idx";
DROP INDEX IF EXISTS "supabase_functions"."supabase_functions_hooks_h_table_id_h_name_idx";
DROP INDEX IF EXISTS "storage"."objects_bucket_id_level_idx";
DROP INDEX IF EXISTS "storage"."name_prefix_search";
DROP INDEX IF EXISTS "storage"."idx_prefixes_lower_name";
DROP INDEX IF EXISTS "storage"."idx_objects_lower_name";
DROP INDEX IF EXISTS "storage"."idx_objects_bucket_id_name";
DROP INDEX IF EXISTS "storage"."idx_name_bucket_level_unique";
DROP INDEX IF EXISTS "storage"."idx_multipart_uploads_list";
DROP INDEX IF EXISTS "storage"."idx_iceberg_tables_namespace_id";
DROP INDEX IF EXISTS "storage"."idx_iceberg_namespaces_bucket_id";
DROP INDEX IF EXISTS "storage"."bucketid_objname";
DROP INDEX IF EXISTS "storage"."bname";
DROP INDEX IF EXISTS "realtime"."subscription_subscription_id_entity_filters_key";
DROP INDEX IF EXISTS "realtime"."messages_inserted_at_topic_index";
DROP INDEX IF EXISTS "realtime"."ix_realtime_subscription_entity";
DROP INDEX IF EXISTS "public"."idx_webhook_events_processed_at";
DROP INDEX IF EXISTS "public"."idx_webhook_events_event_name";
DROP INDEX IF EXISTS "public"."idx_webhook_events_event_id";
DROP INDEX IF EXISTS "public"."idx_viewing_sessions_user_video";
DROP INDEX IF EXISTS "public"."idx_viewing_sessions_session";
DROP INDEX IF EXISTS "public"."idx_videos_user_id";
DROP INDEX IF EXISTS "public"."idx_videos_status_visibility";
DROP INDEX IF EXISTS "public"."idx_videos_status";
DROP INDEX IF EXISTS "public"."idx_videos_search";
DROP INDEX IF EXISTS "public"."idx_videos_deleted_at";
DROP INDEX IF EXISTS "public"."idx_videos_course_lecture";
DROP INDEX IF EXISTS "public"."idx_videos_cloudflare_uid";
DROP INDEX IF EXISTS "public"."idx_video_qualities_video_id";
DROP INDEX IF EXISTS "public"."idx_video_permissions_video_user";
DROP INDEX IF EXISTS "public"."idx_video_analytics_video_date";
DROP INDEX IF EXISTS "public"."idx_users_username";
DROP INDEX IF EXISTS "public"."idx_users_role";
DROP INDEX IF EXISTS "public"."idx_users_provider_id";
DROP INDEX IF EXISTS "public"."idx_users_email";
DROP INDEX IF EXISTS "public"."idx_user_payment_methods_user_id";
DROP INDEX IF EXISTS "public"."idx_user_payment_methods_default";
DROP INDEX IF EXISTS "public"."idx_user_payment_methods_active";
DROP INDEX IF EXISTS "public"."idx_upload_sessions_user_status";
DROP INDEX IF EXISTS "public"."idx_upload_sessions_upload_id_unique";
DROP INDEX IF EXISTS "public"."idx_upload_sessions_upload_id";
DROP INDEX IF EXISTS "public"."idx_upload_sessions_expires_at";
DROP INDEX IF EXISTS "public"."idx_upload_sessions_bucket_key";
DROP INDEX IF EXISTS "public"."idx_transactions_webhook_event_id";
DROP INDEX IF EXISTS "public"."idx_transactions_user_id";
DROP INDEX IF EXISTS "public"."idx_transactions_user_course";
DROP INDEX IF EXISTS "public"."idx_transactions_stripe_session_id";
DROP INDEX IF EXISTS "public"."idx_transactions_stripe_payment_intent_id";
DROP INDEX IF EXISTS "public"."idx_transactions_stripe_customer_id";
DROP INDEX IF EXISTS "public"."idx_transactions_stripe_charge_id";
DROP INDEX IF EXISTS "public"."idx_transactions_status_updated";
DROP INDEX IF EXISTS "public"."idx_transactions_status";
DROP INDEX IF EXISTS "public"."idx_transactions_pending";
DROP INDEX IF EXISTS "public"."idx_transactions_payment_verified_at";
DROP INDEX IF EXISTS "public"."idx_transactions_payment_provider";
DROP INDEX IF EXISTS "public"."idx_transactions_lemon_squeezy_order_id";
DROP INDEX IF EXISTS "public"."idx_transactions_lemon_squeezy_checkout_id";
DROP INDEX IF EXISTS "public"."idx_transactions_course_id";
DROP INDEX IF EXISTS "public"."idx_transactions_completed";
DROP INDEX IF EXISTS "public"."idx_transaction_user_status_course";
DROP INDEX IF EXISTS "public"."idx_subscriptions_user_id";
DROP INDEX IF EXISTS "public"."idx_subscriptions_status";
DROP INDEX IF EXISTS "public"."idx_stripe_webhook_events_unprocessed";
DROP INDEX IF EXISTS "public"."idx_stripe_webhook_events_type";
DROP INDEX IF EXISTS "public"."idx_stripe_webhook_events_stripe_event_id";
DROP INDEX IF EXISTS "public"."idx_stripe_webhook_events_processing_attempts";
DROP INDEX IF EXISTS "public"."idx_stripe_webhook_events_processed";
DROP INDEX IF EXISTS "public"."idx_stripe_webhook_events_event_type";
DROP INDEX IF EXISTS "public"."idx_stripe_webhook_events_created_at";
DROP INDEX IF EXISTS "public"."idx_stripe_products_stripe_product_id";
DROP INDEX IF EXISTS "public"."idx_stripe_products_stripe_price_id";
DROP INDEX IF EXISTS "public"."idx_stripe_products_course_id";
DROP INDEX IF EXISTS "public"."idx_stripe_customers_user_id";
DROP INDEX IF EXISTS "public"."idx_stripe_customers_stripe_customer_id";
DROP INDEX IF EXISTS "public"."idx_progress_watch_time";
DROP INDEX IF EXISTS "public"."idx_progress_user_id";
DROP INDEX IF EXISTS "public"."idx_progress_user_course";
DROP INDEX IF EXISTS "public"."idx_progress_progress_percentage";
DROP INDEX IF EXISTS "public"."idx_progress_lecture_id";
DROP INDEX IF EXISTS "public"."idx_progress_last_accessed";
DROP INDEX IF EXISTS "public"."idx_progress_is_completed";
DROP INDEX IF EXISTS "public"."idx_progress_course_id";
DROP INDEX IF EXISTS "public"."idx_progress_completed_at";
DROP INDEX IF EXISTS "public"."idx_progress_completed";
DROP INDEX IF EXISTS "public"."idx_payment_methods_user_id";
DROP INDEX IF EXISTS "public"."idx_payment_methods_stripe_payment_method_id";
DROP INDEX IF EXISTS "public"."idx_payment_methods_stripe_customer_id";
DROP INDEX IF EXISTS "public"."idx_payment_events_user_course";
DROP INDEX IF EXISTS "public"."idx_payment_events_type_processed";
DROP INDEX IF EXISTS "public"."idx_payment_events_transaction_id";
DROP INDEX IF EXISTS "public"."idx_payment_events_provider_event_id";
DROP INDEX IF EXISTS "public"."idx_payment_analytics_date";
DROP INDEX IF EXISTS "public"."idx_oauth_accounts_user_id";
DROP INDEX IF EXISTS "public"."idx_oauth_accounts_provider_id";
DROP INDEX IF EXISTS "public"."idx_oauth_accounts_provider";
DROP INDEX IF EXISTS "public"."idx_notes_user_course_lecture";
DROP INDEX IF EXISTS "public"."idx_notes_timestamp";
DROP INDEX IF EXISTS "public"."idx_notes_created_at";
DROP INDEX IF EXISTS "public"."idx_notes_course_lecture";
DROP INDEX IF EXISTS "public"."idx_network_metrics_user_timestamp";
DROP INDEX IF EXISTS "public"."idx_network_metrics_timestamp";
DROP INDEX IF EXISTS "public"."idx_network_metrics_session";
DROP INDEX IF EXISTS "public"."idx_network_metrics_quality_score";
DROP INDEX IF EXISTS "public"."idx_network_metrics_device_type";
DROP INDEX IF EXISTS "public"."idx_network_metrics_connection_type";
DROP INDEX IF EXISTS "public"."idx_network_events_user";
DROP INDEX IF EXISTS "public"."idx_network_events_unresolved";
DROP INDEX IF EXISTS "public"."idx_network_events_type";
DROP INDEX IF EXISTS "public"."idx_network_events_timestamp";
DROP INDEX IF EXISTS "public"."idx_network_events_severity";
DROP INDEX IF EXISTS "public"."idx_network_events_session";
DROP INDEX IF EXISTS "public"."idx_network_dashboard_video_id";
DROP INDEX IF EXISTS "public"."idx_network_dashboard_session_status";
DROP INDEX IF EXISTS "public"."idx_network_analytics_video";
DROP INDEX IF EXISTS "public"."idx_network_analytics_user";
DROP INDEX IF EXISTS "public"."idx_network_analytics_date";
DROP INDEX IF EXISTS "public"."idx_lemon_squeezy_variants_variant_id";
DROP INDEX IF EXISTS "public"."idx_lemon_squeezy_variants_product_id";
DROP INDEX IF EXISTS "public"."idx_lemon_squeezy_products_product_id";
DROP INDEX IF EXISTS "public"."idx_lectures_status";
DROP INDEX IF EXISTS "public"."idx_lectures_preview_available";
DROP INDEX IF EXISTS "public"."idx_lectures_order_number";
DROP INDEX IF EXISTS "public"."idx_lectures_deleted_at";
DROP INDEX IF EXISTS "public"."idx_lectures_course_id";
DROP INDEX IF EXISTS "public"."idx_lectures_access_level";
DROP INDEX IF EXISTS "public"."idx_lecture_resources_type";
DROP INDEX IF EXISTS "public"."idx_lecture_resources_order";
DROP INDEX IF EXISTS "public"."idx_lecture_resources_lecture_id";
DROP INDEX IF EXISTS "public"."idx_lecture_resources_file_id";
DROP INDEX IF EXISTS "public"."idx_lecture_preview_sessions_user_lecture";
DROP INDEX IF EXISTS "public"."idx_lecture_preview_sessions_last_accessed";
DROP INDEX IF EXISTS "public"."idx_lecture_preview_sessions_exhausted";
DROP INDEX IF EXISTS "public"."idx_forum_votes_vote_type";
DROP INDEX IF EXISTS "public"."idx_forum_votes_user_id";
DROP INDEX IF EXISTS "public"."idx_forum_votes_post_id";
DROP INDEX IF EXISTS "public"."idx_forum_topics_status";
DROP INDEX IF EXISTS "public"."idx_forum_topics_pin_order";
DROP INDEX IF EXISTS "public"."idx_forum_topics_is_pinned";
DROP INDEX IF EXISTS "public"."idx_forum_topics_creator_id";
DROP INDEX IF EXISTS "public"."idx_forum_topics_course_status";
DROP INDEX IF EXISTS "public"."idx_forum_topics_course_id";
DROP INDEX IF EXISTS "public"."idx_forum_subscriptions_user_id";
DROP INDEX IF EXISTS "public"."idx_forum_subscriptions_topic_id";
DROP INDEX IF EXISTS "public"."idx_forum_posts_topic_status";
DROP INDEX IF EXISTS "public"."idx_forum_posts_topic_id";
DROP INDEX IF EXISTS "public"."idx_forum_posts_status";
DROP INDEX IF EXISTS "public"."idx_forum_posts_pin_order";
DROP INDEX IF EXISTS "public"."idx_forum_posts_parent_id";
DROP INDEX IF EXISTS "public"."idx_forum_posts_is_pinned";
DROP INDEX IF EXISTS "public"."idx_forum_posts_is_answer";
DROP INDEX IF EXISTS "public"."idx_forum_posts_created_at";
DROP INDEX IF EXISTS "public"."idx_forum_posts_author_id";
DROP INDEX IF EXISTS "public"."idx_forum_notifications_user_id";
DROP INDEX IF EXISTS "public"."idx_forum_notifications_unread";
DROP INDEX IF EXISTS "public"."idx_forum_notifications_type";
DROP INDEX IF EXISTS "public"."idx_forum_notifications_reference";
DROP INDEX IF EXISTS "public"."idx_forum_mentions_unread";
DROP INDEX IF EXISTS "public"."idx_forum_mentions_post_id";
DROP INDEX IF EXISTS "public"."idx_forum_mentions_mentioner_user";
DROP INDEX IF EXISTS "public"."idx_forum_mentions_mentioned_user";
DROP INDEX IF EXISTS "public"."idx_files_user_id";
DROP INDEX IF EXISTS "public"."idx_files_is_public";
DROP INDEX IF EXISTS "public"."idx_files_deleted_at";
DROP INDEX IF EXISTS "public"."idx_files_created_at";
DROP INDEX IF EXISTS "public"."idx_files_content_type";
DROP INDEX IF EXISTS "public"."idx_files_bucket_object_active";
DROP INDEX IF EXISTS "public"."idx_files_bucket_key";
DROP INDEX IF EXISTS "public"."idx_file_permissions_user";
DROP INDEX IF EXISTS "public"."idx_file_permissions_unique";
DROP INDEX IF EXISTS "public"."idx_file_permissions_type";
DROP INDEX IF EXISTS "public"."idx_file_permissions_file_user";
DROP INDEX IF EXISTS "public"."idx_enrollments_user_id";
DROP INDEX IF EXISTS "public"."idx_enrollments_updated_at";
DROP INDEX IF EXISTS "public"."idx_enrollments_transaction_id";
DROP INDEX IF EXISTS "public"."idx_enrollments_total_lectures";
DROP INDEX IF EXISTS "public"."idx_enrollments_status";
DROP INDEX IF EXISTS "public"."idx_enrollments_progress";
DROP INDEX IF EXISTS "public"."idx_enrollments_payment_verified_at";
DROP INDEX IF EXISTS "public"."idx_enrollments_payment_status";
DROP INDEX IF EXISTS "public"."idx_enrollments_last_accessed";
DROP INDEX IF EXISTS "public"."idx_enrollments_enrolled_at";
DROP INDEX IF EXISTS "public"."idx_enrollments_deleted_at";
DROP INDEX IF EXISTS "public"."idx_enrollments_created_at";
DROP INDEX IF EXISTS "public"."idx_enrollments_course_id";
DROP INDEX IF EXISTS "public"."idx_enrollments_completed_lectures";
DROP INDEX IF EXISTS "public"."idx_enrollments_access_expires_at";
DROP INDEX IF EXISTS "public"."idx_enrollment_user_payment_status";
DROP INDEX IF EXISTS "public"."idx_courses_tags";
DROP INDEX IF EXISTS "public"."idx_courses_status";
DROP INDEX IF EXISTS "public"."idx_courses_search";
DROP INDEX IF EXISTS "public"."idx_courses_rating";
DROP INDEX IF EXISTS "public"."idx_courses_price";
DROP INDEX IF EXISTS "public"."idx_courses_preview_enabled";
DROP INDEX IF EXISTS "public"."idx_courses_level";
DROP INDEX IF EXISTS "public"."idx_courses_lemon_squeezy_variant_id";
DROP INDEX IF EXISTS "public"."idx_courses_lemon_squeezy_product_id";
DROP INDEX IF EXISTS "public"."idx_courses_is_paid";
DROP INDEX IF EXISTS "public"."idx_courses_instructor_id";
DROP INDEX IF EXISTS "public"."idx_courses_deleted_at";
DROP INDEX IF EXISTS "public"."idx_courses_created_at";
DROP INDEX IF EXISTS "public"."idx_courses_category";
DROP INDEX IF EXISTS "public"."idx_course_resources_type";
DROP INDEX IF EXISTS "public"."idx_course_resources_order";
DROP INDEX IF EXISTS "public"."idx_course_resources_file_id";
DROP INDEX IF EXISTS "public"."idx_course_resources_course_id";
DROP INDEX IF EXISTS "public"."idx_course_access_user_course_time";
DROP INDEX IF EXISTS "public"."idx_course_access_logs_user_lecture";
DROP INDEX IF EXISTS "public"."idx_course_access_logs_user_course";
DROP INDEX IF EXISTS "public"."idx_course_access_logs_created_at";
DROP INDEX IF EXISTS "public"."idx_course_access_logs_access_type";
DROP INDEX IF EXISTS "public"."idx_course_access_cache_user_course";
DROP INDEX IF EXISTS "public"."idx_course_access_cache_expires_at";
DROP INDEX IF EXISTS "public"."idx_course_access_cache_access_level";
DROP INDEX IF EXISTS "public"."idx_course_access_analytics_course";
DROP INDEX IF EXISTS "public"."idx_chat_history_user_id";
DROP INDEX IF EXISTS "public"."idx_chat_history_created_at";
DROP INDEX IF EXISTS "public"."idx_bandwidth_tests_user";
DROP INDEX IF EXISTS "public"."idx_bandwidth_tests_type";
DROP INDEX IF EXISTS "public"."idx_bandwidth_tests_timestamp";
DROP INDEX IF EXISTS "public"."idx_bandwidth_tests_session";
DROP INDEX IF EXISTS "public"."idx_audit_logs_user_id";
DROP INDEX IF EXISTS "public"."idx_audit_logs_created_at";
DROP INDEX IF EXISTS "public"."idx_audit_logs_course_id";
DROP INDEX IF EXISTS "public"."idx_audit_logs_action";
DROP INDEX IF EXISTS "public"."idx_adaptive_rules_priority";
DROP INDEX IF EXISTS "public"."idx_adaptive_rules_active";
DROP INDEX IF EXISTS "auth"."users_is_anonymous_idx";
DROP INDEX IF EXISTS "auth"."users_instance_id_idx";
DROP INDEX IF EXISTS "auth"."users_instance_id_email_idx";
DROP INDEX IF EXISTS "auth"."users_email_partial_key";
DROP INDEX IF EXISTS "auth"."user_id_created_at_idx";
DROP INDEX IF EXISTS "auth"."unique_phone_factor_per_user";
DROP INDEX IF EXISTS "auth"."sso_providers_resource_id_pattern_idx";
DROP INDEX IF EXISTS "auth"."sso_providers_resource_id_idx";
DROP INDEX IF EXISTS "auth"."sso_domains_sso_provider_id_idx";
DROP INDEX IF EXISTS "auth"."sso_domains_domain_idx";
DROP INDEX IF EXISTS "auth"."sessions_user_id_idx";
DROP INDEX IF EXISTS "auth"."sessions_oauth_client_id_idx";
DROP INDEX IF EXISTS "auth"."sessions_not_after_idx";
DROP INDEX IF EXISTS "auth"."saml_relay_states_sso_provider_id_idx";
DROP INDEX IF EXISTS "auth"."saml_relay_states_for_email_idx";
DROP INDEX IF EXISTS "auth"."saml_relay_states_created_at_idx";
DROP INDEX IF EXISTS "auth"."saml_providers_sso_provider_id_idx";
DROP INDEX IF EXISTS "auth"."refresh_tokens_updated_at_idx";
DROP INDEX IF EXISTS "auth"."refresh_tokens_session_id_revoked_idx";
DROP INDEX IF EXISTS "auth"."refresh_tokens_parent_idx";
DROP INDEX IF EXISTS "auth"."refresh_tokens_instance_id_user_id_idx";
DROP INDEX IF EXISTS "auth"."refresh_tokens_instance_id_idx";
DROP INDEX IF EXISTS "auth"."recovery_token_idx";
DROP INDEX IF EXISTS "auth"."reauthentication_token_idx";
DROP INDEX IF EXISTS "auth"."one_time_tokens_user_id_token_type_key";
DROP INDEX IF EXISTS "auth"."one_time_tokens_token_hash_hash_idx";
DROP INDEX IF EXISTS "auth"."one_time_tokens_relates_to_hash_idx";
DROP INDEX IF EXISTS "auth"."oauth_consents_user_order_idx";
DROP INDEX IF EXISTS "auth"."oauth_consents_active_user_client_idx";
DROP INDEX IF EXISTS "auth"."oauth_consents_active_client_idx";
DROP INDEX IF EXISTS "auth"."oauth_clients_deleted_at_idx";
DROP INDEX IF EXISTS "auth"."oauth_auth_pending_exp_idx";
DROP INDEX IF EXISTS "auth"."mfa_factors_user_id_idx";
DROP INDEX IF EXISTS "auth"."mfa_factors_user_friendly_name_unique";
DROP INDEX IF EXISTS "auth"."mfa_challenge_created_at_idx";
DROP INDEX IF EXISTS "auth"."idx_user_id_auth_method";
DROP INDEX IF EXISTS "auth"."idx_auth_code";
DROP INDEX IF EXISTS "auth"."identities_user_id_idx";
DROP INDEX IF EXISTS "auth"."identities_email_idx";
DROP INDEX IF EXISTS "auth"."flow_state_created_at_idx";
DROP INDEX IF EXISTS "auth"."factor_id_created_at_idx";
DROP INDEX IF EXISTS "auth"."email_change_token_new_idx";
DROP INDEX IF EXISTS "auth"."email_change_token_current_idx";
DROP INDEX IF EXISTS "auth"."confirmation_token_idx";
DROP INDEX IF EXISTS "auth"."audit_logs_instance_id_idx";
DROP INDEX IF EXISTS "_realtime"."tenants_external_id_index";
DROP INDEX IF EXISTS "_realtime"."extensions_tenant_external_id_type_index";
DROP INDEX IF EXISTS "_realtime"."extensions_tenant_external_id_index";
ALTER TABLE IF EXISTS ONLY "supabase_functions"."migrations" DROP CONSTRAINT IF EXISTS "migrations_pkey";
ALTER TABLE IF EXISTS ONLY "supabase_functions"."hooks" DROP CONSTRAINT IF EXISTS "hooks_pkey";
ALTER TABLE IF EXISTS ONLY "storage"."s3_multipart_uploads" DROP CONSTRAINT IF EXISTS "s3_multipart_uploads_pkey";
ALTER TABLE IF EXISTS ONLY "storage"."s3_multipart_uploads_parts" DROP CONSTRAINT IF EXISTS "s3_multipart_uploads_parts_pkey";
ALTER TABLE IF EXISTS ONLY "storage"."prefixes" DROP CONSTRAINT IF EXISTS "prefixes_pkey";
ALTER TABLE IF EXISTS ONLY "storage"."objects" DROP CONSTRAINT IF EXISTS "objects_pkey";
ALTER TABLE IF EXISTS ONLY "storage"."migrations" DROP CONSTRAINT IF EXISTS "migrations_pkey";
ALTER TABLE IF EXISTS ONLY "storage"."migrations" DROP CONSTRAINT IF EXISTS "migrations_name_key";
ALTER TABLE IF EXISTS ONLY "storage"."iceberg_tables" DROP CONSTRAINT IF EXISTS "iceberg_tables_pkey";
ALTER TABLE IF EXISTS ONLY "storage"."iceberg_namespaces" DROP CONSTRAINT IF EXISTS "iceberg_namespaces_pkey";
ALTER TABLE IF EXISTS ONLY "storage"."buckets" DROP CONSTRAINT IF EXISTS "buckets_pkey";
ALTER TABLE IF EXISTS ONLY "storage"."buckets_analytics" DROP CONSTRAINT IF EXISTS "buckets_analytics_pkey";
ALTER TABLE IF EXISTS ONLY "realtime"."schema_migrations" DROP CONSTRAINT IF EXISTS "schema_migrations_pkey";
ALTER TABLE IF EXISTS ONLY "realtime"."subscription" DROP CONSTRAINT IF EXISTS "pk_subscription";
ALTER TABLE IF EXISTS ONLY "realtime"."messages_2025_10_10" DROP CONSTRAINT IF EXISTS "messages_2025_10_10_pkey";
ALTER TABLE IF EXISTS ONLY "realtime"."messages_2025_10_09" DROP CONSTRAINT IF EXISTS "messages_2025_10_09_pkey";
ALTER TABLE IF EXISTS ONLY "realtime"."messages_2025_10_08" DROP CONSTRAINT IF EXISTS "messages_2025_10_08_pkey";
ALTER TABLE IF EXISTS ONLY "realtime"."messages_2025_10_07" DROP CONSTRAINT IF EXISTS "messages_2025_10_07_pkey";
ALTER TABLE IF EXISTS ONLY "realtime"."messages_2025_10_06" DROP CONSTRAINT IF EXISTS "messages_2025_10_06_pkey";
ALTER TABLE IF EXISTS ONLY "realtime"."messages_2025_10_05" DROP CONSTRAINT IF EXISTS "messages_2025_10_05_pkey";
ALTER TABLE IF EXISTS ONLY "realtime"."messages_2025_10_04" DROP CONSTRAINT IF EXISTS "messages_2025_10_04_pkey";
ALTER TABLE IF EXISTS ONLY "realtime"."messages" DROP CONSTRAINT IF EXISTS "messages_pkey";
ALTER TABLE IF EXISTS ONLY "public"."webhook_events" DROP CONSTRAINT IF EXISTS "webhook_events_pkey";
ALTER TABLE IF EXISTS ONLY "public"."webhook_events" DROP CONSTRAINT IF EXISTS "webhook_events_event_id_key";
ALTER TABLE IF EXISTS ONLY "public"."viewing_sessions" DROP CONSTRAINT IF EXISTS "viewing_sessions_pkey";
ALTER TABLE IF EXISTS ONLY "public"."videos" DROP CONSTRAINT IF EXISTS "videos_pkey";
ALTER TABLE IF EXISTS ONLY "public"."videos" DROP CONSTRAINT IF EXISTS "videos_cloudflare_uid_key";
ALTER TABLE IF EXISTS ONLY "public"."video_qualities" DROP CONSTRAINT IF EXISTS "video_qualities_pkey";
ALTER TABLE IF EXISTS ONLY "public"."video_permissions" DROP CONSTRAINT IF EXISTS "video_permissions_pkey";
ALTER TABLE IF EXISTS ONLY "public"."video_analytics" DROP CONSTRAINT IF EXISTS "video_analytics_video_id_date_key";
ALTER TABLE IF EXISTS ONLY "public"."video_analytics" DROP CONSTRAINT IF EXISTS "video_analytics_pkey";
ALTER TABLE IF EXISTS ONLY "public"."users" DROP CONSTRAINT IF EXISTS "users_username_key";
ALTER TABLE IF EXISTS ONLY "public"."users" DROP CONSTRAINT IF EXISTS "users_pkey";
ALTER TABLE IF EXISTS ONLY "public"."users" DROP CONSTRAINT IF EXISTS "users_email_key";
ALTER TABLE IF EXISTS ONLY "public"."user_payment_methods" DROP CONSTRAINT IF EXISTS "user_payment_methods_pkey";
ALTER TABLE IF EXISTS ONLY "public"."upload_sessions" DROP CONSTRAINT IF EXISTS "upload_sessions_pkey";
ALTER TABLE IF EXISTS ONLY "public"."transactions" DROP CONSTRAINT IF EXISTS "uk_transactions_stripe_payment_intent_id";
ALTER TABLE IF EXISTS ONLY "public"."lecture_resources" DROP CONSTRAINT IF EXISTS "uk_lecture_resources_lecture_file";
ALTER TABLE IF EXISTS ONLY "public"."transactions" DROP CONSTRAINT IF EXISTS "transactions_transaction_reference_key";
ALTER TABLE IF EXISTS ONLY "public"."transactions" DROP CONSTRAINT IF EXISTS "transactions_pkey";
ALTER TABLE IF EXISTS ONLY "public"."subscriptions" DROP CONSTRAINT IF EXISTS "subscriptions_pkey";
ALTER TABLE IF EXISTS ONLY "public"."stripe_webhook_events" DROP CONSTRAINT IF EXISTS "stripe_webhook_events_stripe_event_id_key";
ALTER TABLE IF EXISTS ONLY "public"."stripe_webhook_events" DROP CONSTRAINT IF EXISTS "stripe_webhook_events_pkey";
ALTER TABLE IF EXISTS ONLY "public"."stripe_products" DROP CONSTRAINT IF EXISTS "stripe_products_stripe_product_id_key";
ALTER TABLE IF EXISTS ONLY "public"."stripe_products" DROP CONSTRAINT IF EXISTS "stripe_products_stripe_price_id_key";
ALTER TABLE IF EXISTS ONLY "public"."stripe_products" DROP CONSTRAINT IF EXISTS "stripe_products_pkey";
ALTER TABLE IF EXISTS ONLY "public"."stripe_customers" DROP CONSTRAINT IF EXISTS "stripe_customers_stripe_customer_id_key";
ALTER TABLE IF EXISTS ONLY "public"."stripe_customers" DROP CONSTRAINT IF EXISTS "stripe_customers_pkey";
ALTER TABLE IF EXISTS ONLY "public"."schema_migrations" DROP CONSTRAINT IF EXISTS "schema_migrations_pkey";
ALTER TABLE IF EXISTS ONLY "public"."progress" DROP CONSTRAINT IF EXISTS "progress_pkey";
ALTER TABLE IF EXISTS ONLY "public"."payment_methods" DROP CONSTRAINT IF EXISTS "payment_methods_pkey";
ALTER TABLE IF EXISTS ONLY "public"."payment_events" DROP CONSTRAINT IF EXISTS "payment_events_provider_provider_event_id_key";
ALTER TABLE IF EXISTS ONLY "public"."payment_events" DROP CONSTRAINT IF EXISTS "payment_events_pkey";
ALTER TABLE IF EXISTS ONLY "public"."oauth_accounts" DROP CONSTRAINT IF EXISTS "oauth_accounts_provider_provider_id_key";
ALTER TABLE IF EXISTS ONLY "public"."oauth_accounts" DROP CONSTRAINT IF EXISTS "oauth_accounts_pkey";
ALTER TABLE IF EXISTS ONLY "public"."notes" DROP CONSTRAINT IF EXISTS "notes_pkey";
ALTER TABLE IF EXISTS ONLY "public"."network_metrics" DROP CONSTRAINT IF EXISTS "network_metrics_pkey";
ALTER TABLE IF EXISTS ONLY "public"."network_events" DROP CONSTRAINT IF EXISTS "network_events_pkey";
ALTER TABLE IF EXISTS ONLY "public"."network_analytics_daily" DROP CONSTRAINT IF EXISTS "network_analytics_daily_pkey";
ALTER TABLE IF EXISTS ONLY "public"."network_analytics_daily" DROP CONSTRAINT IF EXISTS "network_analytics_daily_date_user_id_video_id_key";
ALTER TABLE IF EXISTS ONLY "public"."lemon_squeezy_webhook_events" DROP CONSTRAINT IF EXISTS "lemon_squeezy_webhook_events_pkey";
ALTER TABLE IF EXISTS ONLY "public"."lemon_squeezy_webhook_events" DROP CONSTRAINT IF EXISTS "lemon_squeezy_webhook_events_event_id_key";
ALTER TABLE IF EXISTS ONLY "public"."lemon_squeezy_variants" DROP CONSTRAINT IF EXISTS "lemon_squeezy_variants_pkey";
ALTER TABLE IF EXISTS ONLY "public"."lemon_squeezy_variants" DROP CONSTRAINT IF EXISTS "lemon_squeezy_variants_lemon_squeezy_variant_id_key";
ALTER TABLE IF EXISTS ONLY "public"."lemon_squeezy_products" DROP CONSTRAINT IF EXISTS "lemon_squeezy_products_pkey";
ALTER TABLE IF EXISTS ONLY "public"."lemon_squeezy_products" DROP CONSTRAINT IF EXISTS "lemon_squeezy_products_lemon_squeezy_product_id_key";
ALTER TABLE IF EXISTS ONLY "public"."lectures" DROP CONSTRAINT IF EXISTS "lectures_pkey";
ALTER TABLE IF EXISTS ONLY "public"."lecture_resources" DROP CONSTRAINT IF EXISTS "lecture_resources_pkey";
ALTER TABLE IF EXISTS ONLY "public"."lecture_preview_sessions" DROP CONSTRAINT IF EXISTS "lecture_preview_sessions_user_id_lecture_id_key";
ALTER TABLE IF EXISTS ONLY "public"."lecture_preview_sessions" DROP CONSTRAINT IF EXISTS "lecture_preview_sessions_pkey";
ALTER TABLE IF EXISTS ONLY "public"."forum_votes" DROP CONSTRAINT IF EXISTS "forum_votes_post_id_user_id_key";
ALTER TABLE IF EXISTS ONLY "public"."forum_votes" DROP CONSTRAINT IF EXISTS "forum_votes_pkey";
ALTER TABLE IF EXISTS ONLY "public"."forum_topics" DROP CONSTRAINT IF EXISTS "forum_topics_pkey";
ALTER TABLE IF EXISTS ONLY "public"."forum_topic_subscriptions" DROP CONSTRAINT IF EXISTS "forum_topic_subscriptions_user_id_topic_id_key";
ALTER TABLE IF EXISTS ONLY "public"."forum_topic_subscriptions" DROP CONSTRAINT IF EXISTS "forum_topic_subscriptions_pkey";
ALTER TABLE IF EXISTS ONLY "public"."forum_posts" DROP CONSTRAINT IF EXISTS "forum_posts_pkey";
ALTER TABLE IF EXISTS ONLY "public"."forum_notifications" DROP CONSTRAINT IF EXISTS "forum_notifications_pkey";
ALTER TABLE IF EXISTS ONLY "public"."forum_mentions" DROP CONSTRAINT IF EXISTS "forum_mentions_post_id_mentioned_user_id_key";
ALTER TABLE IF EXISTS ONLY "public"."forum_mentions" DROP CONSTRAINT IF EXISTS "forum_mentions_pkey";
ALTER TABLE IF EXISTS ONLY "public"."files" DROP CONSTRAINT IF EXISTS "files_pkey";
ALTER TABLE IF EXISTS ONLY "public"."file_permissions" DROP CONSTRAINT IF EXISTS "file_permissions_pkey";
ALTER TABLE IF EXISTS ONLY "public"."enrollments" DROP CONSTRAINT IF EXISTS "enrollments_user_id_course_id_key";
ALTER TABLE IF EXISTS ONLY "public"."enrollments" DROP CONSTRAINT IF EXISTS "enrollments_pkey";
ALTER TABLE IF EXISTS ONLY "public"."courses" DROP CONSTRAINT IF EXISTS "courses_pkey";
ALTER TABLE IF EXISTS ONLY "public"."course_resources" DROP CONSTRAINT IF EXISTS "course_resources_pkey";
ALTER TABLE IF EXISTS ONLY "public"."course_resources" DROP CONSTRAINT IF EXISTS "course_resources_course_id_file_id_key";
ALTER TABLE IF EXISTS ONLY "public"."course_access_logs" DROP CONSTRAINT IF EXISTS "course_access_logs_pkey";
ALTER TABLE IF EXISTS ONLY "public"."course_access_cache" DROP CONSTRAINT IF EXISTS "course_access_cache_user_id_course_id_key";
ALTER TABLE IF EXISTS ONLY "public"."course_access_cache" DROP CONSTRAINT IF EXISTS "course_access_cache_pkey";
ALTER TABLE IF EXISTS ONLY "public"."chat_history" DROP CONSTRAINT IF EXISTS "chat_history_pkey";
ALTER TABLE IF EXISTS ONLY "public"."bandwidth_tests" DROP CONSTRAINT IF EXISTS "bandwidth_tests_pkey";
ALTER TABLE IF EXISTS ONLY "public"."audit_logs" DROP CONSTRAINT IF EXISTS "audit_logs_pkey";
ALTER TABLE IF EXISTS ONLY "public"."adaptive_streaming_rules" DROP CONSTRAINT IF EXISTS "adaptive_streaming_rules_pkey";
ALTER TABLE IF EXISTS ONLY "auth"."users" DROP CONSTRAINT IF EXISTS "users_pkey";
ALTER TABLE IF EXISTS ONLY "auth"."users" DROP CONSTRAINT IF EXISTS "users_phone_key";
ALTER TABLE IF EXISTS ONLY "auth"."sso_providers" DROP CONSTRAINT IF EXISTS "sso_providers_pkey";
ALTER TABLE IF EXISTS ONLY "auth"."sso_domains" DROP CONSTRAINT IF EXISTS "sso_domains_pkey";
ALTER TABLE IF EXISTS ONLY "auth"."sessions" DROP CONSTRAINT IF EXISTS "sessions_pkey";
ALTER TABLE IF EXISTS ONLY "auth"."schema_migrations" DROP CONSTRAINT IF EXISTS "schema_migrations_pkey";
ALTER TABLE IF EXISTS ONLY "auth"."saml_relay_states" DROP CONSTRAINT IF EXISTS "saml_relay_states_pkey";
ALTER TABLE IF EXISTS ONLY "auth"."saml_providers" DROP CONSTRAINT IF EXISTS "saml_providers_pkey";
ALTER TABLE IF EXISTS ONLY "auth"."saml_providers" DROP CONSTRAINT IF EXISTS "saml_providers_entity_id_key";
ALTER TABLE IF EXISTS ONLY "auth"."refresh_tokens" DROP CONSTRAINT IF EXISTS "refresh_tokens_token_unique";
ALTER TABLE IF EXISTS ONLY "auth"."refresh_tokens" DROP CONSTRAINT IF EXISTS "refresh_tokens_pkey";
ALTER TABLE IF EXISTS ONLY "auth"."one_time_tokens" DROP CONSTRAINT IF EXISTS "one_time_tokens_pkey";
ALTER TABLE IF EXISTS ONLY "auth"."oauth_consents" DROP CONSTRAINT IF EXISTS "oauth_consents_user_client_unique";
ALTER TABLE IF EXISTS ONLY "auth"."oauth_consents" DROP CONSTRAINT IF EXISTS "oauth_consents_pkey";
ALTER TABLE IF EXISTS ONLY "auth"."oauth_clients" DROP CONSTRAINT IF EXISTS "oauth_clients_pkey";
ALTER TABLE IF EXISTS ONLY "auth"."oauth_authorizations" DROP CONSTRAINT IF EXISTS "oauth_authorizations_pkey";
ALTER TABLE IF EXISTS ONLY "auth"."oauth_authorizations" DROP CONSTRAINT IF EXISTS "oauth_authorizations_authorization_id_key";
ALTER TABLE IF EXISTS ONLY "auth"."oauth_authorizations" DROP CONSTRAINT IF EXISTS "oauth_authorizations_authorization_code_key";
ALTER TABLE IF EXISTS ONLY "auth"."mfa_factors" DROP CONSTRAINT IF EXISTS "mfa_factors_pkey";
ALTER TABLE IF EXISTS ONLY "auth"."mfa_factors" DROP CONSTRAINT IF EXISTS "mfa_factors_last_challenged_at_key";
ALTER TABLE IF EXISTS ONLY "auth"."mfa_challenges" DROP CONSTRAINT IF EXISTS "mfa_challenges_pkey";
ALTER TABLE IF EXISTS ONLY "auth"."mfa_amr_claims" DROP CONSTRAINT IF EXISTS "mfa_amr_claims_session_id_authentication_method_pkey";
ALTER TABLE IF EXISTS ONLY "auth"."instances" DROP CONSTRAINT IF EXISTS "instances_pkey";
ALTER TABLE IF EXISTS ONLY "auth"."identities" DROP CONSTRAINT IF EXISTS "identities_provider_id_provider_unique";
ALTER TABLE IF EXISTS ONLY "auth"."identities" DROP CONSTRAINT IF EXISTS "identities_pkey";
ALTER TABLE IF EXISTS ONLY "auth"."flow_state" DROP CONSTRAINT IF EXISTS "flow_state_pkey";
ALTER TABLE IF EXISTS ONLY "auth"."audit_log_entries" DROP CONSTRAINT IF EXISTS "audit_log_entries_pkey";
ALTER TABLE IF EXISTS ONLY "auth"."mfa_amr_claims" DROP CONSTRAINT IF EXISTS "amr_id_pk";
ALTER TABLE IF EXISTS ONLY "_realtime"."tenants" DROP CONSTRAINT IF EXISTS "tenants_pkey";
ALTER TABLE IF EXISTS ONLY "_realtime"."schema_migrations" DROP CONSTRAINT IF EXISTS "schema_migrations_pkey";
ALTER TABLE IF EXISTS ONLY "_realtime"."extensions" DROP CONSTRAINT IF EXISTS "extensions_pkey";
ALTER TABLE IF EXISTS "supabase_functions"."hooks" ALTER COLUMN "id" DROP DEFAULT;
ALTER TABLE IF EXISTS "auth"."refresh_tokens" ALTER COLUMN "id" DROP DEFAULT;
DROP TABLE IF EXISTS "supabase_functions"."migrations";
DROP SEQUENCE IF EXISTS "supabase_functions"."hooks_id_seq";
DROP TABLE IF EXISTS "supabase_functions"."hooks";
DROP TABLE IF EXISTS "storage"."s3_multipart_uploads_parts";
DROP TABLE IF EXISTS "storage"."s3_multipart_uploads";
DROP TABLE IF EXISTS "storage"."prefixes";
DROP TABLE IF EXISTS "storage"."objects";
DROP TABLE IF EXISTS "storage"."migrations";
DROP TABLE IF EXISTS "storage"."iceberg_tables";
DROP TABLE IF EXISTS "storage"."iceberg_namespaces";
DROP TABLE IF EXISTS "storage"."buckets_analytics";
DROP TABLE IF EXISTS "storage"."buckets";
DROP TABLE IF EXISTS "realtime"."subscription";
DROP TABLE IF EXISTS "realtime"."schema_migrations";
DROP TABLE IF EXISTS "realtime"."messages_2025_10_10";
DROP TABLE IF EXISTS "realtime"."messages_2025_10_09";
DROP TABLE IF EXISTS "realtime"."messages_2025_10_08";
DROP TABLE IF EXISTS "realtime"."messages_2025_10_07";
DROP TABLE IF EXISTS "realtime"."messages_2025_10_06";
DROP TABLE IF EXISTS "realtime"."messages_2025_10_05";
DROP TABLE IF EXISTS "realtime"."messages_2025_10_04";
DROP TABLE IF EXISTS "realtime"."messages";
DROP TABLE IF EXISTS "public"."webhook_events";
DROP TABLE IF EXISTS "public"."videos";
DROP TABLE IF EXISTS "public"."video_qualities";
DROP TABLE IF EXISTS "public"."video_permissions";
DROP TABLE IF EXISTS "public"."video_analytics";
DROP TABLE IF EXISTS "public"."users";
DROP TABLE IF EXISTS "public"."user_payment_methods";
DROP TABLE IF EXISTS "public"."upload_sessions";
DROP TABLE IF EXISTS "public"."subscriptions";
DROP TABLE IF EXISTS "public"."stripe_webhook_events";
DROP TABLE IF EXISTS "public"."stripe_products";
DROP TABLE IF EXISTS "public"."stripe_customers";
DROP TABLE IF EXISTS "public"."schema_migrations";
DROP TABLE IF EXISTS "public"."progress";
DROP TABLE IF EXISTS "public"."payment_methods";
DROP TABLE IF EXISTS "public"."payment_events";
DROP VIEW IF EXISTS "public"."payment_analytics";
DROP TABLE IF EXISTS "public"."oauth_accounts";
DROP TABLE IF EXISTS "public"."notes";
DROP MATERIALIZED VIEW IF EXISTS "public"."network_monitoring_dashboard";
DROP TABLE IF EXISTS "public"."viewing_sessions";
DROP TABLE IF EXISTS "public"."network_metrics";
DROP TABLE IF EXISTS "public"."network_events";
DROP TABLE IF EXISTS "public"."network_analytics_daily";
DROP TABLE IF EXISTS "public"."lemon_squeezy_webhook_events";
DROP TABLE IF EXISTS "public"."lemon_squeezy_variants";
DROP TABLE IF EXISTS "public"."lemon_squeezy_products";
DROP TABLE IF EXISTS "public"."lectures";
DROP TABLE IF EXISTS "public"."lecture_resources";
DROP TABLE IF EXISTS "public"."lecture_preview_sessions";
DROP TABLE IF EXISTS "public"."forum_votes";
DROP TABLE IF EXISTS "public"."forum_topics";
DROP TABLE IF EXISTS "public"."forum_topic_subscriptions";
DROP TABLE IF EXISTS "public"."forum_posts";
DROP TABLE IF EXISTS "public"."forum_notifications";
DROP TABLE IF EXISTS "public"."forum_mentions";
DROP TABLE IF EXISTS "public"."files";
DROP TABLE IF EXISTS "public"."file_permissions";
DROP TABLE IF EXISTS "public"."course_resources";
DROP TABLE IF EXISTS "public"."course_access_cache";
DROP VIEW IF EXISTS "public"."course_access_analytics";
DROP TABLE IF EXISTS "public"."transactions";
DROP TABLE IF EXISTS "public"."enrollments";
DROP TABLE IF EXISTS "public"."courses";
DROP TABLE IF EXISTS "public"."course_access_logs";
DROP TABLE IF EXISTS "public"."chat_history";
DROP TABLE IF EXISTS "public"."bandwidth_tests";
DROP TABLE IF EXISTS "public"."audit_logs";
DROP TABLE IF EXISTS "public"."adaptive_streaming_rules";
DROP TABLE IF EXISTS "auth"."users";
DROP TABLE IF EXISTS "auth"."sso_providers";
DROP TABLE IF EXISTS "auth"."sso_domains";
DROP TABLE IF EXISTS "auth"."sessions";
DROP TABLE IF EXISTS "auth"."schema_migrations";
DROP TABLE IF EXISTS "auth"."saml_relay_states";
DROP TABLE IF EXISTS "auth"."saml_providers";
DROP SEQUENCE IF EXISTS "auth"."refresh_tokens_id_seq";
DROP TABLE IF EXISTS "auth"."refresh_tokens";
DROP TABLE IF EXISTS "auth"."one_time_tokens";
DROP TABLE IF EXISTS "auth"."oauth_consents";
DROP TABLE IF EXISTS "auth"."oauth_clients";
DROP TABLE IF EXISTS "auth"."oauth_authorizations";
DROP TABLE IF EXISTS "auth"."mfa_factors";
DROP TABLE IF EXISTS "auth"."mfa_challenges";
DROP TABLE IF EXISTS "auth"."mfa_amr_claims";
DROP TABLE IF EXISTS "auth"."instances";
DROP TABLE IF EXISTS "auth"."identities";
DROP TABLE IF EXISTS "auth"."flow_state";
DROP TABLE IF EXISTS "auth"."audit_log_entries";
DROP TABLE IF EXISTS "_realtime"."tenants";
DROP TABLE IF EXISTS "_realtime"."schema_migrations";
DROP TABLE IF EXISTS "_realtime"."extensions";
DROP FUNCTION IF EXISTS "supabase_functions"."http_request"();
DROP FUNCTION IF EXISTS "storage"."update_updated_at_column"();
DROP FUNCTION IF EXISTS "storage"."search_v2"("prefix" "text", "bucket_name" "text", "limits" integer, "levels" integer, "start_after" "text", "sort_order" "text", "sort_column" "text", "sort_column_after" "text");
DROP FUNCTION IF EXISTS "storage"."search_v1_optimised"("prefix" "text", "bucketname" "text", "limits" integer, "levels" integer, "offsets" integer, "search" "text", "sortcolumn" "text", "sortorder" "text");
DROP FUNCTION IF EXISTS "storage"."search_legacy_v1"("prefix" "text", "bucketname" "text", "limits" integer, "levels" integer, "offsets" integer, "search" "text", "sortcolumn" "text", "sortorder" "text");
DROP FUNCTION IF EXISTS "storage"."search"("prefix" "text", "bucketname" "text", "limits" integer, "levels" integer, "offsets" integer, "search" "text", "sortcolumn" "text", "sortorder" "text");
DROP FUNCTION IF EXISTS "storage"."prefixes_insert_trigger"();
DROP FUNCTION IF EXISTS "storage"."prefixes_delete_cleanup"();
DROP FUNCTION IF EXISTS "storage"."operation"();
DROP FUNCTION IF EXISTS "storage"."objects_update_prefix_trigger"();
DROP FUNCTION IF EXISTS "storage"."objects_update_level_trigger"();
DROP FUNCTION IF EXISTS "storage"."objects_update_cleanup"();
DROP FUNCTION IF EXISTS "storage"."objects_insert_prefix_trigger"();
DROP FUNCTION IF EXISTS "storage"."objects_delete_cleanup"();
DROP FUNCTION IF EXISTS "storage"."lock_top_prefixes"("bucket_ids" "text"[], "names" "text"[]);
DROP FUNCTION IF EXISTS "storage"."list_objects_with_delimiter"("bucket_id" "text", "prefix_param" "text", "delimiter_param" "text", "max_keys" integer, "start_after" "text", "next_token" "text");
DROP FUNCTION IF EXISTS "storage"."list_multipart_uploads_with_delimiter"("bucket_id" "text", "prefix_param" "text", "delimiter_param" "text", "max_keys" integer, "next_key_token" "text", "next_upload_token" "text");
DROP FUNCTION IF EXISTS "storage"."get_size_by_bucket"();
DROP FUNCTION IF EXISTS "storage"."get_prefixes"("name" "text");
DROP FUNCTION IF EXISTS "storage"."get_prefix"("name" "text");
DROP FUNCTION IF EXISTS "storage"."get_level"("name" "text");
DROP FUNCTION IF EXISTS "storage"."foldername"("name" "text");
DROP FUNCTION IF EXISTS "storage"."filename"("name" "text");
DROP FUNCTION IF EXISTS "storage"."extension"("name" "text");
DROP FUNCTION IF EXISTS "storage"."enforce_bucket_name_length"();
DROP FUNCTION IF EXISTS "storage"."delete_prefix_hierarchy_trigger"();
DROP FUNCTION IF EXISTS "storage"."delete_prefix"("_bucket_id" "text", "_name" "text");
DROP FUNCTION IF EXISTS "storage"."delete_leaf_prefixes"("bucket_ids" "text"[], "names" "text"[]);
DROP FUNCTION IF EXISTS "storage"."can_insert_object"("bucketid" "text", "name" "text", "owner" "uuid", "metadata" "jsonb");
DROP FUNCTION IF EXISTS "storage"."add_prefixes"("_bucket_id" "text", "_name" "text");
DROP FUNCTION IF EXISTS "realtime"."topic"();
DROP FUNCTION IF EXISTS "realtime"."to_regrole"("role_name" "text");
DROP FUNCTION IF EXISTS "realtime"."subscription_check_filters"();
DROP FUNCTION IF EXISTS "realtime"."send"("payload" "jsonb", "event" "text", "topic" "text", "private" boolean);
DROP FUNCTION IF EXISTS "realtime"."quote_wal2json"("entity" "regclass");
DROP FUNCTION IF EXISTS "realtime"."list_changes"("publication" "name", "slot_name" "name", "max_changes" integer, "max_record_bytes" integer);
DROP FUNCTION IF EXISTS "realtime"."is_visible_through_filters"("columns" "realtime"."wal_column"[], "filters" "realtime"."user_defined_filter"[]);
DROP FUNCTION IF EXISTS "realtime"."check_equality_op"("op" "realtime"."equality_op", "type_" "regtype", "val_1" "text", "val_2" "text");
DROP FUNCTION IF EXISTS "realtime"."cast"("val" "text", "type_" "regtype");
DROP FUNCTION IF EXISTS "realtime"."build_prepared_statement_sql"("prepared_statement_name" "text", "entity" "regclass", "columns" "realtime"."wal_column"[]);
DROP FUNCTION IF EXISTS "realtime"."broadcast_changes"("topic_name" "text", "event_name" "text", "operation" "text", "table_name" "text", "table_schema" "text", "new" "record", "old" "record", "level" "text");
DROP FUNCTION IF EXISTS "realtime"."apply_rls"("wal" "jsonb", "max_record_bytes" integer);
DROP FUNCTION IF EXISTS "public"."update_updated_at_column"();
DROP FUNCTION IF EXISTS "public"."update_stripe_products_updated_at"();
DROP FUNCTION IF EXISTS "public"."update_stripe_customers_updated_at"();
DROP FUNCTION IF EXISTS "public"."update_notes_updated_at"();
DROP FUNCTION IF EXISTS "public"."update_network_analytics_daily"();
DROP FUNCTION IF EXISTS "public"."update_enrollment_payment_status"();
DROP FUNCTION IF EXISTS "public"."refresh_network_dashboard"();
DROP FUNCTION IF EXISTS "public"."prevent_course_resources_insert"();
DROP FUNCTION IF EXISTS "public"."cleanup_old_preview_sessions"();
DROP FUNCTION IF EXISTS "public"."cleanup_expired_access_cache"();
DROP FUNCTION IF EXISTS "pgbouncer"."get_auth"("p_usename" "text");
DROP FUNCTION IF EXISTS "extensions"."set_graphql_placeholder"();
DROP FUNCTION IF EXISTS "extensions"."pgrst_drop_watch"();
DROP FUNCTION IF EXISTS "extensions"."pgrst_ddl_watch"();
DROP FUNCTION IF EXISTS "extensions"."grant_pg_net_access"();
DROP FUNCTION IF EXISTS "extensions"."grant_pg_graphql_access"();
DROP FUNCTION IF EXISTS "extensions"."grant_pg_cron_access"();
DROP FUNCTION IF EXISTS "auth"."uid"();
DROP FUNCTION IF EXISTS "auth"."role"();
DROP FUNCTION IF EXISTS "auth"."jwt"();
DROP FUNCTION IF EXISTS "auth"."email"();
DROP TYPE IF EXISTS "storage"."buckettype";
DROP TYPE IF EXISTS "realtime"."wal_rls";
DROP TYPE IF EXISTS "realtime"."wal_column";
DROP TYPE IF EXISTS "realtime"."user_defined_filter";
DROP TYPE IF EXISTS "realtime"."equality_op";
DROP TYPE IF EXISTS "realtime"."action";
DROP TYPE IF EXISTS "public"."user_role";
DROP TYPE IF EXISTS "public"."lecture_status";
DROP TYPE IF EXISTS "public"."enrollment_status";
DROP TYPE IF EXISTS "public"."course_status";
DROP TYPE IF EXISTS "public"."course_level";
DROP TYPE IF EXISTS "auth"."one_time_token_type";
DROP TYPE IF EXISTS "auth"."oauth_response_type";
DROP TYPE IF EXISTS "auth"."oauth_registration_type";
DROP TYPE IF EXISTS "auth"."oauth_client_type";
DROP TYPE IF EXISTS "auth"."oauth_authorization_status";
DROP TYPE IF EXISTS "auth"."factor_type";
DROP TYPE IF EXISTS "auth"."factor_status";
DROP TYPE IF EXISTS "auth"."code_challenge_method";
DROP TYPE IF EXISTS "auth"."aal_level";
DROP EXTENSION IF EXISTS "uuid-ossp";
DROP EXTENSION IF EXISTS "supabase_vault";
DROP EXTENSION IF EXISTS "pgcrypto";
DROP EXTENSION IF EXISTS "pg_stat_statements";
DROP EXTENSION IF EXISTS "pg_graphql";
DROP SCHEMA IF EXISTS "vault";
DROP SCHEMA IF EXISTS "supabase_functions";
DROP SCHEMA IF EXISTS "storage";
DROP SCHEMA IF EXISTS "realtime";
DROP SCHEMA IF EXISTS "pgbouncer";
DROP EXTENSION IF EXISTS "pg_net";
DROP SCHEMA IF EXISTS "graphql_public";
DROP SCHEMA IF EXISTS "graphql";
DROP SCHEMA IF EXISTS "extensions";
DROP SCHEMA IF EXISTS "auth";
DROP SCHEMA IF EXISTS "_realtime";
--
-- Name: _realtime; Type: SCHEMA; Schema: -; Owner: -
--

CREATE SCHEMA "_realtime";


--
-- Name: auth; Type: SCHEMA; Schema: -; Owner: -
--

CREATE SCHEMA "auth";


--
-- Name: extensions; Type: SCHEMA; Schema: -; Owner: -
--

CREATE SCHEMA "extensions";


--
-- Name: graphql; Type: SCHEMA; Schema: -; Owner: -
--

CREATE SCHEMA "graphql";


--
-- Name: graphql_public; Type: SCHEMA; Schema: -; Owner: -
--

CREATE SCHEMA "graphql_public";


--
-- Name: pg_net; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS "pg_net" WITH SCHEMA "extensions";


--
-- Name: EXTENSION "pg_net"; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION "pg_net" IS 'Async HTTP';


--
-- Name: pgbouncer; Type: SCHEMA; Schema: -; Owner: -
--

CREATE SCHEMA "pgbouncer";


--
-- Name: SCHEMA "public"; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON SCHEMA "public" IS 'standard public schema';


--
-- Name: realtime; Type: SCHEMA; Schema: -; Owner: -
--

CREATE SCHEMA "realtime";


--
-- Name: storage; Type: SCHEMA; Schema: -; Owner: -
--

CREATE SCHEMA "storage";


--
-- Name: supabase_functions; Type: SCHEMA; Schema: -; Owner: -
--

CREATE SCHEMA "supabase_functions";


--
-- Name: vault; Type: SCHEMA; Schema: -; Owner: -
--

CREATE SCHEMA "vault";


--
-- Name: pg_graphql; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS "pg_graphql" WITH SCHEMA "graphql";


--
-- Name: EXTENSION "pg_graphql"; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION "pg_graphql" IS 'pg_graphql: GraphQL support';


--
-- Name: pg_stat_statements; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS "pg_stat_statements" WITH SCHEMA "extensions";


--
-- Name: EXTENSION "pg_stat_statements"; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION "pg_stat_statements" IS 'track planning and execution statistics of all SQL statements executed';


--
-- Name: pgcrypto; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS "pgcrypto" WITH SCHEMA "extensions";


--
-- Name: EXTENSION "pgcrypto"; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION "pgcrypto" IS 'cryptographic functions';


--
-- Name: supabase_vault; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS "supabase_vault" WITH SCHEMA "vault";


--
-- Name: EXTENSION "supabase_vault"; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION "supabase_vault" IS 'Supabase Vault Extension';


--
-- Name: uuid-ossp; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA "public";


--
-- Name: EXTENSION "uuid-ossp"; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION "uuid-ossp" IS 'generate universally unique identifiers (UUIDs)';


--
-- Name: aal_level; Type: TYPE; Schema: auth; Owner: -
--

CREATE TYPE "auth"."aal_level" AS ENUM (
    'aal1',
    'aal2',
    'aal3'
);


--
-- Name: code_challenge_method; Type: TYPE; Schema: auth; Owner: -
--

CREATE TYPE "auth"."code_challenge_method" AS ENUM (
    's256',
    'plain'
);


--
-- Name: factor_status; Type: TYPE; Schema: auth; Owner: -
--

CREATE TYPE "auth"."factor_status" AS ENUM (
    'unverified',
    'verified'
);


--
-- Name: factor_type; Type: TYPE; Schema: auth; Owner: -
--

CREATE TYPE "auth"."factor_type" AS ENUM (
    'totp',
    'webauthn',
    'phone'
);


--
-- Name: oauth_authorization_status; Type: TYPE; Schema: auth; Owner: -
--

CREATE TYPE "auth"."oauth_authorization_status" AS ENUM (
    'pending',
    'approved',
    'denied',
    'expired'
);


--
-- Name: oauth_client_type; Type: TYPE; Schema: auth; Owner: -
--

CREATE TYPE "auth"."oauth_client_type" AS ENUM (
    'public',
    'confidential'
);


--
-- Name: oauth_registration_type; Type: TYPE; Schema: auth; Owner: -
--

CREATE TYPE "auth"."oauth_registration_type" AS ENUM (
    'dynamic',
    'manual'
);


--
-- Name: oauth_response_type; Type: TYPE; Schema: auth; Owner: -
--

CREATE TYPE "auth"."oauth_response_type" AS ENUM (
    'code'
);


--
-- Name: one_time_token_type; Type: TYPE; Schema: auth; Owner: -
--

CREATE TYPE "auth"."one_time_token_type" AS ENUM (
    'confirmation_token',
    'reauthentication_token',
    'recovery_token',
    'email_change_token_new',
    'email_change_token_current',
    'phone_change_token'
);


--
-- Name: course_level; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE "public"."course_level" AS ENUM (
    'beginner',
    'intermediate',
    'advanced'
);


--
-- Name: course_status; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE "public"."course_status" AS ENUM (
    'draft',
    'published',
    'archived'
);


--
-- Name: enrollment_status; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE "public"."enrollment_status" AS ENUM (
    'enrolled',
    'completed',
    'cancelled',
    'active',
    'pending'
);


--
-- Name: lecture_status; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE "public"."lecture_status" AS ENUM (
    'draft',
    'published'
);


--
-- Name: user_role; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE "public"."user_role" AS ENUM (
    'admin',
    'instructor',
    'student'
);


--
-- Name: action; Type: TYPE; Schema: realtime; Owner: -
--

CREATE TYPE "realtime"."action" AS ENUM (
    'INSERT',
    'UPDATE',
    'DELETE',
    'TRUNCATE',
    'ERROR'
);


--
-- Name: equality_op; Type: TYPE; Schema: realtime; Owner: -
--

CREATE TYPE "realtime"."equality_op" AS ENUM (
    'eq',
    'neq',
    'lt',
    'lte',
    'gt',
    'gte',
    'in'
);


--
-- Name: user_defined_filter; Type: TYPE; Schema: realtime; Owner: -
--

CREATE TYPE "realtime"."user_defined_filter" AS (
	"column_name" "text",
	"op" "realtime"."equality_op",
	"value" "text"
);


--
-- Name: wal_column; Type: TYPE; Schema: realtime; Owner: -
--

CREATE TYPE "realtime"."wal_column" AS (
	"name" "text",
	"type_name" "text",
	"type_oid" "oid",
	"value" "jsonb",
	"is_pkey" boolean,
	"is_selectable" boolean
);


--
-- Name: wal_rls; Type: TYPE; Schema: realtime; Owner: -
--

CREATE TYPE "realtime"."wal_rls" AS (
	"wal" "jsonb",
	"is_rls_enabled" boolean,
	"subscription_ids" "uuid"[],
	"errors" "text"[]
);


--
-- Name: buckettype; Type: TYPE; Schema: storage; Owner: -
--

CREATE TYPE "storage"."buckettype" AS ENUM (
    'STANDARD',
    'ANALYTICS'
);


--
-- Name: email(); Type: FUNCTION; Schema: auth; Owner: -
--

CREATE FUNCTION "auth"."email"() RETURNS "text"
    LANGUAGE "sql" STABLE
    AS $$
  select 
  coalesce(
    nullif(current_setting('request.jwt.claim.email', true), ''),
    (nullif(current_setting('request.jwt.claims', true), '')::jsonb ->> 'email')
  )::text
$$;


--
-- Name: FUNCTION "email"(); Type: COMMENT; Schema: auth; Owner: -
--

COMMENT ON FUNCTION "auth"."email"() IS 'Deprecated. Use auth.jwt() -> ''email'' instead.';


--
-- Name: jwt(); Type: FUNCTION; Schema: auth; Owner: -
--

CREATE FUNCTION "auth"."jwt"() RETURNS "jsonb"
    LANGUAGE "sql" STABLE
    AS $$
  select 
    coalesce(
        nullif(current_setting('request.jwt.claim', true), ''),
        nullif(current_setting('request.jwt.claims', true), '')
    )::jsonb
$$;


--
-- Name: role(); Type: FUNCTION; Schema: auth; Owner: -
--

CREATE FUNCTION "auth"."role"() RETURNS "text"
    LANGUAGE "sql" STABLE
    AS $$
  select 
  coalesce(
    nullif(current_setting('request.jwt.claim.role', true), ''),
    (nullif(current_setting('request.jwt.claims', true), '')::jsonb ->> 'role')
  )::text
$$;


--
-- Name: FUNCTION "role"(); Type: COMMENT; Schema: auth; Owner: -
--

COMMENT ON FUNCTION "auth"."role"() IS 'Deprecated. Use auth.jwt() -> ''role'' instead.';


--
-- Name: uid(); Type: FUNCTION; Schema: auth; Owner: -
--

CREATE FUNCTION "auth"."uid"() RETURNS "uuid"
    LANGUAGE "sql" STABLE
    AS $$
  select 
  coalesce(
    nullif(current_setting('request.jwt.claim.sub', true), ''),
    (nullif(current_setting('request.jwt.claims', true), '')::jsonb ->> 'sub')
  )::uuid
$$;


--
-- Name: FUNCTION "uid"(); Type: COMMENT; Schema: auth; Owner: -
--

COMMENT ON FUNCTION "auth"."uid"() IS 'Deprecated. Use auth.jwt() -> ''sub'' instead.';


--
-- Name: grant_pg_cron_access(); Type: FUNCTION; Schema: extensions; Owner: -
--

CREATE FUNCTION "extensions"."grant_pg_cron_access"() RETURNS "event_trigger"
    LANGUAGE "plpgsql"
    AS $$
BEGIN
  IF EXISTS (
    SELECT
    FROM pg_event_trigger_ddl_commands() AS ev
    JOIN pg_extension AS ext
    ON ev.objid = ext.oid
    WHERE ext.extname = 'pg_cron'
  )
  THEN
    grant usage on schema cron to postgres with grant option;

    alter default privileges in schema cron grant all on tables to postgres with grant option;
    alter default privileges in schema cron grant all on functions to postgres with grant option;
    alter default privileges in schema cron grant all on sequences to postgres with grant option;

    alter default privileges for user supabase_admin in schema cron grant all
        on sequences to postgres with grant option;
    alter default privileges for user supabase_admin in schema cron grant all
        on tables to postgres with grant option;
    alter default privileges for user supabase_admin in schema cron grant all
        on functions to postgres with grant option;

    grant all privileges on all tables in schema cron to postgres with grant option;
    revoke all on table cron.job from postgres;
    grant select on table cron.job to postgres with grant option;
  END IF;
END;
$$;


--
-- Name: FUNCTION "grant_pg_cron_access"(); Type: COMMENT; Schema: extensions; Owner: -
--

COMMENT ON FUNCTION "extensions"."grant_pg_cron_access"() IS 'Grants access to pg_cron';


--
-- Name: grant_pg_graphql_access(); Type: FUNCTION; Schema: extensions; Owner: -
--

CREATE FUNCTION "extensions"."grant_pg_graphql_access"() RETURNS "event_trigger"
    LANGUAGE "plpgsql"
    AS $_$
DECLARE
    func_is_graphql_resolve bool;
BEGIN
    func_is_graphql_resolve = (
        SELECT n.proname = 'resolve'
        FROM pg_event_trigger_ddl_commands() AS ev
        LEFT JOIN pg_catalog.pg_proc AS n
        ON ev.objid = n.oid
    );

    IF func_is_graphql_resolve
    THEN
        -- Update public wrapper to pass all arguments through to the pg_graphql resolve func
        DROP FUNCTION IF EXISTS graphql_public.graphql;
        create or replace function graphql_public.graphql(
            "operationName" text default null,
            query text default null,
            variables jsonb default null,
            extensions jsonb default null
        )
            returns jsonb
            language sql
        as $$
            select graphql.resolve(
                query := query,
                variables := coalesce(variables, '{}'),
                "operationName" := "operationName",
                extensions := extensions
            );
        $$;

        -- This hook executes when `graphql.resolve` is created. That is not necessarily the last
        -- function in the extension so we need to grant permissions on existing entities AND
        -- update default permissions to any others that are created after `graphql.resolve`
        grant usage on schema graphql to postgres, anon, authenticated, service_role;
        grant select on all tables in schema graphql to postgres, anon, authenticated, service_role;
        grant execute on all functions in schema graphql to postgres, anon, authenticated, service_role;
        grant all on all sequences in schema graphql to postgres, anon, authenticated, service_role;
        alter default privileges in schema graphql grant all on tables to postgres, anon, authenticated, service_role;
        alter default privileges in schema graphql grant all on functions to postgres, anon, authenticated, service_role;
        alter default privileges in schema graphql grant all on sequences to postgres, anon, authenticated, service_role;

        -- Allow postgres role to allow granting usage on graphql and graphql_public schemas to custom roles
        grant usage on schema graphql_public to postgres with grant option;
        grant usage on schema graphql to postgres with grant option;
    END IF;

END;
$_$;


--
-- Name: FUNCTION "grant_pg_graphql_access"(); Type: COMMENT; Schema: extensions; Owner: -
--

COMMENT ON FUNCTION "extensions"."grant_pg_graphql_access"() IS 'Grants access to pg_graphql';


--
-- Name: grant_pg_net_access(); Type: FUNCTION; Schema: extensions; Owner: -
--

CREATE FUNCTION "extensions"."grant_pg_net_access"() RETURNS "event_trigger"
    LANGUAGE "plpgsql"
    AS $$
BEGIN
  IF EXISTS (
    SELECT 1
    FROM pg_event_trigger_ddl_commands() AS ev
    JOIN pg_extension AS ext
    ON ev.objid = ext.oid
    WHERE ext.extname = 'pg_net'
  )
  THEN
    GRANT USAGE ON SCHEMA net TO supabase_functions_admin, postgres, anon, authenticated, service_role;

    ALTER function net.http_get(url text, params jsonb, headers jsonb, timeout_milliseconds integer) SECURITY DEFINER;
    ALTER function net.http_post(url text, body jsonb, params jsonb, headers jsonb, timeout_milliseconds integer) SECURITY DEFINER;

    ALTER function net.http_get(url text, params jsonb, headers jsonb, timeout_milliseconds integer) SET search_path = net;
    ALTER function net.http_post(url text, body jsonb, params jsonb, headers jsonb, timeout_milliseconds integer) SET search_path = net;

    REVOKE ALL ON FUNCTION net.http_get(url text, params jsonb, headers jsonb, timeout_milliseconds integer) FROM PUBLIC;
    REVOKE ALL ON FUNCTION net.http_post(url text, body jsonb, params jsonb, headers jsonb, timeout_milliseconds integer) FROM PUBLIC;

    GRANT EXECUTE ON FUNCTION net.http_get(url text, params jsonb, headers jsonb, timeout_milliseconds integer) TO supabase_functions_admin, postgres, anon, authenticated, service_role;
    GRANT EXECUTE ON FUNCTION net.http_post(url text, body jsonb, params jsonb, headers jsonb, timeout_milliseconds integer) TO supabase_functions_admin, postgres, anon, authenticated, service_role;
  END IF;
END;
$$;


--
-- Name: FUNCTION "grant_pg_net_access"(); Type: COMMENT; Schema: extensions; Owner: -
--

COMMENT ON FUNCTION "extensions"."grant_pg_net_access"() IS 'Grants access to pg_net';


--
-- Name: pgrst_ddl_watch(); Type: FUNCTION; Schema: extensions; Owner: -
--

CREATE FUNCTION "extensions"."pgrst_ddl_watch"() RETURNS "event_trigger"
    LANGUAGE "plpgsql"
    AS $$
DECLARE
  cmd record;
BEGIN
  FOR cmd IN SELECT * FROM pg_event_trigger_ddl_commands()
  LOOP
    IF cmd.command_tag IN (
      'CREATE SCHEMA', 'ALTER SCHEMA'
    , 'CREATE TABLE', 'CREATE TABLE AS', 'SELECT INTO', 'ALTER TABLE'
    , 'CREATE FOREIGN TABLE', 'ALTER FOREIGN TABLE'
    , 'CREATE VIEW', 'ALTER VIEW'
    , 'CREATE MATERIALIZED VIEW', 'ALTER MATERIALIZED VIEW'
    , 'CREATE FUNCTION', 'ALTER FUNCTION'
    , 'CREATE TRIGGER'
    , 'CREATE TYPE', 'ALTER TYPE'
    , 'CREATE RULE'
    , 'COMMENT'
    )
    -- don't notify in case of CREATE TEMP table or other objects created on pg_temp
    AND cmd.schema_name is distinct from 'pg_temp'
    THEN
      NOTIFY pgrst, 'reload schema';
    END IF;
  END LOOP;
END; $$;


--
-- Name: pgrst_drop_watch(); Type: FUNCTION; Schema: extensions; Owner: -
--

CREATE FUNCTION "extensions"."pgrst_drop_watch"() RETURNS "event_trigger"
    LANGUAGE "plpgsql"
    AS $$
DECLARE
  obj record;
BEGIN
  FOR obj IN SELECT * FROM pg_event_trigger_dropped_objects()
  LOOP
    IF obj.object_type IN (
      'schema'
    , 'table'
    , 'foreign table'
    , 'view'
    , 'materialized view'
    , 'function'
    , 'trigger'
    , 'type'
    , 'rule'
    )
    AND obj.is_temporary IS false -- no pg_temp objects
    THEN
      NOTIFY pgrst, 'reload schema';
    END IF;
  END LOOP;
END; $$;


--
-- Name: set_graphql_placeholder(); Type: FUNCTION; Schema: extensions; Owner: -
--

CREATE FUNCTION "extensions"."set_graphql_placeholder"() RETURNS "event_trigger"
    LANGUAGE "plpgsql"
    AS $_$
    DECLARE
    graphql_is_dropped bool;
    BEGIN
    graphql_is_dropped = (
        SELECT ev.schema_name = 'graphql_public'
        FROM pg_event_trigger_dropped_objects() AS ev
        WHERE ev.schema_name = 'graphql_public'
    );

    IF graphql_is_dropped
    THEN
        create or replace function graphql_public.graphql(
            "operationName" text default null,
            query text default null,
            variables jsonb default null,
            extensions jsonb default null
        )
            returns jsonb
            language plpgsql
        as $$
            DECLARE
                server_version float;
            BEGIN
                server_version = (SELECT (SPLIT_PART((select version()), ' ', 2))::float);

                IF server_version >= 14 THEN
                    RETURN jsonb_build_object(
                        'errors', jsonb_build_array(
                            jsonb_build_object(
                                'message', 'pg_graphql extension is not enabled.'
                            )
                        )
                    );
                ELSE
                    RETURN jsonb_build_object(
                        'errors', jsonb_build_array(
                            jsonb_build_object(
                                'message', 'pg_graphql is only available on projects running Postgres 14 onwards.'
                            )
                        )
                    );
                END IF;
            END;
        $$;
    END IF;

    END;
$_$;


--
-- Name: FUNCTION "set_graphql_placeholder"(); Type: COMMENT; Schema: extensions; Owner: -
--

COMMENT ON FUNCTION "extensions"."set_graphql_placeholder"() IS 'Reintroduces placeholder function for graphql_public.graphql';


--
-- Name: get_auth("text"); Type: FUNCTION; Schema: pgbouncer; Owner: -
--

CREATE FUNCTION "pgbouncer"."get_auth"("p_usename" "text") RETURNS TABLE("username" "text", "password" "text")
    LANGUAGE "plpgsql" SECURITY DEFINER
    AS $_$
begin
    raise debug 'PgBouncer auth request: %', p_usename;

    return query
    select 
        rolname::text, 
        case when rolvaliduntil < now() 
            then null 
            else rolpassword::text 
        end 
    from pg_authid 
    where rolname=$1 and rolcanlogin;
end;
$_$;


--
-- Name: cleanup_expired_access_cache(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION "public"."cleanup_expired_access_cache"() RETURNS "void"
    LANGUAGE "plpgsql"
    AS $$
BEGIN
    DELETE FROM course_access_cache WHERE expires_at < NOW();
END;
$$;


--
-- Name: cleanup_old_preview_sessions(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION "public"."cleanup_old_preview_sessions"() RETURNS "void"
    LANGUAGE "plpgsql"
    AS $$
BEGIN
    -- Mark sessions as exhausted if they've exceeded their time limit
    UPDATE lecture_preview_sessions
    SET preview_exhausted = TRUE,
        updated_at = NOW()
    WHERE preview_exhausted = FALSE
      AND (EXTRACT(EPOCH FROM (NOW() - session_started_at)) > preview_limit_seconds);

    -- Clean up old exhausted sessions (older than 30 days)
    DELETE FROM lecture_preview_sessions
    WHERE preview_exhausted = TRUE
      AND created_at < NOW() - INTERVAL '30 days';
END;
$$;


--
-- Name: prevent_course_resources_insert(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION "public"."prevent_course_resources_insert"() RETURNS "trigger"
    LANGUAGE "plpgsql"
    AS $$
BEGIN
    RAISE EXCEPTION 'Resources should now be attached to lectures, not courses. Use lecture_resources table instead.';
END;
$$;


--
-- Name: refresh_network_dashboard(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION "public"."refresh_network_dashboard"() RETURNS "void"
    LANGUAGE "plpgsql"
    AS $$
BEGIN
    REFRESH MATERIALIZED VIEW CONCURRENTLY network_monitoring_dashboard;
END;
$$;


--
-- Name: update_enrollment_payment_status(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION "public"."update_enrollment_payment_status"() RETURNS "trigger"
    LANGUAGE "plpgsql"
    AS $$
BEGIN
    -- When a transaction is completed, update corresponding enrollment
    IF NEW.status = 'completed' AND OLD.status != 'completed' AND NEW.course_id IS NOT NULL THEN
        UPDATE enrollments
        SET payment_status = 'paid',
            payment_verified_at = NEW.payment_verified_at,
            transaction_id = NEW.id,
            updated_at = NOW()
        WHERE user_id = NEW.user_id AND course_id = NEW.course_id;

        -- Clear access cache for this user/course combination
        DELETE FROM course_access_cache
        WHERE user_id = NEW.user_id AND course_id = NEW.course_id;
    END IF;

    -- When a transaction is refunded, update enrollment status
    IF NEW.status = 'refunded' AND OLD.status != 'refunded' AND NEW.course_id IS NOT NULL THEN
        UPDATE enrollments
        SET payment_status = 'refunded',
            updated_at = NOW()
        WHERE user_id = NEW.user_id AND course_id = NEW.course_id;

        -- Clear access cache
        DELETE FROM course_access_cache
        WHERE user_id = NEW.user_id AND course_id = NEW.course_id;
    END IF;

    RETURN NEW;
END;
$$;


--
-- Name: update_network_analytics_daily(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION "public"."update_network_analytics_daily"() RETURNS "trigger"
    LANGUAGE "plpgsql"
    AS $$
BEGIN
    INSERT INTO network_analytics_daily (
        date, user_id, video_id, session_count, avg_bandwidth_mbps,
        avg_latency_ms, avg_packet_loss, total_quality_changes
    )
    VALUES (
        CURRENT_DATE, NEW.user_id,
        (SELECT video_id FROM viewing_sessions WHERE session_id = NEW.session_id LIMIT 1),
        1, NEW.bandwidth_mbps, NEW.latency_ms, NEW.packet_loss_percent, 0
    )
    ON CONFLICT (date, user_id, video_id)
    DO UPDATE SET
        session_count = network_analytics_daily.session_count + 1,
        avg_bandwidth_mbps = (network_analytics_daily.avg_bandwidth_mbps * network_analytics_daily.session_count + NEW.bandwidth_mbps) / (network_analytics_daily.session_count + 1),
        avg_latency_ms = (network_analytics_daily.avg_latency_ms * network_analytics_daily.session_count + NEW.latency_ms) / (network_analytics_daily.session_count + 1),
        avg_packet_loss = (network_analytics_daily.avg_packet_loss * network_analytics_daily.session_count + NEW.packet_loss_percent) / (network_analytics_daily.session_count + 1),
        updated_at = NOW();

    RETURN NEW;
END;
$$;


--
-- Name: update_notes_updated_at(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION "public"."update_notes_updated_at"() RETURNS "trigger"
    LANGUAGE "plpgsql"
    AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$;


--
-- Name: update_stripe_customers_updated_at(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION "public"."update_stripe_customers_updated_at"() RETURNS "trigger"
    LANGUAGE "plpgsql"
    AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$;


--
-- Name: update_stripe_products_updated_at(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION "public"."update_stripe_products_updated_at"() RETURNS "trigger"
    LANGUAGE "plpgsql"
    AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$;


--
-- Name: update_updated_at_column(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION "public"."update_updated_at_column"() RETURNS "trigger"
    LANGUAGE "plpgsql"
    AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$;


--
-- Name: apply_rls("jsonb", integer); Type: FUNCTION; Schema: realtime; Owner: -
--

CREATE FUNCTION "realtime"."apply_rls"("wal" "jsonb", "max_record_bytes" integer DEFAULT (1024 * 1024)) RETURNS SETOF "realtime"."wal_rls"
    LANGUAGE "plpgsql"
    AS $$
declare
-- Regclass of the table e.g. public.notes
entity_ regclass = (quote_ident(wal ->> 'schema') || '.' || quote_ident(wal ->> 'table'))::regclass;

-- I, U, D, T: insert, update ...
action realtime.action = (
    case wal ->> 'action'
        when 'I' then 'INSERT'
        when 'U' then 'UPDATE'
        when 'D' then 'DELETE'
        else 'ERROR'
    end
);

-- Is row level security enabled for the table
is_rls_enabled bool = relrowsecurity from pg_class where oid = entity_;

subscriptions realtime.subscription[] = array_agg(subs)
    from
        realtime.subscription subs
    where
        subs.entity = entity_;

-- Subscription vars
roles regrole[] = array_agg(distinct us.claims_role::text)
    from
        unnest(subscriptions) us;

working_role regrole;
claimed_role regrole;
claims jsonb;

subscription_id uuid;
subscription_has_access bool;
visible_to_subscription_ids uuid[] = '{}';

-- structured info for wal's columns
columns realtime.wal_column[];
-- previous identity values for update/delete
old_columns realtime.wal_column[];

error_record_exceeds_max_size boolean = octet_length(wal::text) > max_record_bytes;

-- Primary jsonb output for record
output jsonb;

begin
perform set_config('role', null, true);

columns =
    array_agg(
        (
            x->>'name',
            x->>'type',
            x->>'typeoid',
            realtime.cast(
                (x->'value') #>> '{}',
                coalesce(
                    (x->>'typeoid')::regtype, -- null when wal2json version <= 2.4
                    (x->>'type')::regtype
                )
            ),
            (pks ->> 'name') is not null,
            true
        )::realtime.wal_column
    )
    from
        jsonb_array_elements(wal -> 'columns') x
        left join jsonb_array_elements(wal -> 'pk') pks
            on (x ->> 'name') = (pks ->> 'name');

old_columns =
    array_agg(
        (
            x->>'name',
            x->>'type',
            x->>'typeoid',
            realtime.cast(
                (x->'value') #>> '{}',
                coalesce(
                    (x->>'typeoid')::regtype, -- null when wal2json version <= 2.4
                    (x->>'type')::regtype
                )
            ),
            (pks ->> 'name') is not null,
            true
        )::realtime.wal_column
    )
    from
        jsonb_array_elements(wal -> 'identity') x
        left join jsonb_array_elements(wal -> 'pk') pks
            on (x ->> 'name') = (pks ->> 'name');

for working_role in select * from unnest(roles) loop

    -- Update `is_selectable` for columns and old_columns
    columns =
        array_agg(
            (
                c.name,
                c.type_name,
                c.type_oid,
                c.value,
                c.is_pkey,
                pg_catalog.has_column_privilege(working_role, entity_, c.name, 'SELECT')
            )::realtime.wal_column
        )
        from
            unnest(columns) c;

    old_columns =
            array_agg(
                (
                    c.name,
                    c.type_name,
                    c.type_oid,
                    c.value,
                    c.is_pkey,
                    pg_catalog.has_column_privilege(working_role, entity_, c.name, 'SELECT')
                )::realtime.wal_column
            )
            from
                unnest(old_columns) c;

    if action <> 'DELETE' and count(1) = 0 from unnest(columns) c where c.is_pkey then
        return next (
            jsonb_build_object(
                'schema', wal ->> 'schema',
                'table', wal ->> 'table',
                'type', action
            ),
            is_rls_enabled,
            -- subscriptions is already filtered by entity
            (select array_agg(s.subscription_id) from unnest(subscriptions) as s where claims_role = working_role),
            array['Error 400: Bad Request, no primary key']
        )::realtime.wal_rls;

    -- The claims role does not have SELECT permission to the primary key of entity
    elsif action <> 'DELETE' and sum(c.is_selectable::int) <> count(1) from unnest(columns) c where c.is_pkey then
        return next (
            jsonb_build_object(
                'schema', wal ->> 'schema',
                'table', wal ->> 'table',
                'type', action
            ),
            is_rls_enabled,
            (select array_agg(s.subscription_id) from unnest(subscriptions) as s where claims_role = working_role),
            array['Error 401: Unauthorized']
        )::realtime.wal_rls;

    else
        output = jsonb_build_object(
            'schema', wal ->> 'schema',
            'table', wal ->> 'table',
            'type', action,
            'commit_timestamp', to_char(
                ((wal ->> 'timestamp')::timestamptz at time zone 'utc'),
                'YYYY-MM-DD"T"HH24:MI:SS.MS"Z"'
            ),
            'columns', (
                select
                    jsonb_agg(
                        jsonb_build_object(
                            'name', pa.attname,
                            'type', pt.typname
                        )
                        order by pa.attnum asc
                    )
                from
                    pg_attribute pa
                    join pg_type pt
                        on pa.atttypid = pt.oid
                where
                    attrelid = entity_
                    and attnum > 0
                    and pg_catalog.has_column_privilege(working_role, entity_, pa.attname, 'SELECT')
            )
        )
        -- Add "record" key for insert and update
        || case
            when action in ('INSERT', 'UPDATE') then
                jsonb_build_object(
                    'record',
                    (
                        select
                            jsonb_object_agg(
                                -- if unchanged toast, get column name and value from old record
                                coalesce((c).name, (oc).name),
                                case
                                    when (c).name is null then (oc).value
                                    else (c).value
                                end
                            )
                        from
                            unnest(columns) c
                            full outer join unnest(old_columns) oc
                                on (c).name = (oc).name
                        where
                            coalesce((c).is_selectable, (oc).is_selectable)
                            and ( not error_record_exceeds_max_size or (octet_length((c).value::text) <= 64))
                    )
                )
            else '{}'::jsonb
        end
        -- Add "old_record" key for update and delete
        || case
            when action = 'UPDATE' then
                jsonb_build_object(
                        'old_record',
                        (
                            select jsonb_object_agg((c).name, (c).value)
                            from unnest(old_columns) c
                            where
                                (c).is_selectable
                                and ( not error_record_exceeds_max_size or (octet_length((c).value::text) <= 64))
                        )
                    )
            when action = 'DELETE' then
                jsonb_build_object(
                    'old_record',
                    (
                        select jsonb_object_agg((c).name, (c).value)
                        from unnest(old_columns) c
                        where
                            (c).is_selectable
                            and ( not error_record_exceeds_max_size or (octet_length((c).value::text) <= 64))
                            and ( not is_rls_enabled or (c).is_pkey ) -- if RLS enabled, we can't secure deletes so filter to pkey
                    )
                )
            else '{}'::jsonb
        end;

        -- Create the prepared statement
        if is_rls_enabled and action <> 'DELETE' then
            if (select 1 from pg_prepared_statements where name = 'walrus_rls_stmt' limit 1) > 0 then
                deallocate walrus_rls_stmt;
            end if;
            execute realtime.build_prepared_statement_sql('walrus_rls_stmt', entity_, columns);
        end if;

        visible_to_subscription_ids = '{}';

        for subscription_id, claims in (
                select
                    subs.subscription_id,
                    subs.claims
                from
                    unnest(subscriptions) subs
                where
                    subs.entity = entity_
                    and subs.claims_role = working_role
                    and (
                        realtime.is_visible_through_filters(columns, subs.filters)
                        or (
                          action = 'DELETE'
                          and realtime.is_visible_through_filters(old_columns, subs.filters)
                        )
                    )
        ) loop

            if not is_rls_enabled or action = 'DELETE' then
                visible_to_subscription_ids = visible_to_subscription_ids || subscription_id;
            else
                -- Check if RLS allows the role to see the record
                perform
                    -- Trim leading and trailing quotes from working_role because set_config
                    -- doesn't recognize the role as valid if they are included
                    set_config('role', trim(both '"' from working_role::text), true),
                    set_config('request.jwt.claims', claims::text, true);

                execute 'execute walrus_rls_stmt' into subscription_has_access;

                if subscription_has_access then
                    visible_to_subscription_ids = visible_to_subscription_ids || subscription_id;
                end if;
            end if;
        end loop;

        perform set_config('role', null, true);

        return next (
            output,
            is_rls_enabled,
            visible_to_subscription_ids,
            case
                when error_record_exceeds_max_size then array['Error 413: Payload Too Large']
                else '{}'
            end
        )::realtime.wal_rls;

    end if;
end loop;

perform set_config('role', null, true);
end;
$$;


--
-- Name: broadcast_changes("text", "text", "text", "text", "text", "record", "record", "text"); Type: FUNCTION; Schema: realtime; Owner: -
--

CREATE FUNCTION "realtime"."broadcast_changes"("topic_name" "text", "event_name" "text", "operation" "text", "table_name" "text", "table_schema" "text", "new" "record", "old" "record", "level" "text" DEFAULT 'ROW'::"text") RETURNS "void"
    LANGUAGE "plpgsql"
    AS $$
DECLARE
    -- Declare a variable to hold the JSONB representation of the row
    row_data jsonb := '{}'::jsonb;
BEGIN
    IF level = 'STATEMENT' THEN
        RAISE EXCEPTION 'function can only be triggered for each row, not for each statement';
    END IF;
    -- Check the operation type and handle accordingly
    IF operation = 'INSERT' OR operation = 'UPDATE' OR operation = 'DELETE' THEN
        row_data := jsonb_build_object('old_record', OLD, 'record', NEW, 'operation', operation, 'table', table_name, 'schema', table_schema);
        PERFORM realtime.send (row_data, event_name, topic_name);
    ELSE
        RAISE EXCEPTION 'Unexpected operation type: %', operation;
    END IF;
EXCEPTION
    WHEN OTHERS THEN
        RAISE EXCEPTION 'Failed to process the row: %', SQLERRM;
END;

$$;


--
-- Name: build_prepared_statement_sql("text", "regclass", "realtime"."wal_column"[]); Type: FUNCTION; Schema: realtime; Owner: -
--

CREATE FUNCTION "realtime"."build_prepared_statement_sql"("prepared_statement_name" "text", "entity" "regclass", "columns" "realtime"."wal_column"[]) RETURNS "text"
    LANGUAGE "sql"
    AS $$
      /*
      Builds a sql string that, if executed, creates a prepared statement to
      tests retrive a row from *entity* by its primary key columns.
      Example
          select realtime.build_prepared_statement_sql('public.notes', '{"id"}'::text[], '{"bigint"}'::text[])
      */
          select
      'prepare ' || prepared_statement_name || ' as
          select
              exists(
                  select
                      1
                  from
                      ' || entity || '
                  where
                      ' || string_agg(quote_ident(pkc.name) || '=' || quote_nullable(pkc.value #>> '{}') , ' and ') || '
              )'
          from
              unnest(columns) pkc
          where
              pkc.is_pkey
          group by
              entity
      $$;


--
-- Name: cast("text", "regtype"); Type: FUNCTION; Schema: realtime; Owner: -
--

CREATE FUNCTION "realtime"."cast"("val" "text", "type_" "regtype") RETURNS "jsonb"
    LANGUAGE "plpgsql" IMMUTABLE
    AS $$
    declare
      res jsonb;
    begin
      execute format('select to_jsonb(%L::'|| type_::text || ')', val)  into res;
      return res;
    end
    $$;


--
-- Name: check_equality_op("realtime"."equality_op", "regtype", "text", "text"); Type: FUNCTION; Schema: realtime; Owner: -
--

CREATE FUNCTION "realtime"."check_equality_op"("op" "realtime"."equality_op", "type_" "regtype", "val_1" "text", "val_2" "text") RETURNS boolean
    LANGUAGE "plpgsql" IMMUTABLE
    AS $$
      /*
      Casts *val_1* and *val_2* as type *type_* and check the *op* condition for truthiness
      */
      declare
          op_symbol text = (
              case
                  when op = 'eq' then '='
                  when op = 'neq' then '!='
                  when op = 'lt' then '<'
                  when op = 'lte' then '<='
                  when op = 'gt' then '>'
                  when op = 'gte' then '>='
                  when op = 'in' then '= any'
                  else 'UNKNOWN OP'
              end
          );
          res boolean;
      begin
          execute format(
              'select %L::'|| type_::text || ' ' || op_symbol
              || ' ( %L::'
              || (
                  case
                      when op = 'in' then type_::text || '[]'
                      else type_::text end
              )
              || ')', val_1, val_2) into res;
          return res;
      end;
      $$;


--
-- Name: is_visible_through_filters("realtime"."wal_column"[], "realtime"."user_defined_filter"[]); Type: FUNCTION; Schema: realtime; Owner: -
--

CREATE FUNCTION "realtime"."is_visible_through_filters"("columns" "realtime"."wal_column"[], "filters" "realtime"."user_defined_filter"[]) RETURNS boolean
    LANGUAGE "sql" IMMUTABLE
    AS $_$
    /*
    Should the record be visible (true) or filtered out (false) after *filters* are applied
    */
        select
            -- Default to allowed when no filters present
            $2 is null -- no filters. this should not happen because subscriptions has a default
            or array_length($2, 1) is null -- array length of an empty array is null
            or bool_and(
                coalesce(
                    realtime.check_equality_op(
                        op:=f.op,
                        type_:=coalesce(
                            col.type_oid::regtype, -- null when wal2json version <= 2.4
                            col.type_name::regtype
                        ),
                        -- cast jsonb to text
                        val_1:=col.value #>> '{}',
                        val_2:=f.value
                    ),
                    false -- if null, filter does not match
                )
            )
        from
            unnest(filters) f
            join unnest(columns) col
                on f.column_name = col.name;
    $_$;


--
-- Name: list_changes("name", "name", integer, integer); Type: FUNCTION; Schema: realtime; Owner: -
--

CREATE FUNCTION "realtime"."list_changes"("publication" "name", "slot_name" "name", "max_changes" integer, "max_record_bytes" integer) RETURNS SETOF "realtime"."wal_rls"
    LANGUAGE "sql"
    SET "log_min_messages" TO 'fatal'
    AS $$
      with pub as (
        select
          concat_ws(
            ',',
            case when bool_or(pubinsert) then 'insert' else null end,
            case when bool_or(pubupdate) then 'update' else null end,
            case when bool_or(pubdelete) then 'delete' else null end
          ) as w2j_actions,
          coalesce(
            string_agg(
              realtime.quote_wal2json(format('%I.%I', schemaname, tablename)::regclass),
              ','
            ) filter (where ppt.tablename is not null and ppt.tablename not like '% %'),
            ''
          ) w2j_add_tables
        from
          pg_publication pp
          left join pg_publication_tables ppt
            on pp.pubname = ppt.pubname
        where
          pp.pubname = publication
        group by
          pp.pubname
        limit 1
      ),
      w2j as (
        select
          x.*, pub.w2j_add_tables
        from
          pub,
          pg_logical_slot_get_changes(
            slot_name, null, max_changes,
            'include-pk', 'true',
            'include-transaction', 'false',
            'include-timestamp', 'true',
            'include-type-oids', 'true',
            'format-version', '2',
            'actions', pub.w2j_actions,
            'add-tables', pub.w2j_add_tables
          ) x
      )
      select
        xyz.wal,
        xyz.is_rls_enabled,
        xyz.subscription_ids,
        xyz.errors
      from
        w2j,
        realtime.apply_rls(
          wal := w2j.data::jsonb,
          max_record_bytes := max_record_bytes
        ) xyz(wal, is_rls_enabled, subscription_ids, errors)
      where
        w2j.w2j_add_tables <> ''
        and xyz.subscription_ids[1] is not null
    $$;


--
-- Name: quote_wal2json("regclass"); Type: FUNCTION; Schema: realtime; Owner: -
--

CREATE FUNCTION "realtime"."quote_wal2json"("entity" "regclass") RETURNS "text"
    LANGUAGE "sql" IMMUTABLE STRICT
    AS $$
      select
        (
          select string_agg('' || ch,'')
          from unnest(string_to_array(nsp.nspname::text, null)) with ordinality x(ch, idx)
          where
            not (x.idx = 1 and x.ch = '"')
            and not (
              x.idx = array_length(string_to_array(nsp.nspname::text, null), 1)
              and x.ch = '"'
            )
        )
        || '.'
        || (
          select string_agg('' || ch,'')
          from unnest(string_to_array(pc.relname::text, null)) with ordinality x(ch, idx)
          where
            not (x.idx = 1 and x.ch = '"')
            and not (
              x.idx = array_length(string_to_array(nsp.nspname::text, null), 1)
              and x.ch = '"'
            )
          )
      from
        pg_class pc
        join pg_namespace nsp
          on pc.relnamespace = nsp.oid
      where
        pc.oid = entity
    $$;


--
-- Name: send("jsonb", "text", "text", boolean); Type: FUNCTION; Schema: realtime; Owner: -
--

CREATE FUNCTION "realtime"."send"("payload" "jsonb", "event" "text", "topic" "text", "private" boolean DEFAULT true) RETURNS "void"
    LANGUAGE "plpgsql"
    AS $$
BEGIN
  BEGIN
    -- Set the topic configuration
    EXECUTE format('SET LOCAL realtime.topic TO %L', topic);

    -- Attempt to insert the message
    INSERT INTO realtime.messages (payload, event, topic, private, extension)
    VALUES (payload, event, topic, private, 'broadcast');
  EXCEPTION
    WHEN OTHERS THEN
      -- Capture and notify the error
      RAISE WARNING 'ErrorSendingBroadcastMessage: %', SQLERRM;
  END;
END;
$$;


--
-- Name: subscription_check_filters(); Type: FUNCTION; Schema: realtime; Owner: -
--

CREATE FUNCTION "realtime"."subscription_check_filters"() RETURNS "trigger"
    LANGUAGE "plpgsql"
    AS $$
    /*
    Validates that the user defined filters for a subscription:
    - refer to valid columns that the claimed role may access
    - values are coercable to the correct column type
    */
    declare
        col_names text[] = coalesce(
                array_agg(c.column_name order by c.ordinal_position),
                '{}'::text[]
            )
            from
                information_schema.columns c
            where
                format('%I.%I', c.table_schema, c.table_name)::regclass = new.entity
                and pg_catalog.has_column_privilege(
                    (new.claims ->> 'role'),
                    format('%I.%I', c.table_schema, c.table_name)::regclass,
                    c.column_name,
                    'SELECT'
                );
        filter realtime.user_defined_filter;
        col_type regtype;

        in_val jsonb;
    begin
        for filter in select * from unnest(new.filters) loop
            -- Filtered column is valid
            if not filter.column_name = any(col_names) then
                raise exception 'invalid column for filter %', filter.column_name;
            end if;

            -- Type is sanitized and safe for string interpolation
            col_type = (
                select atttypid::regtype
                from pg_catalog.pg_attribute
                where attrelid = new.entity
                      and attname = filter.column_name
            );
            if col_type is null then
                raise exception 'failed to lookup type for column %', filter.column_name;
            end if;

            -- Set maximum number of entries for in filter
            if filter.op = 'in'::realtime.equality_op then
                in_val = realtime.cast(filter.value, (col_type::text || '[]')::regtype);
                if coalesce(jsonb_array_length(in_val), 0) > 100 then
                    raise exception 'too many values for `in` filter. Maximum 100';
                end if;
            else
                -- raises an exception if value is not coercable to type
                perform realtime.cast(filter.value, col_type);
            end if;

        end loop;

        -- Apply consistent order to filters so the unique constraint on
        -- (subscription_id, entity, filters) can't be tricked by a different filter order
        new.filters = coalesce(
            array_agg(f order by f.column_name, f.op, f.value),
            '{}'
        ) from unnest(new.filters) f;

        return new;
    end;
    $$;


--
-- Name: to_regrole("text"); Type: FUNCTION; Schema: realtime; Owner: -
--

CREATE FUNCTION "realtime"."to_regrole"("role_name" "text") RETURNS "regrole"
    LANGUAGE "sql" IMMUTABLE
    AS $$ select role_name::regrole $$;


--
-- Name: topic(); Type: FUNCTION; Schema: realtime; Owner: -
--

CREATE FUNCTION "realtime"."topic"() RETURNS "text"
    LANGUAGE "sql" STABLE
    AS $$
select nullif(current_setting('realtime.topic', true), '')::text;
$$;


--
-- Name: add_prefixes("text", "text"); Type: FUNCTION; Schema: storage; Owner: -
--

CREATE FUNCTION "storage"."add_prefixes"("_bucket_id" "text", "_name" "text") RETURNS "void"
    LANGUAGE "plpgsql" SECURITY DEFINER
    AS $$
DECLARE
    prefixes text[];
BEGIN
    prefixes := "storage"."get_prefixes"("_name");

    IF array_length(prefixes, 1) > 0 THEN
        INSERT INTO storage.prefixes (name, bucket_id)
        SELECT UNNEST(prefixes) as name, "_bucket_id" ON CONFLICT DO NOTHING;
    END IF;
END;
$$;


--
-- Name: can_insert_object("text", "text", "uuid", "jsonb"); Type: FUNCTION; Schema: storage; Owner: -
--

CREATE FUNCTION "storage"."can_insert_object"("bucketid" "text", "name" "text", "owner" "uuid", "metadata" "jsonb") RETURNS "void"
    LANGUAGE "plpgsql"
    AS $$
BEGIN
  INSERT INTO "storage"."objects" ("bucket_id", "name", "owner", "metadata") VALUES (bucketid, name, owner, metadata);
  -- hack to rollback the successful insert
  RAISE sqlstate 'PT200' using
  message = 'ROLLBACK',
  detail = 'rollback successful insert';
END
$$;


--
-- Name: delete_leaf_prefixes("text"[], "text"[]); Type: FUNCTION; Schema: storage; Owner: -
--

CREATE FUNCTION "storage"."delete_leaf_prefixes"("bucket_ids" "text"[], "names" "text"[]) RETURNS "void"
    LANGUAGE "plpgsql" SECURITY DEFINER
    AS $$
DECLARE
    v_rows_deleted integer;
BEGIN
    LOOP
        WITH candidates AS (
            SELECT DISTINCT
                t.bucket_id,
                unnest(storage.get_prefixes(t.name)) AS name
            FROM unnest(bucket_ids, names) AS t(bucket_id, name)
        ),
        uniq AS (
             SELECT
                 bucket_id,
                 name,
                 storage.get_level(name) AS level
             FROM candidates
             WHERE name <> ''
             GROUP BY bucket_id, name
        ),
        leaf AS (
             SELECT
                 p.bucket_id,
                 p.name,
                 p.level
             FROM storage.prefixes AS p
                  JOIN uniq AS u
                       ON u.bucket_id = p.bucket_id
                           AND u.name = p.name
                           AND u.level = p.level
             WHERE NOT EXISTS (
                 SELECT 1
                 FROM storage.objects AS o
                 WHERE o.bucket_id = p.bucket_id
                   AND o.level = p.level + 1
                   AND o.name COLLATE "C" LIKE p.name || '/%'
             )
             AND NOT EXISTS (
                 SELECT 1
                 FROM storage.prefixes AS c
                 WHERE c.bucket_id = p.bucket_id
                   AND c.level = p.level + 1
                   AND c.name COLLATE "C" LIKE p.name || '/%'
             )
        )
        DELETE
        FROM storage.prefixes AS p
            USING leaf AS l
        WHERE p.bucket_id = l.bucket_id
          AND p.name = l.name
          AND p.level = l.level;

        GET DIAGNOSTICS v_rows_deleted = ROW_COUNT;
        EXIT WHEN v_rows_deleted = 0;
    END LOOP;
END;
$$;


--
-- Name: delete_prefix("text", "text"); Type: FUNCTION; Schema: storage; Owner: -
--

CREATE FUNCTION "storage"."delete_prefix"("_bucket_id" "text", "_name" "text") RETURNS boolean
    LANGUAGE "plpgsql" SECURITY DEFINER
    AS $$
BEGIN
    -- Check if we can delete the prefix
    IF EXISTS(
        SELECT FROM "storage"."prefixes"
        WHERE "prefixes"."bucket_id" = "_bucket_id"
          AND level = "storage"."get_level"("_name") + 1
          AND "prefixes"."name" COLLATE "C" LIKE "_name" || '/%'
        LIMIT 1
    )
    OR EXISTS(
        SELECT FROM "storage"."objects"
        WHERE "objects"."bucket_id" = "_bucket_id"
          AND "storage"."get_level"("objects"."name") = "storage"."get_level"("_name") + 1
          AND "objects"."name" COLLATE "C" LIKE "_name" || '/%'
        LIMIT 1
    ) THEN
    -- There are sub-objects, skip deletion
    RETURN false;
    ELSE
        DELETE FROM "storage"."prefixes"
        WHERE "prefixes"."bucket_id" = "_bucket_id"
          AND level = "storage"."get_level"("_name")
          AND "prefixes"."name" = "_name";
        RETURN true;
    END IF;
END;
$$;


--
-- Name: delete_prefix_hierarchy_trigger(); Type: FUNCTION; Schema: storage; Owner: -
--

CREATE FUNCTION "storage"."delete_prefix_hierarchy_trigger"() RETURNS "trigger"
    LANGUAGE "plpgsql"
    AS $$
DECLARE
    prefix text;
BEGIN
    prefix := "storage"."get_prefix"(OLD."name");

    IF coalesce(prefix, '') != '' THEN
        PERFORM "storage"."delete_prefix"(OLD."bucket_id", prefix);
    END IF;

    RETURN OLD;
END;
$$;


--
-- Name: enforce_bucket_name_length(); Type: FUNCTION; Schema: storage; Owner: -
--

CREATE FUNCTION "storage"."enforce_bucket_name_length"() RETURNS "trigger"
    LANGUAGE "plpgsql"
    AS $$
begin
    if length(new.name) > 100 then
        raise exception 'bucket name "%" is too long (% characters). Max is 100.', new.name, length(new.name);
    end if;
    return new;
end;
$$;


--
-- Name: extension("text"); Type: FUNCTION; Schema: storage; Owner: -
--

CREATE FUNCTION "storage"."extension"("name" "text") RETURNS "text"
    LANGUAGE "plpgsql" IMMUTABLE
    AS $$
DECLARE
    _parts text[];
    _filename text;
BEGIN
    SELECT string_to_array(name, '/') INTO _parts;
    SELECT _parts[array_length(_parts,1)] INTO _filename;
    RETURN reverse(split_part(reverse(_filename), '.', 1));
END
$$;


--
-- Name: filename("text"); Type: FUNCTION; Schema: storage; Owner: -
--

CREATE FUNCTION "storage"."filename"("name" "text") RETURNS "text"
    LANGUAGE "plpgsql"
    AS $$
DECLARE
_parts text[];
BEGIN
	select string_to_array(name, '/') into _parts;
	return _parts[array_length(_parts,1)];
END
$$;


--
-- Name: foldername("text"); Type: FUNCTION; Schema: storage; Owner: -
--

CREATE FUNCTION "storage"."foldername"("name" "text") RETURNS "text"[]
    LANGUAGE "plpgsql" IMMUTABLE
    AS $$
DECLARE
    _parts text[];
BEGIN
    -- Split on "/" to get path segments
    SELECT string_to_array(name, '/') INTO _parts;
    -- Return everything except the last segment
    RETURN _parts[1 : array_length(_parts,1) - 1];
END
$$;


--
-- Name: get_level("text"); Type: FUNCTION; Schema: storage; Owner: -
--

CREATE FUNCTION "storage"."get_level"("name" "text") RETURNS integer
    LANGUAGE "sql" IMMUTABLE STRICT
    AS $$
SELECT array_length(string_to_array("name", '/'), 1);
$$;


--
-- Name: get_prefix("text"); Type: FUNCTION; Schema: storage; Owner: -
--

CREATE FUNCTION "storage"."get_prefix"("name" "text") RETURNS "text"
    LANGUAGE "sql" IMMUTABLE STRICT
    AS $_$
SELECT
    CASE WHEN strpos("name", '/') > 0 THEN
             regexp_replace("name", '[\/]{1}[^\/]+\/?$', '')
         ELSE
             ''
        END;
$_$;


--
-- Name: get_prefixes("text"); Type: FUNCTION; Schema: storage; Owner: -
--

CREATE FUNCTION "storage"."get_prefixes"("name" "text") RETURNS "text"[]
    LANGUAGE "plpgsql" IMMUTABLE STRICT
    AS $$
DECLARE
    parts text[];
    prefixes text[];
    prefix text;
BEGIN
    -- Split the name into parts by '/'
    parts := string_to_array("name", '/');
    prefixes := '{}';

    -- Construct the prefixes, stopping one level below the last part
    FOR i IN 1..array_length(parts, 1) - 1 LOOP
            prefix := array_to_string(parts[1:i], '/');
            prefixes := array_append(prefixes, prefix);
    END LOOP;

    RETURN prefixes;
END;
$$;


--
-- Name: get_size_by_bucket(); Type: FUNCTION; Schema: storage; Owner: -
--

CREATE FUNCTION "storage"."get_size_by_bucket"() RETURNS TABLE("size" bigint, "bucket_id" "text")
    LANGUAGE "plpgsql" STABLE
    AS $$
BEGIN
    return query
        select sum((metadata->>'size')::bigint) as size, obj.bucket_id
        from "storage".objects as obj
        group by obj.bucket_id;
END
$$;


--
-- Name: list_multipart_uploads_with_delimiter("text", "text", "text", integer, "text", "text"); Type: FUNCTION; Schema: storage; Owner: -
--

CREATE FUNCTION "storage"."list_multipart_uploads_with_delimiter"("bucket_id" "text", "prefix_param" "text", "delimiter_param" "text", "max_keys" integer DEFAULT 100, "next_key_token" "text" DEFAULT ''::"text", "next_upload_token" "text" DEFAULT ''::"text") RETURNS TABLE("key" "text", "id" "text", "created_at" timestamp with time zone)
    LANGUAGE "plpgsql"
    AS $_$
BEGIN
    RETURN QUERY EXECUTE
        'SELECT DISTINCT ON(key COLLATE "C") * from (
            SELECT
                CASE
                    WHEN position($2 IN substring(key from length($1) + 1)) > 0 THEN
                        substring(key from 1 for length($1) + position($2 IN substring(key from length($1) + 1)))
                    ELSE
                        key
                END AS key, id, created_at
            FROM
                storage.s3_multipart_uploads
            WHERE
                bucket_id = $5 AND
                key ILIKE $1 || ''%'' AND
                CASE
                    WHEN $4 != '''' AND $6 = '''' THEN
                        CASE
                            WHEN position($2 IN substring(key from length($1) + 1)) > 0 THEN
                                substring(key from 1 for length($1) + position($2 IN substring(key from length($1) + 1))) COLLATE "C" > $4
                            ELSE
                                key COLLATE "C" > $4
                            END
                    ELSE
                        true
                END AND
                CASE
                    WHEN $6 != '''' THEN
                        id COLLATE "C" > $6
                    ELSE
                        true
                    END
            ORDER BY
                key COLLATE "C" ASC, created_at ASC) as e order by key COLLATE "C" LIMIT $3'
        USING prefix_param, delimiter_param, max_keys, next_key_token, bucket_id, next_upload_token;
END;
$_$;


--
-- Name: list_objects_with_delimiter("text", "text", "text", integer, "text", "text"); Type: FUNCTION; Schema: storage; Owner: -
--

CREATE FUNCTION "storage"."list_objects_with_delimiter"("bucket_id" "text", "prefix_param" "text", "delimiter_param" "text", "max_keys" integer DEFAULT 100, "start_after" "text" DEFAULT ''::"text", "next_token" "text" DEFAULT ''::"text") RETURNS TABLE("name" "text", "id" "uuid", "metadata" "jsonb", "updated_at" timestamp with time zone)
    LANGUAGE "plpgsql"
    AS $_$
BEGIN
    RETURN QUERY EXECUTE
        'SELECT DISTINCT ON(name COLLATE "C") * from (
            SELECT
                CASE
                    WHEN position($2 IN substring(name from length($1) + 1)) > 0 THEN
                        substring(name from 1 for length($1) + position($2 IN substring(name from length($1) + 1)))
                    ELSE
                        name
                END AS name, id, metadata, updated_at
            FROM
                storage.objects
            WHERE
                bucket_id = $5 AND
                name ILIKE $1 || ''%'' AND
                CASE
                    WHEN $6 != '''' THEN
                    name COLLATE "C" > $6
                ELSE true END
                AND CASE
                    WHEN $4 != '''' THEN
                        CASE
                            WHEN position($2 IN substring(name from length($1) + 1)) > 0 THEN
                                substring(name from 1 for length($1) + position($2 IN substring(name from length($1) + 1))) COLLATE "C" > $4
                            ELSE
                                name COLLATE "C" > $4
                            END
                    ELSE
                        true
                END
            ORDER BY
                name COLLATE "C" ASC) as e order by name COLLATE "C" LIMIT $3'
        USING prefix_param, delimiter_param, max_keys, next_token, bucket_id, start_after;
END;
$_$;


--
-- Name: lock_top_prefixes("text"[], "text"[]); Type: FUNCTION; Schema: storage; Owner: -
--

CREATE FUNCTION "storage"."lock_top_prefixes"("bucket_ids" "text"[], "names" "text"[]) RETURNS "void"
    LANGUAGE "plpgsql" SECURITY DEFINER
    AS $$
DECLARE
    v_bucket text;
    v_top text;
BEGIN
    FOR v_bucket, v_top IN
        SELECT DISTINCT t.bucket_id,
            split_part(t.name, '/', 1) AS top
        FROM unnest(bucket_ids, names) AS t(bucket_id, name)
        WHERE t.name <> ''
        ORDER BY 1, 2
        LOOP
            PERFORM pg_advisory_xact_lock(hashtextextended(v_bucket || '/' || v_top, 0));
        END LOOP;
END;
$$;


--
-- Name: objects_delete_cleanup(); Type: FUNCTION; Schema: storage; Owner: -
--

CREATE FUNCTION "storage"."objects_delete_cleanup"() RETURNS "trigger"
    LANGUAGE "plpgsql" SECURITY DEFINER
    AS $$
DECLARE
    v_bucket_ids text[];
    v_names      text[];
BEGIN
    IF current_setting('storage.gc.prefixes', true) = '1' THEN
        RETURN NULL;
    END IF;

    PERFORM set_config('storage.gc.prefixes', '1', true);

    SELECT COALESCE(array_agg(d.bucket_id), '{}'),
           COALESCE(array_agg(d.name), '{}')
    INTO v_bucket_ids, v_names
    FROM deleted AS d
    WHERE d.name <> '';

    PERFORM storage.lock_top_prefixes(v_bucket_ids, v_names);
    PERFORM storage.delete_leaf_prefixes(v_bucket_ids, v_names);

    RETURN NULL;
END;
$$;


--
-- Name: objects_insert_prefix_trigger(); Type: FUNCTION; Schema: storage; Owner: -
--

CREATE FUNCTION "storage"."objects_insert_prefix_trigger"() RETURNS "trigger"
    LANGUAGE "plpgsql"
    AS $$
BEGIN
    PERFORM "storage"."add_prefixes"(NEW."bucket_id", NEW."name");
    NEW.level := "storage"."get_level"(NEW."name");

    RETURN NEW;
END;
$$;


--
-- Name: objects_update_cleanup(); Type: FUNCTION; Schema: storage; Owner: -
--

CREATE FUNCTION "storage"."objects_update_cleanup"() RETURNS "trigger"
    LANGUAGE "plpgsql" SECURITY DEFINER
    AS $$
DECLARE
    -- NEW - OLD (destinations to create prefixes for)
    v_add_bucket_ids text[];
    v_add_names      text[];

    -- OLD - NEW (sources to prune)
    v_src_bucket_ids text[];
    v_src_names      text[];
BEGIN
    IF TG_OP <> 'UPDATE' THEN
        RETURN NULL;
    END IF;

    -- 1) Compute NEW−OLD (added paths) and OLD−NEW (moved-away paths)
    WITH added AS (
        SELECT n.bucket_id, n.name
        FROM new_rows n
        WHERE n.name <> '' AND position('/' in n.name) > 0
        EXCEPT
        SELECT o.bucket_id, o.name FROM old_rows o WHERE o.name <> ''
    ),
    moved AS (
         SELECT o.bucket_id, o.name
         FROM old_rows o
         WHERE o.name <> ''
         EXCEPT
         SELECT n.bucket_id, n.name FROM new_rows n WHERE n.name <> ''
    )
    SELECT
        -- arrays for ADDED (dest) in stable order
        COALESCE( (SELECT array_agg(a.bucket_id ORDER BY a.bucket_id, a.name) FROM added a), '{}' ),
        COALESCE( (SELECT array_agg(a.name      ORDER BY a.bucket_id, a.name) FROM added a), '{}' ),
        -- arrays for MOVED (src) in stable order
        COALESCE( (SELECT array_agg(m.bucket_id ORDER BY m.bucket_id, m.name) FROM moved m), '{}' ),
        COALESCE( (SELECT array_agg(m.name      ORDER BY m.bucket_id, m.name) FROM moved m), '{}' )
    INTO v_add_bucket_ids, v_add_names, v_src_bucket_ids, v_src_names;

    -- Nothing to do?
    IF (array_length(v_add_bucket_ids, 1) IS NULL) AND (array_length(v_src_bucket_ids, 1) IS NULL) THEN
        RETURN NULL;
    END IF;

    -- 2) Take per-(bucket, top) locks: ALL prefixes in consistent global order to prevent deadlocks
    DECLARE
        v_all_bucket_ids text[];
        v_all_names text[];
    BEGIN
        -- Combine source and destination arrays for consistent lock ordering
        v_all_bucket_ids := COALESCE(v_src_bucket_ids, '{}') || COALESCE(v_add_bucket_ids, '{}');
        v_all_names := COALESCE(v_src_names, '{}') || COALESCE(v_add_names, '{}');

        -- Single lock call ensures consistent global ordering across all transactions
        IF array_length(v_all_bucket_ids, 1) IS NOT NULL THEN
            PERFORM storage.lock_top_prefixes(v_all_bucket_ids, v_all_names);
        END IF;
    END;

    -- 3) Create destination prefixes (NEW−OLD) BEFORE pruning sources
    IF array_length(v_add_bucket_ids, 1) IS NOT NULL THEN
        WITH candidates AS (
            SELECT DISTINCT t.bucket_id, unnest(storage.get_prefixes(t.name)) AS name
            FROM unnest(v_add_bucket_ids, v_add_names) AS t(bucket_id, name)
            WHERE name <> ''
        )
        INSERT INTO storage.prefixes (bucket_id, name)
        SELECT c.bucket_id, c.name
        FROM candidates c
        ON CONFLICT DO NOTHING;
    END IF;

    -- 4) Prune source prefixes bottom-up for OLD−NEW
    IF array_length(v_src_bucket_ids, 1) IS NOT NULL THEN
        -- re-entrancy guard so DELETE on prefixes won't recurse
        IF current_setting('storage.gc.prefixes', true) <> '1' THEN
            PERFORM set_config('storage.gc.prefixes', '1', true);
        END IF;

        PERFORM storage.delete_leaf_prefixes(v_src_bucket_ids, v_src_names);
    END IF;

    RETURN NULL;
END;
$$;


--
-- Name: objects_update_level_trigger(); Type: FUNCTION; Schema: storage; Owner: -
--

CREATE FUNCTION "storage"."objects_update_level_trigger"() RETURNS "trigger"
    LANGUAGE "plpgsql"
    AS $$
BEGIN
    -- Ensure this is an update operation and the name has changed
    IF TG_OP = 'UPDATE' AND (NEW."name" <> OLD."name" OR NEW."bucket_id" <> OLD."bucket_id") THEN
        -- Set the new level
        NEW."level" := "storage"."get_level"(NEW."name");
    END IF;
    RETURN NEW;
END;
$$;


--
-- Name: objects_update_prefix_trigger(); Type: FUNCTION; Schema: storage; Owner: -
--

CREATE FUNCTION "storage"."objects_update_prefix_trigger"() RETURNS "trigger"
    LANGUAGE "plpgsql"
    AS $$
DECLARE
    old_prefixes TEXT[];
BEGIN
    -- Ensure this is an update operation and the name has changed
    IF TG_OP = 'UPDATE' AND (NEW."name" <> OLD."name" OR NEW."bucket_id" <> OLD."bucket_id") THEN
        -- Retrieve old prefixes
        old_prefixes := "storage"."get_prefixes"(OLD."name");

        -- Remove old prefixes that are only used by this object
        WITH all_prefixes as (
            SELECT unnest(old_prefixes) as prefix
        ),
        can_delete_prefixes as (
             SELECT prefix
             FROM all_prefixes
             WHERE NOT EXISTS (
                 SELECT 1 FROM "storage"."objects"
                 WHERE "bucket_id" = OLD."bucket_id"
                   AND "name" <> OLD."name"
                   AND "name" LIKE (prefix || '%')
             )
         )
        DELETE FROM "storage"."prefixes" WHERE name IN (SELECT prefix FROM can_delete_prefixes);

        -- Add new prefixes
        PERFORM "storage"."add_prefixes"(NEW."bucket_id", NEW."name");
    END IF;
    -- Set the new level
    NEW."level" := "storage"."get_level"(NEW."name");

    RETURN NEW;
END;
$$;


--
-- Name: operation(); Type: FUNCTION; Schema: storage; Owner: -
--

CREATE FUNCTION "storage"."operation"() RETURNS "text"
    LANGUAGE "plpgsql" STABLE
    AS $$
BEGIN
    RETURN current_setting('storage.operation', true);
END;
$$;


--
-- Name: prefixes_delete_cleanup(); Type: FUNCTION; Schema: storage; Owner: -
--

CREATE FUNCTION "storage"."prefixes_delete_cleanup"() RETURNS "trigger"
    LANGUAGE "plpgsql" SECURITY DEFINER
    AS $$
DECLARE
    v_bucket_ids text[];
    v_names      text[];
BEGIN
    IF current_setting('storage.gc.prefixes', true) = '1' THEN
        RETURN NULL;
    END IF;

    PERFORM set_config('storage.gc.prefixes', '1', true);

    SELECT COALESCE(array_agg(d.bucket_id), '{}'),
           COALESCE(array_agg(d.name), '{}')
    INTO v_bucket_ids, v_names
    FROM deleted AS d
    WHERE d.name <> '';

    PERFORM storage.lock_top_prefixes(v_bucket_ids, v_names);
    PERFORM storage.delete_leaf_prefixes(v_bucket_ids, v_names);

    RETURN NULL;
END;
$$;


--
-- Name: prefixes_insert_trigger(); Type: FUNCTION; Schema: storage; Owner: -
--

CREATE FUNCTION "storage"."prefixes_insert_trigger"() RETURNS "trigger"
    LANGUAGE "plpgsql"
    AS $$
BEGIN
    PERFORM "storage"."add_prefixes"(NEW."bucket_id", NEW."name");
    RETURN NEW;
END;
$$;


--
-- Name: search("text", "text", integer, integer, integer, "text", "text", "text"); Type: FUNCTION; Schema: storage; Owner: -
--

CREATE FUNCTION "storage"."search"("prefix" "text", "bucketname" "text", "limits" integer DEFAULT 100, "levels" integer DEFAULT 1, "offsets" integer DEFAULT 0, "search" "text" DEFAULT ''::"text", "sortcolumn" "text" DEFAULT 'name'::"text", "sortorder" "text" DEFAULT 'asc'::"text") RETURNS TABLE("name" "text", "id" "uuid", "updated_at" timestamp with time zone, "created_at" timestamp with time zone, "last_accessed_at" timestamp with time zone, "metadata" "jsonb")
    LANGUAGE "plpgsql"
    AS $$
declare
    can_bypass_rls BOOLEAN;
begin
    SELECT rolbypassrls
    INTO can_bypass_rls
    FROM pg_roles
    WHERE rolname = coalesce(nullif(current_setting('role', true), 'none'), current_user);

    IF can_bypass_rls THEN
        RETURN QUERY SELECT * FROM storage.search_v1_optimised(prefix, bucketname, limits, levels, offsets, search, sortcolumn, sortorder);
    ELSE
        RETURN QUERY SELECT * FROM storage.search_legacy_v1(prefix, bucketname, limits, levels, offsets, search, sortcolumn, sortorder);
    END IF;
end;
$$;


--
-- Name: search_legacy_v1("text", "text", integer, integer, integer, "text", "text", "text"); Type: FUNCTION; Schema: storage; Owner: -
--

CREATE FUNCTION "storage"."search_legacy_v1"("prefix" "text", "bucketname" "text", "limits" integer DEFAULT 100, "levels" integer DEFAULT 1, "offsets" integer DEFAULT 0, "search" "text" DEFAULT ''::"text", "sortcolumn" "text" DEFAULT 'name'::"text", "sortorder" "text" DEFAULT 'asc'::"text") RETURNS TABLE("name" "text", "id" "uuid", "updated_at" timestamp with time zone, "created_at" timestamp with time zone, "last_accessed_at" timestamp with time zone, "metadata" "jsonb")
    LANGUAGE "plpgsql" STABLE
    AS $_$
declare
    v_order_by text;
    v_sort_order text;
begin
    case
        when sortcolumn = 'name' then
            v_order_by = 'name';
        when sortcolumn = 'updated_at' then
            v_order_by = 'updated_at';
        when sortcolumn = 'created_at' then
            v_order_by = 'created_at';
        when sortcolumn = 'last_accessed_at' then
            v_order_by = 'last_accessed_at';
        else
            v_order_by = 'name';
        end case;

    case
        when sortorder = 'asc' then
            v_sort_order = 'asc';
        when sortorder = 'desc' then
            v_sort_order = 'desc';
        else
            v_sort_order = 'asc';
        end case;

    v_order_by = v_order_by || ' ' || v_sort_order;

    return query execute
        'with folders as (
           select path_tokens[$1] as folder
           from storage.objects
             where objects.name ilike $2 || $3 || ''%''
               and bucket_id = $4
               and array_length(objects.path_tokens, 1) <> $1
           group by folder
           order by folder ' || v_sort_order || '
     )
     (select folder as "name",
            null as id,
            null as updated_at,
            null as created_at,
            null as last_accessed_at,
            null as metadata from folders)
     union all
     (select path_tokens[$1] as "name",
            id,
            updated_at,
            created_at,
            last_accessed_at,
            metadata
     from storage.objects
     where objects.name ilike $2 || $3 || ''%''
       and bucket_id = $4
       and array_length(objects.path_tokens, 1) = $1
     order by ' || v_order_by || ')
     limit $5
     offset $6' using levels, prefix, search, bucketname, limits, offsets;
end;
$_$;


--
-- Name: search_v1_optimised("text", "text", integer, integer, integer, "text", "text", "text"); Type: FUNCTION; Schema: storage; Owner: -
--

CREATE FUNCTION "storage"."search_v1_optimised"("prefix" "text", "bucketname" "text", "limits" integer DEFAULT 100, "levels" integer DEFAULT 1, "offsets" integer DEFAULT 0, "search" "text" DEFAULT ''::"text", "sortcolumn" "text" DEFAULT 'name'::"text", "sortorder" "text" DEFAULT 'asc'::"text") RETURNS TABLE("name" "text", "id" "uuid", "updated_at" timestamp with time zone, "created_at" timestamp with time zone, "last_accessed_at" timestamp with time zone, "metadata" "jsonb")
    LANGUAGE "plpgsql" STABLE
    AS $_$
declare
    v_order_by text;
    v_sort_order text;
begin
    case
        when sortcolumn = 'name' then
            v_order_by = 'name';
        when sortcolumn = 'updated_at' then
            v_order_by = 'updated_at';
        when sortcolumn = 'created_at' then
            v_order_by = 'created_at';
        when sortcolumn = 'last_accessed_at' then
            v_order_by = 'last_accessed_at';
        else
            v_order_by = 'name';
        end case;

    case
        when sortorder = 'asc' then
            v_sort_order = 'asc';
        when sortorder = 'desc' then
            v_sort_order = 'desc';
        else
            v_sort_order = 'asc';
        end case;

    v_order_by = v_order_by || ' ' || v_sort_order;

    return query execute
        'with folders as (
           select (string_to_array(name, ''/''))[level] as name
           from storage.prefixes
             where lower(prefixes.name) like lower($2 || $3) || ''%''
               and bucket_id = $4
               and level = $1
           order by name ' || v_sort_order || '
     )
     (select name,
            null as id,
            null as updated_at,
            null as created_at,
            null as last_accessed_at,
            null as metadata from folders)
     union all
     (select path_tokens[level] as "name",
            id,
            updated_at,
            created_at,
            last_accessed_at,
            metadata
     from storage.objects
     where lower(objects.name) like lower($2 || $3) || ''%''
       and bucket_id = $4
       and level = $1
     order by ' || v_order_by || ')
     limit $5
     offset $6' using levels, prefix, search, bucketname, limits, offsets;
end;
$_$;


--
-- Name: search_v2("text", "text", integer, integer, "text", "text", "text", "text"); Type: FUNCTION; Schema: storage; Owner: -
--

CREATE FUNCTION "storage"."search_v2"("prefix" "text", "bucket_name" "text", "limits" integer DEFAULT 100, "levels" integer DEFAULT 1, "start_after" "text" DEFAULT ''::"text", "sort_order" "text" DEFAULT 'asc'::"text", "sort_column" "text" DEFAULT 'name'::"text", "sort_column_after" "text" DEFAULT ''::"text") RETURNS TABLE("key" "text", "name" "text", "id" "uuid", "updated_at" timestamp with time zone, "created_at" timestamp with time zone, "last_accessed_at" timestamp with time zone, "metadata" "jsonb")
    LANGUAGE "plpgsql" STABLE
    AS $_$
DECLARE
    sort_col text;
    sort_ord text;
    cursor_op text;
    cursor_expr text;
    sort_expr text;
BEGIN
    -- Validate sort_order
    sort_ord := lower(sort_order);
    IF sort_ord NOT IN ('asc', 'desc') THEN
        sort_ord := 'asc';
    END IF;

    -- Determine cursor comparison operator
    IF sort_ord = 'asc' THEN
        cursor_op := '>';
    ELSE
        cursor_op := '<';
    END IF;
    
    sort_col := lower(sort_column);
    -- Validate sort column  
    IF sort_col IN ('updated_at', 'created_at') THEN
        cursor_expr := format(
            '($5 = '''' OR ROW(date_trunc(''milliseconds'', %I), name COLLATE "C") %s ROW(COALESCE(NULLIF($6, '''')::timestamptz, ''epoch''::timestamptz), $5))',
            sort_col, cursor_op
        );
        sort_expr := format(
            'COALESCE(date_trunc(''milliseconds'', %I), ''epoch''::timestamptz) %s, name COLLATE "C" %s',
            sort_col, sort_ord, sort_ord
        );
    ELSE
        cursor_expr := format('($5 = '''' OR name COLLATE "C" %s $5)', cursor_op);
        sort_expr := format('name COLLATE "C" %s', sort_ord);
    END IF;

    RETURN QUERY EXECUTE format(
        $sql$
        SELECT * FROM (
            (
                SELECT
                    split_part(name, '/', $4) AS key,
                    name,
                    NULL::uuid AS id,
                    updated_at,
                    created_at,
                    NULL::timestamptz AS last_accessed_at,
                    NULL::jsonb AS metadata
                FROM storage.prefixes
                WHERE name COLLATE "C" LIKE $1 || '%%'
                    AND bucket_id = $2
                    AND level = $4
                    AND %s
                ORDER BY %s
                LIMIT $3
            )
            UNION ALL
            (
                SELECT
                    split_part(name, '/', $4) AS key,
                    name,
                    id,
                    updated_at,
                    created_at,
                    last_accessed_at,
                    metadata
                FROM storage.objects
                WHERE name COLLATE "C" LIKE $1 || '%%'
                    AND bucket_id = $2
                    AND level = $4
                    AND %s
                ORDER BY %s
                LIMIT $3
            )
        ) obj
        ORDER BY %s
        LIMIT $3
        $sql$,
        cursor_expr,    -- prefixes WHERE
        sort_expr,      -- prefixes ORDER BY
        cursor_expr,    -- objects WHERE
        sort_expr,      -- objects ORDER BY
        sort_expr       -- final ORDER BY
    )
    USING prefix, bucket_name, limits, levels, start_after, sort_column_after;
END;
$_$;


--
-- Name: update_updated_at_column(); Type: FUNCTION; Schema: storage; Owner: -
--

CREATE FUNCTION "storage"."update_updated_at_column"() RETURNS "trigger"
    LANGUAGE "plpgsql"
    AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW; 
END;
$$;


--
-- Name: http_request(); Type: FUNCTION; Schema: supabase_functions; Owner: -
--

CREATE FUNCTION "supabase_functions"."http_request"() RETURNS "trigger"
    LANGUAGE "plpgsql" SECURITY DEFINER
    SET "search_path" TO 'supabase_functions'
    AS $$
  DECLARE
    request_id bigint;
    payload jsonb;
    url text := TG_ARGV[0]::text;
    method text := TG_ARGV[1]::text;
    headers jsonb DEFAULT '{}'::jsonb;
    params jsonb DEFAULT '{}'::jsonb;
    timeout_ms integer DEFAULT 1000;
  BEGIN
    IF url IS NULL OR url = 'null' THEN
      RAISE EXCEPTION 'url argument is missing';
    END IF;

    IF method IS NULL OR method = 'null' THEN
      RAISE EXCEPTION 'method argument is missing';
    END IF;

    IF TG_ARGV[2] IS NULL OR TG_ARGV[2] = 'null' THEN
      headers = '{"Content-Type": "application/json"}'::jsonb;
    ELSE
      headers = TG_ARGV[2]::jsonb;
    END IF;

    IF TG_ARGV[3] IS NULL OR TG_ARGV[3] = 'null' THEN
      params = '{}'::jsonb;
    ELSE
      params = TG_ARGV[3]::jsonb;
    END IF;

    IF TG_ARGV[4] IS NULL OR TG_ARGV[4] = 'null' THEN
      timeout_ms = 1000;
    ELSE
      timeout_ms = TG_ARGV[4]::integer;
    END IF;

    CASE
      WHEN method = 'GET' THEN
        SELECT http_get INTO request_id FROM net.http_get(
          url,
          params,
          headers,
          timeout_ms
        );
      WHEN method = 'POST' THEN
        payload = jsonb_build_object(
          'old_record', OLD,
          'record', NEW,
          'type', TG_OP,
          'table', TG_TABLE_NAME,
          'schema', TG_TABLE_SCHEMA
        );

        SELECT http_post INTO request_id FROM net.http_post(
          url,
          payload,
          params,
          headers,
          timeout_ms
        );
      ELSE
        RAISE EXCEPTION 'method argument % is invalid', method;
    END CASE;

    INSERT INTO supabase_functions.hooks
      (hook_table_id, hook_name, request_id)
    VALUES
      (TG_RELID, TG_NAME, request_id);

    RETURN NEW;
  END
$$;


SET default_tablespace = '';

SET default_table_access_method = "heap";

--
-- Name: extensions; Type: TABLE; Schema: _realtime; Owner: -
--

CREATE TABLE "_realtime"."extensions" (
    "id" "uuid" NOT NULL,
    "type" "text",
    "settings" "jsonb",
    "tenant_external_id" "text",
    "inserted_at" timestamp(0) without time zone NOT NULL,
    "updated_at" timestamp(0) without time zone NOT NULL
);


--
-- Name: schema_migrations; Type: TABLE; Schema: _realtime; Owner: -
--

CREATE TABLE "_realtime"."schema_migrations" (
    "version" bigint NOT NULL,
    "inserted_at" timestamp(0) without time zone
);


--
-- Name: tenants; Type: TABLE; Schema: _realtime; Owner: -
--

CREATE TABLE "_realtime"."tenants" (
    "id" "uuid" NOT NULL,
    "name" "text",
    "external_id" "text",
    "jwt_secret" "text",
    "max_concurrent_users" integer DEFAULT 200 NOT NULL,
    "inserted_at" timestamp(0) without time zone NOT NULL,
    "updated_at" timestamp(0) without time zone NOT NULL,
    "max_events_per_second" integer DEFAULT 100 NOT NULL,
    "postgres_cdc_default" "text" DEFAULT 'postgres_cdc_rls'::"text",
    "max_bytes_per_second" integer DEFAULT 100000 NOT NULL,
    "max_channels_per_client" integer DEFAULT 100 NOT NULL,
    "max_joins_per_second" integer DEFAULT 500 NOT NULL,
    "suspend" boolean DEFAULT false,
    "jwt_jwks" "jsonb",
    "notify_private_alpha" boolean DEFAULT false,
    "private_only" boolean DEFAULT false NOT NULL,
    "migrations_ran" integer DEFAULT 0,
    "broadcast_adapter" character varying(255) DEFAULT 'gen_rpc'::character varying,
    "max_presence_events_per_second" integer DEFAULT 10000,
    "max_payload_size_in_kb" integer DEFAULT 3000
);


--
-- Name: audit_log_entries; Type: TABLE; Schema: auth; Owner: -
--

CREATE TABLE "auth"."audit_log_entries" (
    "instance_id" "uuid",
    "id" "uuid" NOT NULL,
    "payload" json,
    "created_at" timestamp with time zone,
    "ip_address" character varying(64) DEFAULT ''::character varying NOT NULL
);


--
-- Name: TABLE "audit_log_entries"; Type: COMMENT; Schema: auth; Owner: -
--

COMMENT ON TABLE "auth"."audit_log_entries" IS 'Auth: Audit trail for user actions.';


--
-- Name: flow_state; Type: TABLE; Schema: auth; Owner: -
--

CREATE TABLE "auth"."flow_state" (
    "id" "uuid" NOT NULL,
    "user_id" "uuid",
    "auth_code" "text" NOT NULL,
    "code_challenge_method" "auth"."code_challenge_method" NOT NULL,
    "code_challenge" "text" NOT NULL,
    "provider_type" "text" NOT NULL,
    "provider_access_token" "text",
    "provider_refresh_token" "text",
    "created_at" timestamp with time zone,
    "updated_at" timestamp with time zone,
    "authentication_method" "text" NOT NULL,
    "auth_code_issued_at" timestamp with time zone
);


--
-- Name: TABLE "flow_state"; Type: COMMENT; Schema: auth; Owner: -
--

COMMENT ON TABLE "auth"."flow_state" IS 'stores metadata for pkce logins';


--
-- Name: identities; Type: TABLE; Schema: auth; Owner: -
--

CREATE TABLE "auth"."identities" (
    "provider_id" "text" NOT NULL,
    "user_id" "uuid" NOT NULL,
    "identity_data" "jsonb" NOT NULL,
    "provider" "text" NOT NULL,
    "last_sign_in_at" timestamp with time zone,
    "created_at" timestamp with time zone,
    "updated_at" timestamp with time zone,
    "email" "text" GENERATED ALWAYS AS ("lower"(("identity_data" ->> 'email'::"text"))) STORED,
    "id" "uuid" DEFAULT "gen_random_uuid"() NOT NULL
);


--
-- Name: TABLE "identities"; Type: COMMENT; Schema: auth; Owner: -
--

COMMENT ON TABLE "auth"."identities" IS 'Auth: Stores identities associated to a user.';


--
-- Name: COLUMN "identities"."email"; Type: COMMENT; Schema: auth; Owner: -
--

COMMENT ON COLUMN "auth"."identities"."email" IS 'Auth: Email is a generated column that references the optional email property in the identity_data';


--
-- Name: instances; Type: TABLE; Schema: auth; Owner: -
--

CREATE TABLE "auth"."instances" (
    "id" "uuid" NOT NULL,
    "uuid" "uuid",
    "raw_base_config" "text",
    "created_at" timestamp with time zone,
    "updated_at" timestamp with time zone
);


--
-- Name: TABLE "instances"; Type: COMMENT; Schema: auth; Owner: -
--

COMMENT ON TABLE "auth"."instances" IS 'Auth: Manages users across multiple sites.';


--
-- Name: mfa_amr_claims; Type: TABLE; Schema: auth; Owner: -
--

CREATE TABLE "auth"."mfa_amr_claims" (
    "session_id" "uuid" NOT NULL,
    "created_at" timestamp with time zone NOT NULL,
    "updated_at" timestamp with time zone NOT NULL,
    "authentication_method" "text" NOT NULL,
    "id" "uuid" NOT NULL
);


--
-- Name: TABLE "mfa_amr_claims"; Type: COMMENT; Schema: auth; Owner: -
--

COMMENT ON TABLE "auth"."mfa_amr_claims" IS 'auth: stores authenticator method reference claims for multi factor authentication';


--
-- Name: mfa_challenges; Type: TABLE; Schema: auth; Owner: -
--

CREATE TABLE "auth"."mfa_challenges" (
    "id" "uuid" NOT NULL,
    "factor_id" "uuid" NOT NULL,
    "created_at" timestamp with time zone NOT NULL,
    "verified_at" timestamp with time zone,
    "ip_address" "inet" NOT NULL,
    "otp_code" "text",
    "web_authn_session_data" "jsonb"
);


--
-- Name: TABLE "mfa_challenges"; Type: COMMENT; Schema: auth; Owner: -
--

COMMENT ON TABLE "auth"."mfa_challenges" IS 'auth: stores metadata about challenge requests made';


--
-- Name: mfa_factors; Type: TABLE; Schema: auth; Owner: -
--

CREATE TABLE "auth"."mfa_factors" (
    "id" "uuid" NOT NULL,
    "user_id" "uuid" NOT NULL,
    "friendly_name" "text",
    "factor_type" "auth"."factor_type" NOT NULL,
    "status" "auth"."factor_status" NOT NULL,
    "created_at" timestamp with time zone NOT NULL,
    "updated_at" timestamp with time zone NOT NULL,
    "secret" "text",
    "phone" "text",
    "last_challenged_at" timestamp with time zone,
    "web_authn_credential" "jsonb",
    "web_authn_aaguid" "uuid"
);


--
-- Name: TABLE "mfa_factors"; Type: COMMENT; Schema: auth; Owner: -
--

COMMENT ON TABLE "auth"."mfa_factors" IS 'auth: stores metadata about factors';


--
-- Name: oauth_authorizations; Type: TABLE; Schema: auth; Owner: -
--

CREATE TABLE "auth"."oauth_authorizations" (
    "id" "uuid" NOT NULL,
    "authorization_id" "text" NOT NULL,
    "client_id" "uuid" NOT NULL,
    "user_id" "uuid",
    "redirect_uri" "text" NOT NULL,
    "scope" "text" NOT NULL,
    "state" "text",
    "resource" "text",
    "code_challenge" "text",
    "code_challenge_method" "auth"."code_challenge_method",
    "response_type" "auth"."oauth_response_type" DEFAULT 'code'::"auth"."oauth_response_type" NOT NULL,
    "status" "auth"."oauth_authorization_status" DEFAULT 'pending'::"auth"."oauth_authorization_status" NOT NULL,
    "authorization_code" "text",
    "created_at" timestamp with time zone DEFAULT "now"() NOT NULL,
    "expires_at" timestamp with time zone DEFAULT ("now"() + '00:03:00'::interval) NOT NULL,
    "approved_at" timestamp with time zone,
    CONSTRAINT "oauth_authorizations_authorization_code_length" CHECK (("char_length"("authorization_code") <= 255)),
    CONSTRAINT "oauth_authorizations_code_challenge_length" CHECK (("char_length"("code_challenge") <= 128)),
    CONSTRAINT "oauth_authorizations_expires_at_future" CHECK (("expires_at" > "created_at")),
    CONSTRAINT "oauth_authorizations_redirect_uri_length" CHECK (("char_length"("redirect_uri") <= 2048)),
    CONSTRAINT "oauth_authorizations_resource_length" CHECK (("char_length"("resource") <= 2048)),
    CONSTRAINT "oauth_authorizations_scope_length" CHECK (("char_length"("scope") <= 4096)),
    CONSTRAINT "oauth_authorizations_state_length" CHECK (("char_length"("state") <= 4096))
);


--
-- Name: oauth_clients; Type: TABLE; Schema: auth; Owner: -
--

CREATE TABLE "auth"."oauth_clients" (
    "id" "uuid" NOT NULL,
    "client_secret_hash" "text",
    "registration_type" "auth"."oauth_registration_type" NOT NULL,
    "redirect_uris" "text" NOT NULL,
    "grant_types" "text" NOT NULL,
    "client_name" "text",
    "client_uri" "text",
    "logo_uri" "text",
    "created_at" timestamp with time zone DEFAULT "now"() NOT NULL,
    "updated_at" timestamp with time zone DEFAULT "now"() NOT NULL,
    "deleted_at" timestamp with time zone,
    "client_type" "auth"."oauth_client_type" DEFAULT 'confidential'::"auth"."oauth_client_type" NOT NULL,
    CONSTRAINT "oauth_clients_client_name_length" CHECK (("char_length"("client_name") <= 1024)),
    CONSTRAINT "oauth_clients_client_uri_length" CHECK (("char_length"("client_uri") <= 2048)),
    CONSTRAINT "oauth_clients_logo_uri_length" CHECK (("char_length"("logo_uri") <= 2048))
);


--
-- Name: oauth_consents; Type: TABLE; Schema: auth; Owner: -
--

CREATE TABLE "auth"."oauth_consents" (
    "id" "uuid" NOT NULL,
    "user_id" "uuid" NOT NULL,
    "client_id" "uuid" NOT NULL,
    "scopes" "text" NOT NULL,
    "granted_at" timestamp with time zone DEFAULT "now"() NOT NULL,
    "revoked_at" timestamp with time zone,
    CONSTRAINT "oauth_consents_revoked_after_granted" CHECK ((("revoked_at" IS NULL) OR ("revoked_at" >= "granted_at"))),
    CONSTRAINT "oauth_consents_scopes_length" CHECK (("char_length"("scopes") <= 2048)),
    CONSTRAINT "oauth_consents_scopes_not_empty" CHECK (("char_length"(TRIM(BOTH FROM "scopes")) > 0))
);


--
-- Name: one_time_tokens; Type: TABLE; Schema: auth; Owner: -
--

CREATE TABLE "auth"."one_time_tokens" (
    "id" "uuid" NOT NULL,
    "user_id" "uuid" NOT NULL,
    "token_type" "auth"."one_time_token_type" NOT NULL,
    "token_hash" "text" NOT NULL,
    "relates_to" "text" NOT NULL,
    "created_at" timestamp without time zone DEFAULT "now"() NOT NULL,
    "updated_at" timestamp without time zone DEFAULT "now"() NOT NULL,
    CONSTRAINT "one_time_tokens_token_hash_check" CHECK (("char_length"("token_hash") > 0))
);


--
-- Name: refresh_tokens; Type: TABLE; Schema: auth; Owner: -
--

CREATE TABLE "auth"."refresh_tokens" (
    "instance_id" "uuid",
    "id" bigint NOT NULL,
    "token" character varying(255),
    "user_id" character varying(255),
    "revoked" boolean,
    "created_at" timestamp with time zone,
    "updated_at" timestamp with time zone,
    "parent" character varying(255),
    "session_id" "uuid"
);


--
-- Name: TABLE "refresh_tokens"; Type: COMMENT; Schema: auth; Owner: -
--

COMMENT ON TABLE "auth"."refresh_tokens" IS 'Auth: Store of tokens used to refresh JWT tokens once they expire.';


--
-- Name: refresh_tokens_id_seq; Type: SEQUENCE; Schema: auth; Owner: -
--

CREATE SEQUENCE "auth"."refresh_tokens_id_seq"
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: refresh_tokens_id_seq; Type: SEQUENCE OWNED BY; Schema: auth; Owner: -
--

ALTER SEQUENCE "auth"."refresh_tokens_id_seq" OWNED BY "auth"."refresh_tokens"."id";


--
-- Name: saml_providers; Type: TABLE; Schema: auth; Owner: -
--

CREATE TABLE "auth"."saml_providers" (
    "id" "uuid" NOT NULL,
    "sso_provider_id" "uuid" NOT NULL,
    "entity_id" "text" NOT NULL,
    "metadata_xml" "text" NOT NULL,
    "metadata_url" "text",
    "attribute_mapping" "jsonb",
    "created_at" timestamp with time zone,
    "updated_at" timestamp with time zone,
    "name_id_format" "text",
    CONSTRAINT "entity_id not empty" CHECK (("char_length"("entity_id") > 0)),
    CONSTRAINT "metadata_url not empty" CHECK ((("metadata_url" = NULL::"text") OR ("char_length"("metadata_url") > 0))),
    CONSTRAINT "metadata_xml not empty" CHECK (("char_length"("metadata_xml") > 0))
);


--
-- Name: TABLE "saml_providers"; Type: COMMENT; Schema: auth; Owner: -
--

COMMENT ON TABLE "auth"."saml_providers" IS 'Auth: Manages SAML Identity Provider connections.';


--
-- Name: saml_relay_states; Type: TABLE; Schema: auth; Owner: -
--

CREATE TABLE "auth"."saml_relay_states" (
    "id" "uuid" NOT NULL,
    "sso_provider_id" "uuid" NOT NULL,
    "request_id" "text" NOT NULL,
    "for_email" "text",
    "redirect_to" "text",
    "created_at" timestamp with time zone,
    "updated_at" timestamp with time zone,
    "flow_state_id" "uuid",
    CONSTRAINT "request_id not empty" CHECK (("char_length"("request_id") > 0))
);


--
-- Name: TABLE "saml_relay_states"; Type: COMMENT; Schema: auth; Owner: -
--

COMMENT ON TABLE "auth"."saml_relay_states" IS 'Auth: Contains SAML Relay State information for each Service Provider initiated login.';


--
-- Name: schema_migrations; Type: TABLE; Schema: auth; Owner: -
--

CREATE TABLE "auth"."schema_migrations" (
    "version" character varying(255) NOT NULL
);


--
-- Name: TABLE "schema_migrations"; Type: COMMENT; Schema: auth; Owner: -
--

COMMENT ON TABLE "auth"."schema_migrations" IS 'Auth: Manages updates to the auth system.';


--
-- Name: sessions; Type: TABLE; Schema: auth; Owner: -
--

CREATE TABLE "auth"."sessions" (
    "id" "uuid" NOT NULL,
    "user_id" "uuid" NOT NULL,
    "created_at" timestamp with time zone,
    "updated_at" timestamp with time zone,
    "factor_id" "uuid",
    "aal" "auth"."aal_level",
    "not_after" timestamp with time zone,
    "refreshed_at" timestamp without time zone,
    "user_agent" "text",
    "ip" "inet",
    "tag" "text",
    "oauth_client_id" "uuid"
);


--
-- Name: TABLE "sessions"; Type: COMMENT; Schema: auth; Owner: -
--

COMMENT ON TABLE "auth"."sessions" IS 'Auth: Stores session data associated to a user.';


--
-- Name: COLUMN "sessions"."not_after"; Type: COMMENT; Schema: auth; Owner: -
--

COMMENT ON COLUMN "auth"."sessions"."not_after" IS 'Auth: Not after is a nullable column that contains a timestamp after which the session should be regarded as expired.';


--
-- Name: sso_domains; Type: TABLE; Schema: auth; Owner: -
--

CREATE TABLE "auth"."sso_domains" (
    "id" "uuid" NOT NULL,
    "sso_provider_id" "uuid" NOT NULL,
    "domain" "text" NOT NULL,
    "created_at" timestamp with time zone,
    "updated_at" timestamp with time zone,
    CONSTRAINT "domain not empty" CHECK (("char_length"("domain") > 0))
);


--
-- Name: TABLE "sso_domains"; Type: COMMENT; Schema: auth; Owner: -
--

COMMENT ON TABLE "auth"."sso_domains" IS 'Auth: Manages SSO email address domain mapping to an SSO Identity Provider.';


--
-- Name: sso_providers; Type: TABLE; Schema: auth; Owner: -
--

CREATE TABLE "auth"."sso_providers" (
    "id" "uuid" NOT NULL,
    "resource_id" "text",
    "created_at" timestamp with time zone,
    "updated_at" timestamp with time zone,
    "disabled" boolean,
    CONSTRAINT "resource_id not empty" CHECK ((("resource_id" = NULL::"text") OR ("char_length"("resource_id") > 0)))
);


--
-- Name: TABLE "sso_providers"; Type: COMMENT; Schema: auth; Owner: -
--

COMMENT ON TABLE "auth"."sso_providers" IS 'Auth: Manages SSO identity provider information; see saml_providers for SAML.';


--
-- Name: COLUMN "sso_providers"."resource_id"; Type: COMMENT; Schema: auth; Owner: -
--

COMMENT ON COLUMN "auth"."sso_providers"."resource_id" IS 'Auth: Uniquely identifies a SSO provider according to a user-chosen resource ID (case insensitive), useful in infrastructure as code.';


--
-- Name: users; Type: TABLE; Schema: auth; Owner: -
--

CREATE TABLE "auth"."users" (
    "instance_id" "uuid",
    "id" "uuid" NOT NULL,
    "aud" character varying(255),
    "role" character varying(255),
    "email" character varying(255),
    "encrypted_password" character varying(255),
    "email_confirmed_at" timestamp with time zone,
    "invited_at" timestamp with time zone,
    "confirmation_token" character varying(255),
    "confirmation_sent_at" timestamp with time zone,
    "recovery_token" character varying(255),
    "recovery_sent_at" timestamp with time zone,
    "email_change_token_new" character varying(255),
    "email_change" character varying(255),
    "email_change_sent_at" timestamp with time zone,
    "last_sign_in_at" timestamp with time zone,
    "raw_app_meta_data" "jsonb",
    "raw_user_meta_data" "jsonb",
    "is_super_admin" boolean,
    "created_at" timestamp with time zone,
    "updated_at" timestamp with time zone,
    "phone" "text" DEFAULT NULL::character varying,
    "phone_confirmed_at" timestamp with time zone,
    "phone_change" "text" DEFAULT ''::character varying,
    "phone_change_token" character varying(255) DEFAULT ''::character varying,
    "phone_change_sent_at" timestamp with time zone,
    "confirmed_at" timestamp with time zone GENERATED ALWAYS AS (LEAST("email_confirmed_at", "phone_confirmed_at")) STORED,
    "email_change_token_current" character varying(255) DEFAULT ''::character varying,
    "email_change_confirm_status" smallint DEFAULT 0,
    "banned_until" timestamp with time zone,
    "reauthentication_token" character varying(255) DEFAULT ''::character varying,
    "reauthentication_sent_at" timestamp with time zone,
    "is_sso_user" boolean DEFAULT false NOT NULL,
    "deleted_at" timestamp with time zone,
    "is_anonymous" boolean DEFAULT false NOT NULL,
    CONSTRAINT "users_email_change_confirm_status_check" CHECK ((("email_change_confirm_status" >= 0) AND ("email_change_confirm_status" <= 2)))
);


--
-- Name: TABLE "users"; Type: COMMENT; Schema: auth; Owner: -
--

COMMENT ON TABLE "auth"."users" IS 'Auth: Stores user login data within a secure schema.';


--
-- Name: COLUMN "users"."is_sso_user"; Type: COMMENT; Schema: auth; Owner: -
--

COMMENT ON COLUMN "auth"."users"."is_sso_user" IS 'Auth: Set this column to true when the account comes from SSO. These accounts can have duplicate emails.';


--
-- Name: adaptive_streaming_rules; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE "public"."adaptive_streaming_rules" (
    "id" "uuid" DEFAULT "gen_random_uuid"() NOT NULL,
    "rule_name" character varying(100) NOT NULL,
    "condition_bandwidth_min" numeric(10,2),
    "condition_bandwidth_max" numeric(10,2),
    "condition_latency_max" integer,
    "condition_packet_loss_max" numeric(5,2),
    "condition_connection_types" "text"[],
    "condition_device_types" "text"[],
    "recommended_quality" character varying(10) NOT NULL,
    "buffer_target_seconds" integer DEFAULT 10,
    "preload_enabled" boolean DEFAULT true,
    "active" boolean DEFAULT true,
    "priority" integer DEFAULT 50,
    "created_at" timestamp without time zone DEFAULT "now"(),
    "updated_at" timestamp without time zone DEFAULT "now"()
);


--
-- Name: audit_logs; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE "public"."audit_logs" (
    "id" "uuid" DEFAULT "gen_random_uuid"() NOT NULL,
    "action" character varying(50) NOT NULL,
    "user_id" "uuid" NOT NULL,
    "course_id" "uuid",
    "lecture_id" "uuid",
    "transaction_id" "uuid",
    "details" "jsonb",
    "ip_address" "inet",
    "user_agent" "text",
    "created_at" timestamp without time zone DEFAULT "now"() NOT NULL
);


--
-- Name: TABLE "audit_logs"; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TABLE "public"."audit_logs" IS 'General audit logging for payment and access events';


--
-- Name: bandwidth_tests; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE "public"."bandwidth_tests" (
    "id" "uuid" DEFAULT "gen_random_uuid"() NOT NULL,
    "session_id" character varying(255) NOT NULL,
    "user_id" "uuid" NOT NULL,
    "test_type" character varying(20) NOT NULL,
    "download_mbps" numeric(10,2),
    "upload_mbps" numeric(10,2),
    "latency_ms" integer,
    "jitter_ms" integer,
    "packet_loss_percent" numeric(5,4),
    "test_duration_seconds" integer,
    "test_server_location" character varying(100),
    "confidence_score" numeric(3,2),
    "timestamp" timestamp without time zone DEFAULT "now"(),
    "created_at" timestamp without time zone DEFAULT "now"()
);


--
-- Name: chat_history; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE "public"."chat_history" (
    "id" "uuid" DEFAULT "gen_random_uuid"() NOT NULL,
    "user_id" "uuid",
    "message" "text" NOT NULL,
    "is_user" boolean NOT NULL,
    "created_at" timestamp without time zone DEFAULT "now"() NOT NULL
);


--
-- Name: course_access_logs; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE "public"."course_access_logs" (
    "id" "uuid" DEFAULT "gen_random_uuid"() NOT NULL,
    "user_id" "uuid" NOT NULL,
    "course_id" "uuid" NOT NULL,
    "lecture_id" "uuid",
    "access_type" character varying(20) NOT NULL,
    "access_granted" boolean NOT NULL,
    "payment_required" boolean NOT NULL,
    "payment_verified" boolean DEFAULT false,
    "transaction_id" "uuid",
    "client_ip" "inet",
    "user_agent" "text",
    "access_duration_seconds" integer,
    "created_at" timestamp without time zone DEFAULT "now"() NOT NULL,
    CONSTRAINT "course_access_logs_access_type_check" CHECK ((("access_type")::"text" = ANY (ARRAY[('full'::character varying)::"text", ('preview'::character varying)::"text", ('denied'::character varying)::"text"])))
);


--
-- Name: TABLE "course_access_logs"; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TABLE "public"."course_access_logs" IS 'Comprehensive audit trail of all course access attempts';


--
-- Name: courses; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE "public"."courses" (
    "id" "uuid" DEFAULT "gen_random_uuid"() NOT NULL,
    "title" character varying(255) NOT NULL,
    "description" "text" NOT NULL,
    "instructor_id" "uuid" NOT NULL,
    "instructor_name" character varying(255) NOT NULL,
    "category" character varying(100) NOT NULL,
    "level" "public"."course_level" DEFAULT 'beginner'::"public"."course_level" NOT NULL,
    "price" numeric(10,2) DEFAULT 0 NOT NULL,
    "currency" character varying(10) DEFAULT 'USD'::character varying NOT NULL,
    "thumbnail_url" "text",
    "status" "public"."course_status" DEFAULT 'draft'::"public"."course_status" NOT NULL,
    "duration_minutes" integer DEFAULT 0 NOT NULL,
    "enrollment_count" integer DEFAULT 0 NOT NULL,
    "rating" numeric(3,2) DEFAULT 0 NOT NULL,
    "rating_count" integer DEFAULT 0 NOT NULL,
    "tags" "text"[],
    "created_at" timestamp without time zone DEFAULT "now"() NOT NULL,
    "updated_at" timestamp without time zone DEFAULT "now"() NOT NULL,
    "lemon_squeezy_product_id" character varying(100),
    "lemon_squeezy_variant_id" character varying(100),
    "is_paid" boolean DEFAULT false,
    "deleted_at" timestamp without time zone,
    "learning_outcomes" "text"[],
    "requirements" "text"[],
    "language" character varying(10) DEFAULT 'en'::character varying,
    "auto_approve_enrollment" boolean DEFAULT true,
    "allow_previews" boolean DEFAULT true,
    "has_certificate" boolean DEFAULT false,
    "mobile_access" boolean DEFAULT true,
    "difficulty_level" character varying(20) DEFAULT 'intermediate'::character varying,
    "estimated_duration_hours" integer DEFAULT 0,
    "preview_enabled" boolean DEFAULT true,
    "preview_duration_minutes" integer DEFAULT 10,
    "requires_enrollment_approval" boolean DEFAULT false
);


--
-- Name: enrollments; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE "public"."enrollments" (
    "id" "uuid" DEFAULT "gen_random_uuid"() NOT NULL,
    "user_id" "uuid" NOT NULL,
    "course_id" "uuid" NOT NULL,
    "status" "public"."enrollment_status" DEFAULT 'enrolled'::"public"."enrollment_status" NOT NULL,
    "progress_percentage" numeric(5,2) DEFAULT 0 NOT NULL,
    "enrolled_at" timestamp without time zone DEFAULT "now"() NOT NULL,
    "completed_at" timestamp without time zone,
    "last_accessed" timestamp without time zone,
    "completed_lectures" integer DEFAULT 0,
    "total_lectures" integer DEFAULT 0,
    "total_watch_time_seconds" integer DEFAULT 0,
    "created_at" timestamp without time zone DEFAULT "now"(),
    "updated_at" timestamp without time zone DEFAULT "now"(),
    "payment_required" boolean DEFAULT false,
    "payment_status" character varying(50) DEFAULT 'pending'::character varying,
    "lemon_squeezy_order_id" character varying(255),
    "payment_amount" numeric(10,2),
    "payment_currency" character varying(10),
    "paid_at" timestamp without time zone,
    "deleted_at" timestamp without time zone,
    "payment_verified_at" timestamp without time zone,
    "access_expires_at" timestamp without time zone,
    "transaction_id" "uuid"
);


--
-- Name: COLUMN "enrollments"."payment_status"; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN "public"."enrollments"."payment_status" IS 'Current payment status: pending, paid, refunded, expired';


--
-- Name: COLUMN "enrollments"."payment_verified_at"; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN "public"."enrollments"."payment_verified_at" IS 'Timestamp when payment was verified by provider';


--
-- Name: COLUMN "enrollments"."access_expires_at"; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN "public"."enrollments"."access_expires_at" IS 'When access expires (for subscription-based courses)';


--
-- Name: transactions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE "public"."transactions" (
    "id" "uuid" DEFAULT "gen_random_uuid"() NOT NULL,
    "user_id" "uuid",
    "course_id" "uuid",
    "payment_method_id" "uuid",
    "amount" numeric(10,2) NOT NULL,
    "currency" character varying(3) NOT NULL,
    "status" character varying(20) NOT NULL,
    "transaction_reference" character varying(100),
    "created_at" timestamp without time zone DEFAULT "now"() NOT NULL,
    "updated_at" timestamp without time zone DEFAULT "now"() NOT NULL,
    "lemon_squeezy_order_id" character varying(100),
    "lemon_squeezy_checkout_id" character varying(100),
    "webhook_event_id" character varying(100),
    "custom_data" "jsonb",
    "stripe_payment_intent_id" character varying(100),
    "stripe_customer_id" character varying(100),
    "stripe_charge_id" character varying(100),
    "stripe_session_id" character varying(100),
    "stripe_invoice_id" character varying(100),
    "stripe_subscription_id" character varying(100),
    "payment_provider" character varying(50) DEFAULT 'lemonsqueezy'::character varying,
    "payment_verified_at" timestamp without time zone,
    "refunded_at" timestamp without time zone,
    "expires_at" timestamp without time zone,
    "payment_method_details" "jsonb",
    "risk_score" numeric(3,2) DEFAULT 0.0,
    CONSTRAINT "chk_transactions_stripe_customer_id" CHECK (((("stripe_customer_id")::"text" ~ '^cus_[a-zA-Z0-9]+$'::"text") OR ("stripe_customer_id" IS NULL))),
    CONSTRAINT "chk_transactions_stripe_payment_intent_id" CHECK (((("stripe_payment_intent_id")::"text" ~ '^pi_[a-zA-Z0-9]+$'::"text") OR ("stripe_payment_intent_id" IS NULL)))
);


--
-- Name: COLUMN "transactions"."lemon_squeezy_order_id"; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN "public"."transactions"."lemon_squeezy_order_id" IS 'LemonSqueezy order ID for alternative payment processor';


--
-- Name: COLUMN "transactions"."lemon_squeezy_checkout_id"; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN "public"."transactions"."lemon_squeezy_checkout_id" IS 'LemonSqueezy checkout ID for alternative payment processor';


--
-- Name: COLUMN "transactions"."custom_data"; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN "public"."transactions"."custom_data" IS 'Additional payment metadata in JSON format';


--
-- Name: COLUMN "transactions"."stripe_payment_intent_id"; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN "public"."transactions"."stripe_payment_intent_id" IS 'Stripe payment intent ID for tracking payments';


--
-- Name: COLUMN "transactions"."stripe_customer_id"; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN "public"."transactions"."stripe_customer_id" IS 'Stripe customer ID associated with the transaction';


--
-- Name: COLUMN "transactions"."stripe_charge_id"; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN "public"."transactions"."stripe_charge_id" IS 'Stripe charge ID for completed payments';


--
-- Name: COLUMN "transactions"."payment_verified_at"; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN "public"."transactions"."payment_verified_at" IS 'When payment was verified by webhook or API';


--
-- Name: COLUMN "transactions"."risk_score"; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN "public"."transactions"."risk_score" IS 'Fraud risk score from 0.0 (safe) to 1.0 (high risk)';


--
-- Name: course_access_analytics; Type: VIEW; Schema: public; Owner: -
--

CREATE VIEW "public"."course_access_analytics" AS
 SELECT "c"."id" AS "course_id",
    "c"."title" AS "course_title",
    "c"."is_paid",
    "count"(DISTINCT "cal"."user_id") AS "total_access_attempts",
    "count"(DISTINCT
        CASE
            WHEN ("cal"."access_granted" = true) THEN "cal"."user_id"
            ELSE NULL::"uuid"
        END) AS "successful_accesses",
    "count"(DISTINCT
        CASE
            WHEN (("cal"."access_type")::"text" = 'preview'::"text") THEN "cal"."user_id"
            ELSE NULL::"uuid"
        END) AS "preview_accesses",
    "count"(DISTINCT "e"."user_id") AS "total_enrollments",
    "count"(DISTINCT
        CASE
            WHEN (("e"."payment_status")::"text" = 'paid'::"text") THEN "e"."user_id"
            ELSE NULL::"uuid"
        END) AS "paid_enrollments",
    "avg"("cal"."access_duration_seconds") AS "avg_access_duration",
    "sum"(
        CASE
            WHEN (("t"."status")::"text" = 'completed'::"text") THEN "t"."amount"
            ELSE (0)::numeric
        END) AS "total_revenue"
   FROM ((("public"."courses" "c"
     LEFT JOIN "public"."course_access_logs" "cal" ON (("c"."id" = "cal"."course_id")))
     LEFT JOIN "public"."enrollments" "e" ON (("c"."id" = "e"."course_id")))
     LEFT JOIN "public"."transactions" "t" ON ((("c"."id" = "t"."course_id") AND (("t"."status")::"text" = 'completed'::"text"))))
  GROUP BY "c"."id", "c"."title", "c"."is_paid"
  ORDER BY ("sum"(
        CASE
            WHEN (("t"."status")::"text" = 'completed'::"text") THEN "t"."amount"
            ELSE (0)::numeric
        END)) DESC;


--
-- Name: course_access_cache; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE "public"."course_access_cache" (
    "id" "uuid" DEFAULT "gen_random_uuid"() NOT NULL,
    "user_id" "uuid" NOT NULL,
    "course_id" "uuid" NOT NULL,
    "access_level" character varying(20) NOT NULL,
    "payment_verified" boolean NOT NULL,
    "transaction_id" "uuid",
    "expires_at" timestamp without time zone NOT NULL,
    "created_at" timestamp without time zone DEFAULT "now"() NOT NULL,
    "updated_at" timestamp without time zone DEFAULT "now"() NOT NULL,
    CONSTRAINT "check_cache_expires_future" CHECK (("expires_at" > "created_at")),
    CONSTRAINT "course_access_cache_access_level_check" CHECK ((("access_level")::"text" = ANY (ARRAY[('full'::character varying)::"text", ('preview'::character varying)::"text", ('denied'::character varying)::"text"])))
);


--
-- Name: TABLE "course_access_cache"; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TABLE "public"."course_access_cache" IS 'Performance cache for course access validation results';


--
-- Name: course_resources; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE "public"."course_resources" (
    "id" "uuid" DEFAULT "gen_random_uuid"() NOT NULL,
    "course_id" "uuid" NOT NULL,
    "file_id" "uuid" NOT NULL,
    "resource_type" character varying(50) DEFAULT 'document'::character varying NOT NULL,
    "display_order" integer DEFAULT 1 NOT NULL,
    "is_required" boolean DEFAULT false NOT NULL,
    "created_at" timestamp without time zone DEFAULT "now"() NOT NULL,
    "updated_at" timestamp without time zone DEFAULT "now"() NOT NULL
);


--
-- Name: TABLE "course_resources"; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TABLE "public"."course_resources" IS 'DEPRECATED: Use lecture_resources instead. Resources are now attached to individual lectures.';


--
-- Name: file_permissions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE "public"."file_permissions" (
    "id" "uuid" DEFAULT "gen_random_uuid"() NOT NULL,
    "file_id" "uuid",
    "user_id" "uuid",
    "permission_type" character varying(20) NOT NULL,
    "granted_by" "uuid" NOT NULL,
    "created_at" timestamp without time zone DEFAULT "now"(),
    CONSTRAINT "file_permissions_permission_type_check" CHECK ((("permission_type")::"text" = ANY (ARRAY[('read'::character varying)::"text", ('write'::character varying)::"text", ('delete'::character varying)::"text"])))
);


--
-- Name: files; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE "public"."files" (
    "id" "uuid" DEFAULT "gen_random_uuid"() NOT NULL,
    "filename" character varying(255) NOT NULL,
    "original_filename" character varying(255) NOT NULL,
    "content_type" character varying(100) NOT NULL,
    "size_bytes" bigint NOT NULL,
    "bucket_name" character varying(100) NOT NULL,
    "object_key" character varying(500) NOT NULL,
    "upload_user_id" "uuid" NOT NULL,
    "is_public" boolean DEFAULT false,
    "metadata" "jsonb",
    "checksum" character varying(64),
    "thumbnail_url" character varying(500),
    "created_at" timestamp without time zone DEFAULT "now"(),
    "updated_at" timestamp without time zone DEFAULT "now"(),
    "deleted_at" timestamp without time zone
);


--
-- Name: forum_mentions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE "public"."forum_mentions" (
    "id" "uuid" DEFAULT "gen_random_uuid"() NOT NULL,
    "post_id" "uuid" NOT NULL,
    "mentioned_user_id" "uuid" NOT NULL,
    "mentioner_user_id" "uuid" NOT NULL,
    "is_read" boolean DEFAULT false,
    "created_at" timestamp without time zone DEFAULT "now"() NOT NULL
);


--
-- Name: TABLE "forum_mentions"; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TABLE "public"."forum_mentions" IS 'Tracks @username mentions in forum posts';


--
-- Name: forum_notifications; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE "public"."forum_notifications" (
    "id" "uuid" DEFAULT "gen_random_uuid"() NOT NULL,
    "user_id" "uuid" NOT NULL,
    "type" character varying(50) NOT NULL,
    "title" character varying(255) NOT NULL,
    "message" "text" NOT NULL,
    "reference_id" "uuid",
    "reference_type" character varying(50),
    "is_read" boolean DEFAULT false,
    "created_at" timestamp without time zone DEFAULT "now"() NOT NULL,
    CONSTRAINT "forum_notifications_reference_type_check" CHECK ((("reference_type")::"text" = ANY (ARRAY[('post'::character varying)::"text", ('topic'::character varying)::"text", ('mention'::character varying)::"text"]))),
    CONSTRAINT "forum_notifications_type_check" CHECK ((("type")::"text" = ANY (ARRAY[('mention'::character varying)::"text", ('topic_approved'::character varying)::"text", ('post_approved'::character varying)::"text", ('topic_reply'::character varying)::"text"])))
);


--
-- Name: TABLE "forum_notifications"; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TABLE "public"."forum_notifications" IS 'User notifications for forum events (mentions, approvals, replies)';


--
-- Name: forum_posts; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE "public"."forum_posts" (
    "id" "uuid" DEFAULT "gen_random_uuid"() NOT NULL,
    "topic_id" "uuid",
    "author_id" "uuid",
    "content" "text" NOT NULL,
    "created_at" timestamp without time zone DEFAULT "now"() NOT NULL,
    "updated_at" timestamp without time zone DEFAULT "now"() NOT NULL,
    "status" character varying(20) DEFAULT 'pending'::character varying,
    "pin_order" integer,
    "parent_id" "uuid",
    "is_edited" boolean DEFAULT false,
    "edited_at" timestamp without time zone,
    "up_votes" integer DEFAULT 0,
    "down_votes" integer DEFAULT 0,
    "is_answer" boolean DEFAULT false,
    "is_pinned" boolean DEFAULT false,
    "deleted_at" timestamp without time zone,
    CONSTRAINT "forum_posts_status_check" CHECK ((("status")::"text" = ANY (ARRAY[('pending'::character varying)::"text", ('approved'::character varying)::"text", ('rejected'::character varying)::"text"])))
);


--
-- Name: COLUMN "forum_posts"."status"; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN "public"."forum_posts"."status" IS 'Approval status: pending (needs approval), approved (visible), rejected (hidden)';


--
-- Name: COLUMN "forum_posts"."pin_order"; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN "public"."forum_posts"."pin_order" IS 'Order for pinned posts within a topic (lower numbers appear first)';


--
-- Name: forum_topic_subscriptions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE "public"."forum_topic_subscriptions" (
    "id" "uuid" DEFAULT "gen_random_uuid"() NOT NULL,
    "user_id" "uuid" NOT NULL,
    "topic_id" "uuid" NOT NULL,
    "subscribed_at" timestamp without time zone DEFAULT "now"() NOT NULL
);


--
-- Name: TABLE "forum_topic_subscriptions"; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TABLE "public"."forum_topic_subscriptions" IS 'User subscriptions to forum topics for notifications';


--
-- Name: forum_topics; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE "public"."forum_topics" (
    "id" "uuid" DEFAULT "gen_random_uuid"() NOT NULL,
    "title" character varying(200) NOT NULL,
    "course_id" "uuid",
    "creator_id" "uuid",
    "is_pinned" boolean DEFAULT false,
    "is_locked" boolean DEFAULT false,
    "created_at" timestamp without time zone DEFAULT "now"() NOT NULL,
    "updated_at" timestamp without time zone DEFAULT "now"() NOT NULL,
    "deleted_at" timestamp without time zone,
    "created_by_id" "uuid",
    "description" "text",
    "category" character varying(100),
    "tags" "text"[],
    "is_sticky" boolean DEFAULT false,
    "view_count" integer DEFAULT 0,
    "post_count" integer DEFAULT 0,
    "last_post_at" timestamp without time zone,
    "last_post_by_id" "uuid",
    "status" character varying(20) DEFAULT 'pending'::character varying,
    "pin_order" integer,
    CONSTRAINT "forum_topics_status_check" CHECK ((("status")::"text" = ANY (ARRAY[('pending'::character varying)::"text", ('approved'::character varying)::"text", ('rejected'::character varying)::"text"])))
);


--
-- Name: COLUMN "forum_topics"."status"; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN "public"."forum_topics"."status" IS 'Approval status: pending (needs approval), approved (visible), rejected (hidden)';


--
-- Name: COLUMN "forum_topics"."pin_order"; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN "public"."forum_topics"."pin_order" IS 'Order for pinned topics (lower numbers appear first)';


--
-- Name: forum_votes; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE "public"."forum_votes" (
    "id" "uuid" DEFAULT "gen_random_uuid"() NOT NULL,
    "post_id" "uuid" NOT NULL,
    "user_id" "uuid" NOT NULL,
    "vote_type" character varying(10) NOT NULL,
    "created_at" timestamp without time zone DEFAULT "now"() NOT NULL,
    "updated_at" timestamp without time zone DEFAULT "now"() NOT NULL,
    "voted_at" timestamp without time zone DEFAULT "now"(),
    CONSTRAINT "forum_votes_vote_type_check" CHECK ((("vote_type")::"text" = ANY (ARRAY[('up'::character varying)::"text", ('down'::character varying)::"text"])))
);


--
-- Name: TABLE "forum_votes"; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TABLE "public"."forum_votes" IS 'User votes (upvotes/downvotes) on forum posts';


--
-- Name: COLUMN "forum_votes"."vote_type"; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN "public"."forum_votes"."vote_type" IS 'Type of vote: up or down';


--
-- Name: lecture_preview_sessions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE "public"."lecture_preview_sessions" (
    "id" "uuid" DEFAULT "gen_random_uuid"() NOT NULL,
    "user_id" "uuid" NOT NULL,
    "lecture_id" "uuid" NOT NULL,
    "session_started_at" timestamp without time zone DEFAULT "now"() NOT NULL,
    "session_duration_seconds" integer DEFAULT 0 NOT NULL,
    "preview_limit_seconds" integer DEFAULT 600 NOT NULL,
    "preview_exhausted" boolean DEFAULT false,
    "ip_address" "inet",
    "last_accessed_at" timestamp without time zone DEFAULT "now"(),
    "created_at" timestamp without time zone DEFAULT "now"() NOT NULL,
    "updated_at" timestamp without time zone DEFAULT "now"() NOT NULL,
    CONSTRAINT "check_preview_limit_positive" CHECK (("preview_limit_seconds" > 0)),
    CONSTRAINT "check_session_duration_valid" CHECK ((("session_duration_seconds" >= 0) AND ("session_duration_seconds" <= "preview_limit_seconds")))
);


--
-- Name: TABLE "lecture_preview_sessions"; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TABLE "public"."lecture_preview_sessions" IS 'Tracking table for lecture preview sessions and time limits';


--
-- Name: lecture_resources; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE "public"."lecture_resources" (
    "id" "uuid" DEFAULT "gen_random_uuid"() NOT NULL,
    "lecture_id" "uuid" NOT NULL,
    "file_id" "uuid" NOT NULL,
    "resource_type" character varying(50) DEFAULT 'document'::character varying NOT NULL,
    "display_order" integer DEFAULT 1 NOT NULL,
    "is_required" boolean DEFAULT false NOT NULL,
    "created_at" timestamp without time zone DEFAULT "now"() NOT NULL,
    "updated_at" timestamp without time zone DEFAULT "now"() NOT NULL,
    CONSTRAINT "check_resource_type" CHECK ((("resource_type")::"text" = ANY (ARRAY[('document'::character varying)::"text", ('pdf'::character varying)::"text", ('video'::character varying)::"text", ('audio'::character varying)::"text", ('image'::character varying)::"text", ('archive'::character varying)::"text", ('code'::character varying)::"text", ('slides'::character varying)::"text", ('worksheet'::character varying)::"text", ('quiz'::character varying)::"text", ('assignment'::character varying)::"text", ('other'::character varying)::"text"])))
);


--
-- Name: TABLE "lecture_resources"; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TABLE "public"."lecture_resources" IS 'Resources attached to individual lectures instead of courses';


--
-- Name: COLUMN "lecture_resources"."lecture_id"; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN "public"."lecture_resources"."lecture_id" IS 'Reference to the lecture this resource belongs to';


--
-- Name: COLUMN "lecture_resources"."file_id"; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN "public"."lecture_resources"."file_id" IS 'Reference to the file in bucket service';


--
-- Name: COLUMN "lecture_resources"."resource_type"; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN "public"."lecture_resources"."resource_type" IS 'Type of resource (document, pdf, video, etc.)';


--
-- Name: COLUMN "lecture_resources"."display_order"; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN "public"."lecture_resources"."display_order" IS 'Order in which resources should be displayed';


--
-- Name: COLUMN "lecture_resources"."is_required"; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN "public"."lecture_resources"."is_required" IS 'Whether this resource is required for course completion';


--
-- Name: lectures; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE "public"."lectures" (
    "id" "uuid" DEFAULT "gen_random_uuid"() NOT NULL,
    "course_id" "uuid" NOT NULL,
    "title" character varying(255) NOT NULL,
    "description" "text",
    "order_number" integer NOT NULL,
    "duration_minutes" integer DEFAULT 0 NOT NULL,
    "video_url" "text",
    "video_id" character varying(255),
    "status" "public"."lecture_status" DEFAULT 'draft'::"public"."lecture_status" NOT NULL,
    "is_free" boolean DEFAULT false NOT NULL,
    "created_at" timestamp without time zone DEFAULT "now"() NOT NULL,
    "updated_at" timestamp without time zone DEFAULT "now"() NOT NULL,
    "deleted_at" timestamp without time zone,
    "type" character varying(50) DEFAULT 'video'::character varying,
    "lecture_type" character varying(50) DEFAULT 'video'::character varying,
    "preview_available" boolean DEFAULT false,
    "access_level" character varying(20) DEFAULT 'paid'::character varying,
    CONSTRAINT "lectures_access_level_check" CHECK ((("access_level")::"text" = ANY (ARRAY[('free'::character varying)::"text", ('paid'::character varying)::"text", ('preview'::character varying)::"text"])))
);


--
-- Name: lemon_squeezy_products; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE "public"."lemon_squeezy_products" (
    "id" "uuid" DEFAULT "gen_random_uuid"() NOT NULL,
    "lemon_squeezy_product_id" character varying(100) NOT NULL,
    "name" character varying(255) NOT NULL,
    "description" "text",
    "status" character varying(50) NOT NULL,
    "store_id" character varying(100) NOT NULL,
    "created_at" timestamp without time zone DEFAULT "now"() NOT NULL,
    "updated_at" timestamp without time zone DEFAULT "now"() NOT NULL
);


--
-- Name: lemon_squeezy_variants; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE "public"."lemon_squeezy_variants" (
    "id" "uuid" DEFAULT "gen_random_uuid"() NOT NULL,
    "lemon_squeezy_variant_id" character varying(100) NOT NULL,
    "lemon_squeezy_product_id" character varying(100) NOT NULL,
    "name" character varying(255) NOT NULL,
    "price" numeric(10,2) NOT NULL,
    "currency" character varying(3) DEFAULT 'USD'::character varying NOT NULL,
    "status" character varying(50) NOT NULL,
    "created_at" timestamp without time zone DEFAULT "now"() NOT NULL,
    "updated_at" timestamp without time zone DEFAULT "now"() NOT NULL
);


--
-- Name: lemon_squeezy_webhook_events; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE "public"."lemon_squeezy_webhook_events" (
    "id" "uuid" DEFAULT "gen_random_uuid"() NOT NULL,
    "event_id" character varying(100) NOT NULL,
    "event_name" character varying(50) NOT NULL,
    "processed_at" timestamp without time zone,
    "payload" "jsonb" NOT NULL,
    "signature" character varying(255) NOT NULL,
    "created_at" timestamp without time zone DEFAULT "now"() NOT NULL
);


--
-- Name: network_analytics_daily; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE "public"."network_analytics_daily" (
    "id" "uuid" DEFAULT "gen_random_uuid"() NOT NULL,
    "date" "date" NOT NULL,
    "user_id" "uuid",
    "video_id" "uuid",
    "session_count" integer DEFAULT 0,
    "avg_bandwidth_mbps" numeric(10,2),
    "avg_latency_ms" numeric(8,2),
    "avg_packet_loss" numeric(5,4),
    "quality_distribution" "jsonb",
    "connection_type_distribution" "jsonb",
    "device_type_distribution" "jsonb",
    "total_quality_changes" integer DEFAULT 0,
    "total_buffer_events" integer DEFAULT 0,
    "total_network_interruptions" integer DEFAULT 0,
    "avg_stability_score" numeric(4,2),
    "created_at" timestamp without time zone DEFAULT "now"(),
    "updated_at" timestamp without time zone DEFAULT "now"()
);


--
-- Name: network_events; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE "public"."network_events" (
    "id" "uuid" DEFAULT "gen_random_uuid"() NOT NULL,
    "session_id" character varying(255) NOT NULL,
    "user_id" "uuid" NOT NULL,
    "event_type" character varying(50) NOT NULL,
    "event_data" "jsonb",
    "severity" character varying(20) DEFAULT 'info'::character varying,
    "timestamp" timestamp without time zone DEFAULT "now"(),
    "resolved" boolean DEFAULT false,
    "resolution_timestamp" timestamp without time zone,
    "created_at" timestamp without time zone DEFAULT "now"()
);


--
-- Name: network_metrics; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE "public"."network_metrics" (
    "id" "uuid" DEFAULT "gen_random_uuid"() NOT NULL,
    "session_id" character varying(255) NOT NULL,
    "user_id" "uuid" NOT NULL,
    "timestamp" timestamp without time zone DEFAULT "now"(),
    "bandwidth_mbps" numeric(10,2),
    "latency_ms" integer,
    "packet_loss_percent" numeric(5,2),
    "connection_type" character varying(20),
    "quality_score" integer,
    "recommended_quality" character varying(10),
    "buffer_health_seconds" integer,
    "created_at" timestamp without time zone DEFAULT "now"(),
    "device_type" character varying(50),
    "screen_resolution" character varying(20),
    "user_agent" "text",
    "geographic_location" "jsonb",
    "isp_info" "jsonb",
    "network_stability_score" integer DEFAULT 5,
    "jitter_ms" integer DEFAULT 0,
    "download_speed_mbps" numeric(10,2),
    "upload_speed_mbps" numeric(10,2),
    "cpu_usage_percent" numeric(5,2),
    "memory_usage_percent" numeric(5,2),
    "battery_level" integer,
    "thermal_state" character varying(20)
);


--
-- Name: viewing_sessions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE "public"."viewing_sessions" (
    "id" "uuid" DEFAULT "gen_random_uuid"() NOT NULL,
    "session_id" character varying(255) NOT NULL,
    "user_id" "uuid" NOT NULL,
    "video_id" "uuid",
    "started_at" timestamp without time zone DEFAULT "now"(),
    "last_heartbeat" timestamp without time zone DEFAULT "now"(),
    "current_time_seconds" integer DEFAULT 0,
    "current_quality" character varying(10),
    "total_watch_time_seconds" integer DEFAULT 0,
    "completed" boolean DEFAULT false,
    "user_agent" "text",
    "ip_address" "inet",
    "created_at" timestamp without time zone DEFAULT "now"()
);


--
-- Name: network_monitoring_dashboard; Type: MATERIALIZED VIEW; Schema: public; Owner: -
--

CREATE MATERIALIZED VIEW "public"."network_monitoring_dashboard" AS
 SELECT "nm"."session_id",
    "nm"."user_id",
    "vs"."video_id",
    "nm"."timestamp",
    "nm"."bandwidth_mbps",
    "nm"."latency_ms",
    "nm"."packet_loss_percent",
    "nm"."connection_type",
    "nm"."device_type",
    "nm"."quality_score",
    "nm"."recommended_quality",
    "nm"."buffer_health_seconds",
    "nm"."network_stability_score",
    "vs"."current_quality",
    "vs"."current_time_seconds",
        CASE
            WHEN ("nm"."timestamp" > ("now"() - '00:05:00'::interval)) THEN 'active'::"text"
            WHEN ("nm"."timestamp" > ("now"() - '00:30:00'::interval)) THEN 'recent'::"text"
            ELSE 'inactive'::"text"
        END AS "session_status",
    EXTRACT(epoch FROM ("now"() - ("nm"."timestamp")::timestamp with time zone)) AS "seconds_since_last_update"
   FROM ("public"."network_metrics" "nm"
     LEFT JOIN "public"."viewing_sessions" "vs" ON ((("nm"."session_id")::"text" = ("vs"."session_id")::"text")))
  WHERE ("nm"."timestamp" > ("now"() - '04:00:00'::interval))
  WITH NO DATA;


--
-- Name: notes; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE "public"."notes" (
    "id" "uuid" DEFAULT "gen_random_uuid"() NOT NULL,
    "user_id" "uuid" NOT NULL,
    "course_id" "uuid" NOT NULL,
    "lecture_id" "uuid" NOT NULL,
    "title" character varying(500) NOT NULL,
    "content" "text" NOT NULL,
    "timestamp_seconds" integer DEFAULT 0,
    "created_at" timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


--
-- Name: oauth_accounts; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE "public"."oauth_accounts" (
    "id" "uuid" DEFAULT "gen_random_uuid"() NOT NULL,
    "user_id" "uuid",
    "provider" character varying(50) NOT NULL,
    "provider_id" character varying(100) NOT NULL,
    "provider_email" character varying(100),
    "access_token" "text",
    "refresh_token" "text",
    "expires_at" timestamp without time zone,
    "created_at" timestamp without time zone DEFAULT "now"() NOT NULL,
    "updated_at" timestamp without time zone DEFAULT "now"() NOT NULL
);


--
-- Name: payment_analytics; Type: VIEW; Schema: public; Owner: -
--

CREATE VIEW "public"."payment_analytics" AS
 SELECT "date_trunc"('day'::"text", "created_at") AS "payment_date",
    "count"(*) AS "total_transactions",
    "count"(
        CASE
            WHEN (("status")::"text" = 'completed'::"text") THEN 1
            ELSE NULL::integer
        END) AS "successful_payments",
    "count"(
        CASE
            WHEN (("status")::"text" = 'failed'::"text") THEN 1
            ELSE NULL::integer
        END) AS "failed_payments",
    "count"(
        CASE
            WHEN (("status")::"text" = 'refunded'::"text") THEN 1
            ELSE NULL::integer
        END) AS "refunded_payments",
    "sum"(
        CASE
            WHEN (("status")::"text" = 'completed'::"text") THEN "amount"
            ELSE (0)::numeric
        END) AS "total_revenue",
    "avg"(
        CASE
            WHEN (("status")::"text" = 'completed'::"text") THEN "amount"
            ELSE NULL::numeric
        END) AS "avg_transaction_amount",
    "count"(DISTINCT "user_id") AS "unique_customers",
    "count"(DISTINCT "course_id") AS "courses_purchased"
   FROM "public"."transactions" "t"
  WHERE ("course_id" IS NOT NULL)
  GROUP BY ("date_trunc"('day'::"text", "created_at"))
  ORDER BY ("date_trunc"('day'::"text", "created_at")) DESC;


--
-- Name: payment_events; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE "public"."payment_events" (
    "id" "uuid" DEFAULT "gen_random_uuid"() NOT NULL,
    "event_type" character varying(50) NOT NULL,
    "provider" character varying(20) DEFAULT 'lemonsqueezy'::character varying NOT NULL,
    "provider_event_id" character varying(100) NOT NULL,
    "transaction_id" "uuid",
    "user_id" "uuid",
    "course_id" "uuid",
    "payload" "jsonb" NOT NULL,
    "processed" boolean DEFAULT false,
    "processed_at" timestamp without time zone,
    "error_message" "text",
    "retry_count" integer DEFAULT 0,
    "created_at" timestamp without time zone DEFAULT "now"() NOT NULL
);


--
-- Name: TABLE "payment_events"; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TABLE "public"."payment_events" IS 'Webhook events and payment provider communications';


--
-- Name: payment_methods; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE "public"."payment_methods" (
    "id" "uuid" DEFAULT "gen_random_uuid"() NOT NULL,
    "user_id" "uuid",
    "provider" character varying(50) NOT NULL,
    "token" character varying(255) NOT NULL,
    "card_last_four" character varying(4),
    "card_expiry" character varying(7),
    "is_default" boolean DEFAULT false,
    "created_at" timestamp without time zone DEFAULT "now"() NOT NULL,
    "updated_at" timestamp without time zone DEFAULT "now"() NOT NULL,
    "stripe_customer_id" character varying(100),
    "stripe_payment_method_id" character varying(255),
    "card_brand" character varying(20),
    "card_country" character varying(2),
    "card_funding" character varying(10),
    "billing_name" character varying(255),
    CONSTRAINT "chk_payment_methods_stripe_customer_id" CHECK (((("stripe_customer_id")::"text" ~ '^cus_[a-zA-Z0-9]+$'::"text") OR ("stripe_customer_id" IS NULL)))
);


--
-- Name: COLUMN "payment_methods"."stripe_payment_method_id"; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN "public"."payment_methods"."stripe_payment_method_id" IS 'Stripe PaymentMethod ID for saved payment methods';


--
-- Name: progress; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE "public"."progress" (
    "user_id" "uuid" NOT NULL,
    "lecture_id" "uuid" NOT NULL,
    "watched_duration" integer DEFAULT 0,
    "completed" boolean DEFAULT false,
    "last_watched_at" timestamp without time zone,
    "is_completed" boolean DEFAULT false,
    "last_accessed" timestamp without time zone DEFAULT "now"(),
    "completed_at" timestamp without time zone,
    "created_at" timestamp without time zone DEFAULT "now"(),
    "updated_at" timestamp without time zone DEFAULT "now"(),
    "id" "uuid" DEFAULT "gen_random_uuid"(),
    "course_id" "uuid" NOT NULL,
    "progress_percentage" numeric(5,2) DEFAULT 0.0,
    "watch_time_seconds" integer DEFAULT 0
);


--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE "public"."schema_migrations" (
    "version" bigint NOT NULL,
    "dirty" boolean NOT NULL
);


--
-- Name: stripe_customers; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE "public"."stripe_customers" (
    "id" "uuid" DEFAULT "gen_random_uuid"() NOT NULL,
    "user_id" "uuid",
    "stripe_customer_id" character varying(100) NOT NULL,
    "email" character varying(255),
    "name" character varying(255),
    "phone" character varying(50),
    "created_at" timestamp without time zone DEFAULT "now"(),
    "updated_at" timestamp without time zone DEFAULT "now"(),
    CONSTRAINT "chk_stripe_customers_stripe_customer_id" CHECK ((("stripe_customer_id")::"text" ~ '^cus_[a-zA-Z0-9]+$'::"text"))
);


--
-- Name: TABLE "stripe_customers"; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TABLE "public"."stripe_customers" IS 'Maps platform users to Stripe customers';


--
-- Name: stripe_products; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE "public"."stripe_products" (
    "id" "uuid" DEFAULT "gen_random_uuid"() NOT NULL,
    "course_id" "uuid",
    "stripe_product_id" character varying(100) NOT NULL,
    "stripe_price_id" character varying(100) NOT NULL,
    "product_name" character varying(255) NOT NULL,
    "product_description" "text",
    "price_amount" integer NOT NULL,
    "price_currency" character varying(3) DEFAULT 'USD'::character varying NOT NULL,
    "price_type" character varying(20) DEFAULT 'one_time'::character varying NOT NULL,
    "recurring_interval" character varying(20),
    "recurring_interval_count" integer DEFAULT 1,
    "active" boolean DEFAULT true,
    "created_at" timestamp without time zone DEFAULT "now"(),
    "updated_at" timestamp without time zone DEFAULT "now"(),
    CONSTRAINT "chk_stripe_products_stripe_price_id" CHECK ((("stripe_price_id")::"text" ~ '^price_[a-zA-Z0-9]+$'::"text")),
    CONSTRAINT "chk_stripe_products_stripe_product_id" CHECK ((("stripe_product_id")::"text" ~ '^prod_[a-zA-Z0-9]+$'::"text"))
);


--
-- Name: TABLE "stripe_products"; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TABLE "public"."stripe_products" IS 'Stores Stripe product and price information for courses';


--
-- Name: stripe_webhook_events; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE "public"."stripe_webhook_events" (
    "id" "uuid" DEFAULT "gen_random_uuid"() NOT NULL,
    "stripe_event_id" character varying(100) NOT NULL,
    "event_type" character varying(100) NOT NULL,
    "processed" boolean DEFAULT false,
    "processing_attempts" integer DEFAULT 0,
    "event_data" "jsonb" NOT NULL,
    "created_at" timestamp without time zone DEFAULT "now"(),
    "processed_at" timestamp without time zone,
    "error_message" "text"
);


--
-- Name: TABLE "stripe_webhook_events"; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TABLE "public"."stripe_webhook_events" IS 'Stripe webhook events storage for processing and deduplication';


--
-- Name: COLUMN "stripe_webhook_events"."stripe_event_id"; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN "public"."stripe_webhook_events"."stripe_event_id" IS 'Unique event ID from Stripe';


--
-- Name: COLUMN "stripe_webhook_events"."event_type"; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN "public"."stripe_webhook_events"."event_type" IS 'Type of Stripe event (e.g., payment_intent.succeeded)';


--
-- Name: COLUMN "stripe_webhook_events"."processed"; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN "public"."stripe_webhook_events"."processed" IS 'Whether the event has been successfully processed';


--
-- Name: COLUMN "stripe_webhook_events"."processing_attempts"; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN "public"."stripe_webhook_events"."processing_attempts" IS 'Number of times processing was attempted';


--
-- Name: COLUMN "stripe_webhook_events"."event_data"; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN "public"."stripe_webhook_events"."event_data" IS 'Full JSON payload from Stripe webhook';


--
-- Name: COLUMN "stripe_webhook_events"."error_message"; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN "public"."stripe_webhook_events"."error_message" IS 'Error message from last failed processing attempt';


--
-- Name: subscriptions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE "public"."subscriptions" (
    "id" "uuid" DEFAULT "gen_random_uuid"() NOT NULL,
    "user_id" "uuid",
    "payment_method_id" "uuid",
    "plan_name" character varying(50) NOT NULL,
    "status" character varying(20) NOT NULL,
    "billing_period" character varying(20) NOT NULL,
    "next_billing_date" timestamp without time zone,
    "price" numeric(10,2) NOT NULL,
    "created_at" timestamp without time zone DEFAULT "now"() NOT NULL,
    "updated_at" timestamp without time zone DEFAULT "now"() NOT NULL
);


--
-- Name: upload_sessions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE "public"."upload_sessions" (
    "id" "uuid" DEFAULT "gen_random_uuid"() NOT NULL,
    "upload_id" character varying(255) NOT NULL,
    "filename" character varying(255) NOT NULL,
    "content_type" character varying(100) NOT NULL,
    "total_size" bigint NOT NULL,
    "uploaded_size" bigint DEFAULT 0,
    "bucket_name" character varying(100) NOT NULL,
    "object_key" character varying(500) NOT NULL,
    "user_id" "uuid" NOT NULL,
    "status" character varying(20) DEFAULT 'active'::character varying,
    "created_at" timestamp without time zone DEFAULT "now"(),
    "expires_at" timestamp without time zone NOT NULL,
    CONSTRAINT "upload_sessions_status_check" CHECK ((("status")::"text" = ANY (ARRAY[('active'::character varying)::"text", ('completed'::character varying)::"text", ('aborted'::character varying)::"text"])))
);


--
-- Name: user_payment_methods; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE "public"."user_payment_methods" (
    "id" "uuid" DEFAULT "gen_random_uuid"() NOT NULL,
    "user_id" "uuid" NOT NULL,
    "provider" character varying(20) DEFAULT 'lemonsqueezy'::character varying NOT NULL,
    "provider_customer_id" character varying(100),
    "payment_method_type" character varying(20) DEFAULT 'card'::character varying NOT NULL,
    "last_four_digits" character varying(4),
    "expiry_month" integer,
    "expiry_year" integer,
    "brand" character varying(20),
    "is_default" boolean DEFAULT false,
    "is_active" boolean DEFAULT true,
    "created_at" timestamp without time zone DEFAULT "now"() NOT NULL,
    "updated_at" timestamp without time zone DEFAULT "now"() NOT NULL
);


--
-- Name: TABLE "user_payment_methods"; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TABLE "public"."user_payment_methods" IS 'Stored payment methods for users';


--
-- Name: users; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE "public"."users" (
    "id" "uuid" DEFAULT "gen_random_uuid"() NOT NULL,
    "username" character varying(50) NOT NULL,
    "email" character varying(100) NOT NULL,
    "password_hash" character varying(100),
    "role" "public"."user_role" DEFAULT 'student'::"public"."user_role" NOT NULL,
    "created_at" timestamp without time zone DEFAULT "now"() NOT NULL,
    "updated_at" timestamp without time zone DEFAULT "now"() NOT NULL,
    "provider" character varying(50) DEFAULT 'local'::character varying,
    "provider_id" character varying(100),
    "avatar_url" "text",
    "is_email_verified" boolean DEFAULT false,
    CONSTRAINT "check_password_for_local_users" CHECK ((((("provider")::"text" = 'local'::"text") AND ("password_hash" IS NOT NULL)) OR ((("provider")::"text" <> 'local'::"text") AND ("password_hash" IS NULL)) OR ((("provider")::"text" <> 'local'::"text") AND ("password_hash" IS NOT NULL))))
);


--
-- Name: video_analytics; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE "public"."video_analytics" (
    "id" "uuid" DEFAULT "gen_random_uuid"() NOT NULL,
    "video_id" "uuid",
    "date" "date" NOT NULL,
    "total_views" integer DEFAULT 0,
    "unique_viewers" integer DEFAULT 0,
    "total_watch_time_seconds" bigint DEFAULT 0,
    "avg_watch_time_seconds" integer DEFAULT 0,
    "completion_rate" numeric(5,2) DEFAULT 0,
    "quality_distribution" "jsonb",
    "geographic_distribution" "jsonb",
    "device_distribution" "jsonb",
    "created_at" timestamp without time zone DEFAULT "now"(),
    "updated_at" timestamp without time zone DEFAULT "now"()
);


--
-- Name: video_permissions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE "public"."video_permissions" (
    "id" "uuid" DEFAULT "gen_random_uuid"() NOT NULL,
    "video_id" "uuid",
    "user_id" "uuid",
    "role_id" "uuid",
    "permission_type" character varying(20) NOT NULL,
    "granted_by" "uuid" NOT NULL,
    "expires_at" timestamp without time zone,
    "created_at" timestamp without time zone DEFAULT "now"()
);


--
-- Name: video_qualities; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE "public"."video_qualities" (
    "id" "uuid" DEFAULT "gen_random_uuid"() NOT NULL,
    "video_id" "uuid",
    "quality_label" character varying(10) NOT NULL,
    "bitrate_kbps" integer NOT NULL,
    "width" integer NOT NULL,
    "height" integer NOT NULL,
    "fps" integer DEFAULT 30,
    "codec" character varying(20) DEFAULT 'h264'::character varying,
    "url" character varying(500) NOT NULL,
    "file_size_bytes" bigint,
    "created_at" timestamp without time zone DEFAULT "now"()
);


--
-- Name: videos; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE "public"."videos" (
    "id" "uuid" DEFAULT "gen_random_uuid"() NOT NULL,
    "cloudflare_uid" character varying(255) NOT NULL,
    "title" character varying(255) NOT NULL,
    "description" "text",
    "duration_seconds" integer,
    "file_size_bytes" bigint,
    "upload_user_id" "uuid" NOT NULL,
    "course_id" "uuid",
    "lecture_id" "uuid",
    "status" character varying(20) DEFAULT 'processing'::character varying,
    "visibility" character varying(20) DEFAULT 'private'::character varying,
    "thumbnail_url" character varying(500),
    "stream_url" character varying(500),
    "preview_url" character varying(500),
    "metadata" "jsonb",
    "created_at" timestamp without time zone DEFAULT "now"(),
    "updated_at" timestamp without time zone DEFAULT "now"(),
    "deleted_at" timestamp without time zone
);


--
-- Name: webhook_events; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE "public"."webhook_events" (
    "id" "uuid" DEFAULT "gen_random_uuid"() NOT NULL,
    "event_id" character varying(255) NOT NULL,
    "event_type" character varying(100) NOT NULL,
    "object_type" character varying(100) NOT NULL,
    "object_id" character varying(255) NOT NULL,
    "payload" "jsonb" NOT NULL,
    "processed" boolean DEFAULT false,
    "created_at" timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


--
-- Name: messages; Type: TABLE; Schema: realtime; Owner: -
--

CREATE TABLE "realtime"."messages" (
    "topic" "text" NOT NULL,
    "extension" "text" NOT NULL,
    "payload" "jsonb",
    "event" "text",
    "private" boolean DEFAULT false,
    "updated_at" timestamp without time zone DEFAULT "now"() NOT NULL,
    "inserted_at" timestamp without time zone DEFAULT "now"() NOT NULL,
    "id" "uuid" DEFAULT "gen_random_uuid"() NOT NULL
)
PARTITION BY RANGE ("inserted_at");


--
-- Name: messages_2025_10_04; Type: TABLE; Schema: realtime; Owner: -
--

CREATE TABLE "realtime"."messages_2025_10_04" (
    "topic" "text" NOT NULL,
    "extension" "text" NOT NULL,
    "payload" "jsonb",
    "event" "text",
    "private" boolean DEFAULT false,
    "updated_at" timestamp without time zone DEFAULT "now"() NOT NULL,
    "inserted_at" timestamp without time zone DEFAULT "now"() NOT NULL,
    "id" "uuid" DEFAULT "gen_random_uuid"() NOT NULL
);


--
-- Name: messages_2025_10_05; Type: TABLE; Schema: realtime; Owner: -
--

CREATE TABLE "realtime"."messages_2025_10_05" (
    "topic" "text" NOT NULL,
    "extension" "text" NOT NULL,
    "payload" "jsonb",
    "event" "text",
    "private" boolean DEFAULT false,
    "updated_at" timestamp without time zone DEFAULT "now"() NOT NULL,
    "inserted_at" timestamp without time zone DEFAULT "now"() NOT NULL,
    "id" "uuid" DEFAULT "gen_random_uuid"() NOT NULL
);


--
-- Name: messages_2025_10_06; Type: TABLE; Schema: realtime; Owner: -
--

CREATE TABLE "realtime"."messages_2025_10_06" (
    "topic" "text" NOT NULL,
    "extension" "text" NOT NULL,
    "payload" "jsonb",
    "event" "text",
    "private" boolean DEFAULT false,
    "updated_at" timestamp without time zone DEFAULT "now"() NOT NULL,
    "inserted_at" timestamp without time zone DEFAULT "now"() NOT NULL,
    "id" "uuid" DEFAULT "gen_random_uuid"() NOT NULL
);


--
-- Name: messages_2025_10_07; Type: TABLE; Schema: realtime; Owner: -
--

CREATE TABLE "realtime"."messages_2025_10_07" (
    "topic" "text" NOT NULL,
    "extension" "text" NOT NULL,
    "payload" "jsonb",
    "event" "text",
    "private" boolean DEFAULT false,
    "updated_at" timestamp without time zone DEFAULT "now"() NOT NULL,
    "inserted_at" timestamp without time zone DEFAULT "now"() NOT NULL,
    "id" "uuid" DEFAULT "gen_random_uuid"() NOT NULL
);


--
-- Name: messages_2025_10_08; Type: TABLE; Schema: realtime; Owner: -
--

CREATE TABLE "realtime"."messages_2025_10_08" (
    "topic" "text" NOT NULL,
    "extension" "text" NOT NULL,
    "payload" "jsonb",
    "event" "text",
    "private" boolean DEFAULT false,
    "updated_at" timestamp without time zone DEFAULT "now"() NOT NULL,
    "inserted_at" timestamp without time zone DEFAULT "now"() NOT NULL,
    "id" "uuid" DEFAULT "gen_random_uuid"() NOT NULL
);


--
-- Name: messages_2025_10_09; Type: TABLE; Schema: realtime; Owner: -
--

CREATE TABLE "realtime"."messages_2025_10_09" (
    "topic" "text" NOT NULL,
    "extension" "text" NOT NULL,
    "payload" "jsonb",
    "event" "text",
    "private" boolean DEFAULT false,
    "updated_at" timestamp without time zone DEFAULT "now"() NOT NULL,
    "inserted_at" timestamp without time zone DEFAULT "now"() NOT NULL,
    "id" "uuid" DEFAULT "gen_random_uuid"() NOT NULL
);


--
-- Name: messages_2025_10_10; Type: TABLE; Schema: realtime; Owner: -
--

CREATE TABLE "realtime"."messages_2025_10_10" (
    "topic" "text" NOT NULL,
    "extension" "text" NOT NULL,
    "payload" "jsonb",
    "event" "text",
    "private" boolean DEFAULT false,
    "updated_at" timestamp without time zone DEFAULT "now"() NOT NULL,
    "inserted_at" timestamp without time zone DEFAULT "now"() NOT NULL,
    "id" "uuid" DEFAULT "gen_random_uuid"() NOT NULL
);


--
-- Name: schema_migrations; Type: TABLE; Schema: realtime; Owner: -
--

CREATE TABLE "realtime"."schema_migrations" (
    "version" bigint NOT NULL,
    "inserted_at" timestamp(0) without time zone
);


--
-- Name: subscription; Type: TABLE; Schema: realtime; Owner: -
--

CREATE TABLE "realtime"."subscription" (
    "id" bigint NOT NULL,
    "subscription_id" "uuid" NOT NULL,
    "entity" "regclass" NOT NULL,
    "filters" "realtime"."user_defined_filter"[] DEFAULT '{}'::"realtime"."user_defined_filter"[] NOT NULL,
    "claims" "jsonb" NOT NULL,
    "claims_role" "regrole" GENERATED ALWAYS AS ("realtime"."to_regrole"(("claims" ->> 'role'::"text"))) STORED NOT NULL,
    "created_at" timestamp without time zone DEFAULT "timezone"('utc'::"text", "now"()) NOT NULL
);


--
-- Name: subscription_id_seq; Type: SEQUENCE; Schema: realtime; Owner: -
--

ALTER TABLE "realtime"."subscription" ALTER COLUMN "id" ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME "realtime"."subscription_id_seq"
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: buckets; Type: TABLE; Schema: storage; Owner: -
--

CREATE TABLE "storage"."buckets" (
    "id" "text" NOT NULL,
    "name" "text" NOT NULL,
    "owner" "uuid",
    "created_at" timestamp with time zone DEFAULT "now"(),
    "updated_at" timestamp with time zone DEFAULT "now"(),
    "public" boolean DEFAULT false,
    "avif_autodetection" boolean DEFAULT false,
    "file_size_limit" bigint,
    "allowed_mime_types" "text"[],
    "owner_id" "text",
    "type" "storage"."buckettype" DEFAULT 'STANDARD'::"storage"."buckettype" NOT NULL
);


--
-- Name: COLUMN "buckets"."owner"; Type: COMMENT; Schema: storage; Owner: -
--

COMMENT ON COLUMN "storage"."buckets"."owner" IS 'Field is deprecated, use owner_id instead';


--
-- Name: buckets_analytics; Type: TABLE; Schema: storage; Owner: -
--

CREATE TABLE "storage"."buckets_analytics" (
    "id" "text" NOT NULL,
    "type" "storage"."buckettype" DEFAULT 'ANALYTICS'::"storage"."buckettype" NOT NULL,
    "format" "text" DEFAULT 'ICEBERG'::"text" NOT NULL,
    "created_at" timestamp with time zone DEFAULT "now"() NOT NULL,
    "updated_at" timestamp with time zone DEFAULT "now"() NOT NULL
);


--
-- Name: iceberg_namespaces; Type: TABLE; Schema: storage; Owner: -
--

CREATE TABLE "storage"."iceberg_namespaces" (
    "id" "uuid" DEFAULT "gen_random_uuid"() NOT NULL,
    "bucket_id" "text" NOT NULL,
    "name" "text" NOT NULL COLLATE "pg_catalog"."C",
    "created_at" timestamp with time zone DEFAULT "now"() NOT NULL,
    "updated_at" timestamp with time zone DEFAULT "now"() NOT NULL
);


--
-- Name: iceberg_tables; Type: TABLE; Schema: storage; Owner: -
--

CREATE TABLE "storage"."iceberg_tables" (
    "id" "uuid" DEFAULT "gen_random_uuid"() NOT NULL,
    "namespace_id" "uuid" NOT NULL,
    "bucket_id" "text" NOT NULL,
    "name" "text" NOT NULL COLLATE "pg_catalog"."C",
    "location" "text" NOT NULL,
    "created_at" timestamp with time zone DEFAULT "now"() NOT NULL,
    "updated_at" timestamp with time zone DEFAULT "now"() NOT NULL
);


--
-- Name: migrations; Type: TABLE; Schema: storage; Owner: -
--

CREATE TABLE "storage"."migrations" (
    "id" integer NOT NULL,
    "name" character varying(100) NOT NULL,
    "hash" character varying(40) NOT NULL,
    "executed_at" timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


--
-- Name: objects; Type: TABLE; Schema: storage; Owner: -
--

CREATE TABLE "storage"."objects" (
    "id" "uuid" DEFAULT "gen_random_uuid"() NOT NULL,
    "bucket_id" "text",
    "name" "text",
    "owner" "uuid",
    "created_at" timestamp with time zone DEFAULT "now"(),
    "updated_at" timestamp with time zone DEFAULT "now"(),
    "last_accessed_at" timestamp with time zone DEFAULT "now"(),
    "metadata" "jsonb",
    "path_tokens" "text"[] GENERATED ALWAYS AS ("string_to_array"("name", '/'::"text")) STORED,
    "version" "text",
    "owner_id" "text",
    "user_metadata" "jsonb",
    "level" integer
);


--
-- Name: COLUMN "objects"."owner"; Type: COMMENT; Schema: storage; Owner: -
--

COMMENT ON COLUMN "storage"."objects"."owner" IS 'Field is deprecated, use owner_id instead';


--
-- Name: prefixes; Type: TABLE; Schema: storage; Owner: -
--

CREATE TABLE "storage"."prefixes" (
    "bucket_id" "text" NOT NULL,
    "name" "text" NOT NULL COLLATE "pg_catalog"."C",
    "level" integer GENERATED ALWAYS AS ("storage"."get_level"("name")) STORED NOT NULL,
    "created_at" timestamp with time zone DEFAULT "now"(),
    "updated_at" timestamp with time zone DEFAULT "now"()
);


--
-- Name: s3_multipart_uploads; Type: TABLE; Schema: storage; Owner: -
--

CREATE TABLE "storage"."s3_multipart_uploads" (
    "id" "text" NOT NULL,
    "in_progress_size" bigint DEFAULT 0 NOT NULL,
    "upload_signature" "text" NOT NULL,
    "bucket_id" "text" NOT NULL,
    "key" "text" NOT NULL COLLATE "pg_catalog"."C",
    "version" "text" NOT NULL,
    "owner_id" "text",
    "created_at" timestamp with time zone DEFAULT "now"() NOT NULL,
    "user_metadata" "jsonb"
);


--
-- Name: s3_multipart_uploads_parts; Type: TABLE; Schema: storage; Owner: -
--

CREATE TABLE "storage"."s3_multipart_uploads_parts" (
    "id" "uuid" DEFAULT "gen_random_uuid"() NOT NULL,
    "upload_id" "text" NOT NULL,
    "size" bigint DEFAULT 0 NOT NULL,
    "part_number" integer NOT NULL,
    "bucket_id" "text" NOT NULL,
    "key" "text" NOT NULL COLLATE "pg_catalog"."C",
    "etag" "text" NOT NULL,
    "owner_id" "text",
    "version" "text" NOT NULL,
    "created_at" timestamp with time zone DEFAULT "now"() NOT NULL
);


--
-- Name: hooks; Type: TABLE; Schema: supabase_functions; Owner: -
--

CREATE TABLE "supabase_functions"."hooks" (
    "id" bigint NOT NULL,
    "hook_table_id" integer NOT NULL,
    "hook_name" "text" NOT NULL,
    "created_at" timestamp with time zone DEFAULT "now"() NOT NULL,
    "request_id" bigint
);


--
-- Name: TABLE "hooks"; Type: COMMENT; Schema: supabase_functions; Owner: -
--

COMMENT ON TABLE "supabase_functions"."hooks" IS 'Supabase Functions Hooks: Audit trail for triggered hooks.';


--
-- Name: hooks_id_seq; Type: SEQUENCE; Schema: supabase_functions; Owner: -
--

CREATE SEQUENCE "supabase_functions"."hooks_id_seq"
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: hooks_id_seq; Type: SEQUENCE OWNED BY; Schema: supabase_functions; Owner: -
--

ALTER SEQUENCE "supabase_functions"."hooks_id_seq" OWNED BY "supabase_functions"."hooks"."id";


--
-- Name: migrations; Type: TABLE; Schema: supabase_functions; Owner: -
--

CREATE TABLE "supabase_functions"."migrations" (
    "version" "text" NOT NULL,
    "inserted_at" timestamp with time zone DEFAULT "now"() NOT NULL
);


--
-- Name: messages_2025_10_04; Type: TABLE ATTACH; Schema: realtime; Owner: -
--

ALTER TABLE ONLY "realtime"."messages" ATTACH PARTITION "realtime"."messages_2025_10_04" FOR VALUES FROM ('2025-10-04 00:00:00') TO ('2025-10-05 00:00:00');


--
-- Name: messages_2025_10_05; Type: TABLE ATTACH; Schema: realtime; Owner: -
--

ALTER TABLE ONLY "realtime"."messages" ATTACH PARTITION "realtime"."messages_2025_10_05" FOR VALUES FROM ('2025-10-05 00:00:00') TO ('2025-10-06 00:00:00');


--
-- Name: messages_2025_10_06; Type: TABLE ATTACH; Schema: realtime; Owner: -
--

ALTER TABLE ONLY "realtime"."messages" ATTACH PARTITION "realtime"."messages_2025_10_06" FOR VALUES FROM ('2025-10-06 00:00:00') TO ('2025-10-07 00:00:00');


--
-- Name: messages_2025_10_07; Type: TABLE ATTACH; Schema: realtime; Owner: -
--

ALTER TABLE ONLY "realtime"."messages" ATTACH PARTITION "realtime"."messages_2025_10_07" FOR VALUES FROM ('2025-10-07 00:00:00') TO ('2025-10-08 00:00:00');


--
-- Name: messages_2025_10_08; Type: TABLE ATTACH; Schema: realtime; Owner: -
--

ALTER TABLE ONLY "realtime"."messages" ATTACH PARTITION "realtime"."messages_2025_10_08" FOR VALUES FROM ('2025-10-08 00:00:00') TO ('2025-10-09 00:00:00');


--
-- Name: messages_2025_10_09; Type: TABLE ATTACH; Schema: realtime; Owner: -
--

ALTER TABLE ONLY "realtime"."messages" ATTACH PARTITION "realtime"."messages_2025_10_09" FOR VALUES FROM ('2025-10-09 00:00:00') TO ('2025-10-10 00:00:00');


--
-- Name: messages_2025_10_10; Type: TABLE ATTACH; Schema: realtime; Owner: -
--

ALTER TABLE ONLY "realtime"."messages" ATTACH PARTITION "realtime"."messages_2025_10_10" FOR VALUES FROM ('2025-10-10 00:00:00') TO ('2025-10-11 00:00:00');


--
-- Name: refresh_tokens id; Type: DEFAULT; Schema: auth; Owner: -
--

ALTER TABLE ONLY "auth"."refresh_tokens" ALTER COLUMN "id" SET DEFAULT "nextval"('"auth"."refresh_tokens_id_seq"'::"regclass");


--
-- Name: hooks id; Type: DEFAULT; Schema: supabase_functions; Owner: -
--

ALTER TABLE ONLY "supabase_functions"."hooks" ALTER COLUMN "id" SET DEFAULT "nextval"('"supabase_functions"."hooks_id_seq"'::"regclass");


--
-- Name: extensions extensions_pkey; Type: CONSTRAINT; Schema: _realtime; Owner: -
--

ALTER TABLE ONLY "_realtime"."extensions"
    ADD CONSTRAINT "extensions_pkey" PRIMARY KEY ("id");


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: _realtime; Owner: -
--

ALTER TABLE ONLY "_realtime"."schema_migrations"
    ADD CONSTRAINT "schema_migrations_pkey" PRIMARY KEY ("version");


--
-- Name: tenants tenants_pkey; Type: CONSTRAINT; Schema: _realtime; Owner: -
--

ALTER TABLE ONLY "_realtime"."tenants"
    ADD CONSTRAINT "tenants_pkey" PRIMARY KEY ("id");


--
-- Name: mfa_amr_claims amr_id_pk; Type: CONSTRAINT; Schema: auth; Owner: -
--

ALTER TABLE ONLY "auth"."mfa_amr_claims"
    ADD CONSTRAINT "amr_id_pk" PRIMARY KEY ("id");


--
-- Name: audit_log_entries audit_log_entries_pkey; Type: CONSTRAINT; Schema: auth; Owner: -
--

ALTER TABLE ONLY "auth"."audit_log_entries"
    ADD CONSTRAINT "audit_log_entries_pkey" PRIMARY KEY ("id");


--
-- Name: flow_state flow_state_pkey; Type: CONSTRAINT; Schema: auth; Owner: -
--

ALTER TABLE ONLY "auth"."flow_state"
    ADD CONSTRAINT "flow_state_pkey" PRIMARY KEY ("id");


--
-- Name: identities identities_pkey; Type: CONSTRAINT; Schema: auth; Owner: -
--

ALTER TABLE ONLY "auth"."identities"
    ADD CONSTRAINT "identities_pkey" PRIMARY KEY ("id");


--
-- Name: identities identities_provider_id_provider_unique; Type: CONSTRAINT; Schema: auth; Owner: -
--

ALTER TABLE ONLY "auth"."identities"
    ADD CONSTRAINT "identities_provider_id_provider_unique" UNIQUE ("provider_id", "provider");


--
-- Name: instances instances_pkey; Type: CONSTRAINT; Schema: auth; Owner: -
--

ALTER TABLE ONLY "auth"."instances"
    ADD CONSTRAINT "instances_pkey" PRIMARY KEY ("id");


--
-- Name: mfa_amr_claims mfa_amr_claims_session_id_authentication_method_pkey; Type: CONSTRAINT; Schema: auth; Owner: -
--

ALTER TABLE ONLY "auth"."mfa_amr_claims"
    ADD CONSTRAINT "mfa_amr_claims_session_id_authentication_method_pkey" UNIQUE ("session_id", "authentication_method");


--
-- Name: mfa_challenges mfa_challenges_pkey; Type: CONSTRAINT; Schema: auth; Owner: -
--

ALTER TABLE ONLY "auth"."mfa_challenges"
    ADD CONSTRAINT "mfa_challenges_pkey" PRIMARY KEY ("id");


--
-- Name: mfa_factors mfa_factors_last_challenged_at_key; Type: CONSTRAINT; Schema: auth; Owner: -
--

ALTER TABLE ONLY "auth"."mfa_factors"
    ADD CONSTRAINT "mfa_factors_last_challenged_at_key" UNIQUE ("last_challenged_at");


--
-- Name: mfa_factors mfa_factors_pkey; Type: CONSTRAINT; Schema: auth; Owner: -
--

ALTER TABLE ONLY "auth"."mfa_factors"
    ADD CONSTRAINT "mfa_factors_pkey" PRIMARY KEY ("id");


--
-- Name: oauth_authorizations oauth_authorizations_authorization_code_key; Type: CONSTRAINT; Schema: auth; Owner: -
--

ALTER TABLE ONLY "auth"."oauth_authorizations"
    ADD CONSTRAINT "oauth_authorizations_authorization_code_key" UNIQUE ("authorization_code");


--
-- Name: oauth_authorizations oauth_authorizations_authorization_id_key; Type: CONSTRAINT; Schema: auth; Owner: -
--

ALTER TABLE ONLY "auth"."oauth_authorizations"
    ADD CONSTRAINT "oauth_authorizations_authorization_id_key" UNIQUE ("authorization_id");


--
-- Name: oauth_authorizations oauth_authorizations_pkey; Type: CONSTRAINT; Schema: auth; Owner: -
--

ALTER TABLE ONLY "auth"."oauth_authorizations"
    ADD CONSTRAINT "oauth_authorizations_pkey" PRIMARY KEY ("id");


--
-- Name: oauth_clients oauth_clients_pkey; Type: CONSTRAINT; Schema: auth; Owner: -
--

ALTER TABLE ONLY "auth"."oauth_clients"
    ADD CONSTRAINT "oauth_clients_pkey" PRIMARY KEY ("id");


--
-- Name: oauth_consents oauth_consents_pkey; Type: CONSTRAINT; Schema: auth; Owner: -
--

ALTER TABLE ONLY "auth"."oauth_consents"
    ADD CONSTRAINT "oauth_consents_pkey" PRIMARY KEY ("id");


--
-- Name: oauth_consents oauth_consents_user_client_unique; Type: CONSTRAINT; Schema: auth; Owner: -
--

ALTER TABLE ONLY "auth"."oauth_consents"
    ADD CONSTRAINT "oauth_consents_user_client_unique" UNIQUE ("user_id", "client_id");


--
-- Name: one_time_tokens one_time_tokens_pkey; Type: CONSTRAINT; Schema: auth; Owner: -
--

ALTER TABLE ONLY "auth"."one_time_tokens"
    ADD CONSTRAINT "one_time_tokens_pkey" PRIMARY KEY ("id");


--
-- Name: refresh_tokens refresh_tokens_pkey; Type: CONSTRAINT; Schema: auth; Owner: -
--

ALTER TABLE ONLY "auth"."refresh_tokens"
    ADD CONSTRAINT "refresh_tokens_pkey" PRIMARY KEY ("id");


--
-- Name: refresh_tokens refresh_tokens_token_unique; Type: CONSTRAINT; Schema: auth; Owner: -
--

ALTER TABLE ONLY "auth"."refresh_tokens"
    ADD CONSTRAINT "refresh_tokens_token_unique" UNIQUE ("token");


--
-- Name: saml_providers saml_providers_entity_id_key; Type: CONSTRAINT; Schema: auth; Owner: -
--

ALTER TABLE ONLY "auth"."saml_providers"
    ADD CONSTRAINT "saml_providers_entity_id_key" UNIQUE ("entity_id");


--
-- Name: saml_providers saml_providers_pkey; Type: CONSTRAINT; Schema: auth; Owner: -
--

ALTER TABLE ONLY "auth"."saml_providers"
    ADD CONSTRAINT "saml_providers_pkey" PRIMARY KEY ("id");


--
-- Name: saml_relay_states saml_relay_states_pkey; Type: CONSTRAINT; Schema: auth; Owner: -
--

ALTER TABLE ONLY "auth"."saml_relay_states"
    ADD CONSTRAINT "saml_relay_states_pkey" PRIMARY KEY ("id");


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: auth; Owner: -
--

ALTER TABLE ONLY "auth"."schema_migrations"
    ADD CONSTRAINT "schema_migrations_pkey" PRIMARY KEY ("version");


--
-- Name: sessions sessions_pkey; Type: CONSTRAINT; Schema: auth; Owner: -
--

ALTER TABLE ONLY "auth"."sessions"
    ADD CONSTRAINT "sessions_pkey" PRIMARY KEY ("id");


--
-- Name: sso_domains sso_domains_pkey; Type: CONSTRAINT; Schema: auth; Owner: -
--

ALTER TABLE ONLY "auth"."sso_domains"
    ADD CONSTRAINT "sso_domains_pkey" PRIMARY KEY ("id");


--
-- Name: sso_providers sso_providers_pkey; Type: CONSTRAINT; Schema: auth; Owner: -
--

ALTER TABLE ONLY "auth"."sso_providers"
    ADD CONSTRAINT "sso_providers_pkey" PRIMARY KEY ("id");


--
-- Name: users users_phone_key; Type: CONSTRAINT; Schema: auth; Owner: -
--

ALTER TABLE ONLY "auth"."users"
    ADD CONSTRAINT "users_phone_key" UNIQUE ("phone");


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: auth; Owner: -
--

ALTER TABLE ONLY "auth"."users"
    ADD CONSTRAINT "users_pkey" PRIMARY KEY ("id");


--
-- Name: adaptive_streaming_rules adaptive_streaming_rules_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."adaptive_streaming_rules"
    ADD CONSTRAINT "adaptive_streaming_rules_pkey" PRIMARY KEY ("id");


--
-- Name: audit_logs audit_logs_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."audit_logs"
    ADD CONSTRAINT "audit_logs_pkey" PRIMARY KEY ("id");


--
-- Name: bandwidth_tests bandwidth_tests_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."bandwidth_tests"
    ADD CONSTRAINT "bandwidth_tests_pkey" PRIMARY KEY ("id");


--
-- Name: chat_history chat_history_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."chat_history"
    ADD CONSTRAINT "chat_history_pkey" PRIMARY KEY ("id");


--
-- Name: course_access_cache course_access_cache_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."course_access_cache"
    ADD CONSTRAINT "course_access_cache_pkey" PRIMARY KEY ("id");


--
-- Name: course_access_cache course_access_cache_user_id_course_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."course_access_cache"
    ADD CONSTRAINT "course_access_cache_user_id_course_id_key" UNIQUE ("user_id", "course_id");


--
-- Name: course_access_logs course_access_logs_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."course_access_logs"
    ADD CONSTRAINT "course_access_logs_pkey" PRIMARY KEY ("id");


--
-- Name: course_resources course_resources_course_id_file_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."course_resources"
    ADD CONSTRAINT "course_resources_course_id_file_id_key" UNIQUE ("course_id", "file_id");


--
-- Name: course_resources course_resources_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."course_resources"
    ADD CONSTRAINT "course_resources_pkey" PRIMARY KEY ("id");


--
-- Name: courses courses_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."courses"
    ADD CONSTRAINT "courses_pkey" PRIMARY KEY ("id");


--
-- Name: enrollments enrollments_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."enrollments"
    ADD CONSTRAINT "enrollments_pkey" PRIMARY KEY ("id");


--
-- Name: enrollments enrollments_user_id_course_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."enrollments"
    ADD CONSTRAINT "enrollments_user_id_course_id_key" UNIQUE ("user_id", "course_id");


--
-- Name: file_permissions file_permissions_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."file_permissions"
    ADD CONSTRAINT "file_permissions_pkey" PRIMARY KEY ("id");


--
-- Name: files files_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."files"
    ADD CONSTRAINT "files_pkey" PRIMARY KEY ("id");


--
-- Name: forum_mentions forum_mentions_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."forum_mentions"
    ADD CONSTRAINT "forum_mentions_pkey" PRIMARY KEY ("id");


--
-- Name: forum_mentions forum_mentions_post_id_mentioned_user_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."forum_mentions"
    ADD CONSTRAINT "forum_mentions_post_id_mentioned_user_id_key" UNIQUE ("post_id", "mentioned_user_id");


--
-- Name: forum_notifications forum_notifications_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."forum_notifications"
    ADD CONSTRAINT "forum_notifications_pkey" PRIMARY KEY ("id");


--
-- Name: forum_posts forum_posts_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."forum_posts"
    ADD CONSTRAINT "forum_posts_pkey" PRIMARY KEY ("id");


--
-- Name: forum_topic_subscriptions forum_topic_subscriptions_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."forum_topic_subscriptions"
    ADD CONSTRAINT "forum_topic_subscriptions_pkey" PRIMARY KEY ("id");


--
-- Name: forum_topic_subscriptions forum_topic_subscriptions_user_id_topic_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."forum_topic_subscriptions"
    ADD CONSTRAINT "forum_topic_subscriptions_user_id_topic_id_key" UNIQUE ("user_id", "topic_id");


--
-- Name: forum_topics forum_topics_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."forum_topics"
    ADD CONSTRAINT "forum_topics_pkey" PRIMARY KEY ("id");


--
-- Name: forum_votes forum_votes_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."forum_votes"
    ADD CONSTRAINT "forum_votes_pkey" PRIMARY KEY ("id");


--
-- Name: forum_votes forum_votes_post_id_user_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."forum_votes"
    ADD CONSTRAINT "forum_votes_post_id_user_id_key" UNIQUE ("post_id", "user_id");


--
-- Name: lecture_preview_sessions lecture_preview_sessions_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."lecture_preview_sessions"
    ADD CONSTRAINT "lecture_preview_sessions_pkey" PRIMARY KEY ("id");


--
-- Name: lecture_preview_sessions lecture_preview_sessions_user_id_lecture_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."lecture_preview_sessions"
    ADD CONSTRAINT "lecture_preview_sessions_user_id_lecture_id_key" UNIQUE ("user_id", "lecture_id");


--
-- Name: lecture_resources lecture_resources_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."lecture_resources"
    ADD CONSTRAINT "lecture_resources_pkey" PRIMARY KEY ("id");


--
-- Name: lectures lectures_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."lectures"
    ADD CONSTRAINT "lectures_pkey" PRIMARY KEY ("id");


--
-- Name: lemon_squeezy_products lemon_squeezy_products_lemon_squeezy_product_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."lemon_squeezy_products"
    ADD CONSTRAINT "lemon_squeezy_products_lemon_squeezy_product_id_key" UNIQUE ("lemon_squeezy_product_id");


--
-- Name: lemon_squeezy_products lemon_squeezy_products_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."lemon_squeezy_products"
    ADD CONSTRAINT "lemon_squeezy_products_pkey" PRIMARY KEY ("id");


--
-- Name: lemon_squeezy_variants lemon_squeezy_variants_lemon_squeezy_variant_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."lemon_squeezy_variants"
    ADD CONSTRAINT "lemon_squeezy_variants_lemon_squeezy_variant_id_key" UNIQUE ("lemon_squeezy_variant_id");


--
-- Name: lemon_squeezy_variants lemon_squeezy_variants_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."lemon_squeezy_variants"
    ADD CONSTRAINT "lemon_squeezy_variants_pkey" PRIMARY KEY ("id");


--
-- Name: lemon_squeezy_webhook_events lemon_squeezy_webhook_events_event_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."lemon_squeezy_webhook_events"
    ADD CONSTRAINT "lemon_squeezy_webhook_events_event_id_key" UNIQUE ("event_id");


--
-- Name: lemon_squeezy_webhook_events lemon_squeezy_webhook_events_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."lemon_squeezy_webhook_events"
    ADD CONSTRAINT "lemon_squeezy_webhook_events_pkey" PRIMARY KEY ("id");


--
-- Name: network_analytics_daily network_analytics_daily_date_user_id_video_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."network_analytics_daily"
    ADD CONSTRAINT "network_analytics_daily_date_user_id_video_id_key" UNIQUE ("date", "user_id", "video_id");


--
-- Name: network_analytics_daily network_analytics_daily_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."network_analytics_daily"
    ADD CONSTRAINT "network_analytics_daily_pkey" PRIMARY KEY ("id");


--
-- Name: network_events network_events_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."network_events"
    ADD CONSTRAINT "network_events_pkey" PRIMARY KEY ("id");


--
-- Name: network_metrics network_metrics_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."network_metrics"
    ADD CONSTRAINT "network_metrics_pkey" PRIMARY KEY ("id");


--
-- Name: notes notes_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."notes"
    ADD CONSTRAINT "notes_pkey" PRIMARY KEY ("id");


--
-- Name: oauth_accounts oauth_accounts_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."oauth_accounts"
    ADD CONSTRAINT "oauth_accounts_pkey" PRIMARY KEY ("id");


--
-- Name: oauth_accounts oauth_accounts_provider_provider_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."oauth_accounts"
    ADD CONSTRAINT "oauth_accounts_provider_provider_id_key" UNIQUE ("provider", "provider_id");


--
-- Name: payment_events payment_events_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."payment_events"
    ADD CONSTRAINT "payment_events_pkey" PRIMARY KEY ("id");


--
-- Name: payment_events payment_events_provider_provider_event_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."payment_events"
    ADD CONSTRAINT "payment_events_provider_provider_event_id_key" UNIQUE ("provider", "provider_event_id");


--
-- Name: payment_methods payment_methods_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."payment_methods"
    ADD CONSTRAINT "payment_methods_pkey" PRIMARY KEY ("id");


--
-- Name: progress progress_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."progress"
    ADD CONSTRAINT "progress_pkey" PRIMARY KEY ("user_id", "lecture_id");


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."schema_migrations"
    ADD CONSTRAINT "schema_migrations_pkey" PRIMARY KEY ("version");


--
-- Name: stripe_customers stripe_customers_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."stripe_customers"
    ADD CONSTRAINT "stripe_customers_pkey" PRIMARY KEY ("id");


--
-- Name: stripe_customers stripe_customers_stripe_customer_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."stripe_customers"
    ADD CONSTRAINT "stripe_customers_stripe_customer_id_key" UNIQUE ("stripe_customer_id");


--
-- Name: stripe_products stripe_products_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."stripe_products"
    ADD CONSTRAINT "stripe_products_pkey" PRIMARY KEY ("id");


--
-- Name: stripe_products stripe_products_stripe_price_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."stripe_products"
    ADD CONSTRAINT "stripe_products_stripe_price_id_key" UNIQUE ("stripe_price_id");


--
-- Name: stripe_products stripe_products_stripe_product_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."stripe_products"
    ADD CONSTRAINT "stripe_products_stripe_product_id_key" UNIQUE ("stripe_product_id");


--
-- Name: stripe_webhook_events stripe_webhook_events_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."stripe_webhook_events"
    ADD CONSTRAINT "stripe_webhook_events_pkey" PRIMARY KEY ("id");


--
-- Name: stripe_webhook_events stripe_webhook_events_stripe_event_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."stripe_webhook_events"
    ADD CONSTRAINT "stripe_webhook_events_stripe_event_id_key" UNIQUE ("stripe_event_id");


--
-- Name: subscriptions subscriptions_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."subscriptions"
    ADD CONSTRAINT "subscriptions_pkey" PRIMARY KEY ("id");


--
-- Name: transactions transactions_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."transactions"
    ADD CONSTRAINT "transactions_pkey" PRIMARY KEY ("id");


--
-- Name: transactions transactions_transaction_reference_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."transactions"
    ADD CONSTRAINT "transactions_transaction_reference_key" UNIQUE ("transaction_reference");


--
-- Name: lecture_resources uk_lecture_resources_lecture_file; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."lecture_resources"
    ADD CONSTRAINT "uk_lecture_resources_lecture_file" UNIQUE ("lecture_id", "file_id");


--
-- Name: transactions uk_transactions_stripe_payment_intent_id; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."transactions"
    ADD CONSTRAINT "uk_transactions_stripe_payment_intent_id" UNIQUE ("stripe_payment_intent_id") DEFERRABLE INITIALLY DEFERRED;


--
-- Name: upload_sessions upload_sessions_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."upload_sessions"
    ADD CONSTRAINT "upload_sessions_pkey" PRIMARY KEY ("id");


--
-- Name: user_payment_methods user_payment_methods_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."user_payment_methods"
    ADD CONSTRAINT "user_payment_methods_pkey" PRIMARY KEY ("id");


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."users"
    ADD CONSTRAINT "users_email_key" UNIQUE ("email");


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."users"
    ADD CONSTRAINT "users_pkey" PRIMARY KEY ("id");


--
-- Name: users users_username_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."users"
    ADD CONSTRAINT "users_username_key" UNIQUE ("username");


--
-- Name: video_analytics video_analytics_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."video_analytics"
    ADD CONSTRAINT "video_analytics_pkey" PRIMARY KEY ("id");


--
-- Name: video_analytics video_analytics_video_id_date_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."video_analytics"
    ADD CONSTRAINT "video_analytics_video_id_date_key" UNIQUE ("video_id", "date");


--
-- Name: video_permissions video_permissions_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."video_permissions"
    ADD CONSTRAINT "video_permissions_pkey" PRIMARY KEY ("id");


--
-- Name: video_qualities video_qualities_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."video_qualities"
    ADD CONSTRAINT "video_qualities_pkey" PRIMARY KEY ("id");


--
-- Name: videos videos_cloudflare_uid_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."videos"
    ADD CONSTRAINT "videos_cloudflare_uid_key" UNIQUE ("cloudflare_uid");


--
-- Name: videos videos_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."videos"
    ADD CONSTRAINT "videos_pkey" PRIMARY KEY ("id");


--
-- Name: viewing_sessions viewing_sessions_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."viewing_sessions"
    ADD CONSTRAINT "viewing_sessions_pkey" PRIMARY KEY ("id");


--
-- Name: webhook_events webhook_events_event_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."webhook_events"
    ADD CONSTRAINT "webhook_events_event_id_key" UNIQUE ("event_id");


--
-- Name: webhook_events webhook_events_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."webhook_events"
    ADD CONSTRAINT "webhook_events_pkey" PRIMARY KEY ("id");


--
-- Name: messages messages_pkey; Type: CONSTRAINT; Schema: realtime; Owner: -
--

ALTER TABLE ONLY "realtime"."messages"
    ADD CONSTRAINT "messages_pkey" PRIMARY KEY ("id", "inserted_at");


--
-- Name: messages_2025_10_04 messages_2025_10_04_pkey; Type: CONSTRAINT; Schema: realtime; Owner: -
--

ALTER TABLE ONLY "realtime"."messages_2025_10_04"
    ADD CONSTRAINT "messages_2025_10_04_pkey" PRIMARY KEY ("id", "inserted_at");


--
-- Name: messages_2025_10_05 messages_2025_10_05_pkey; Type: CONSTRAINT; Schema: realtime; Owner: -
--

ALTER TABLE ONLY "realtime"."messages_2025_10_05"
    ADD CONSTRAINT "messages_2025_10_05_pkey" PRIMARY KEY ("id", "inserted_at");


--
-- Name: messages_2025_10_06 messages_2025_10_06_pkey; Type: CONSTRAINT; Schema: realtime; Owner: -
--

ALTER TABLE ONLY "realtime"."messages_2025_10_06"
    ADD CONSTRAINT "messages_2025_10_06_pkey" PRIMARY KEY ("id", "inserted_at");


--
-- Name: messages_2025_10_07 messages_2025_10_07_pkey; Type: CONSTRAINT; Schema: realtime; Owner: -
--

ALTER TABLE ONLY "realtime"."messages_2025_10_07"
    ADD CONSTRAINT "messages_2025_10_07_pkey" PRIMARY KEY ("id", "inserted_at");


--
-- Name: messages_2025_10_08 messages_2025_10_08_pkey; Type: CONSTRAINT; Schema: realtime; Owner: -
--

ALTER TABLE ONLY "realtime"."messages_2025_10_08"
    ADD CONSTRAINT "messages_2025_10_08_pkey" PRIMARY KEY ("id", "inserted_at");


--
-- Name: messages_2025_10_09 messages_2025_10_09_pkey; Type: CONSTRAINT; Schema: realtime; Owner: -
--

ALTER TABLE ONLY "realtime"."messages_2025_10_09"
    ADD CONSTRAINT "messages_2025_10_09_pkey" PRIMARY KEY ("id", "inserted_at");


--
-- Name: messages_2025_10_10 messages_2025_10_10_pkey; Type: CONSTRAINT; Schema: realtime; Owner: -
--

ALTER TABLE ONLY "realtime"."messages_2025_10_10"
    ADD CONSTRAINT "messages_2025_10_10_pkey" PRIMARY KEY ("id", "inserted_at");


--
-- Name: subscription pk_subscription; Type: CONSTRAINT; Schema: realtime; Owner: -
--

ALTER TABLE ONLY "realtime"."subscription"
    ADD CONSTRAINT "pk_subscription" PRIMARY KEY ("id");


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: realtime; Owner: -
--

ALTER TABLE ONLY "realtime"."schema_migrations"
    ADD CONSTRAINT "schema_migrations_pkey" PRIMARY KEY ("version");


--
-- Name: buckets_analytics buckets_analytics_pkey; Type: CONSTRAINT; Schema: storage; Owner: -
--

ALTER TABLE ONLY "storage"."buckets_analytics"
    ADD CONSTRAINT "buckets_analytics_pkey" PRIMARY KEY ("id");


--
-- Name: buckets buckets_pkey; Type: CONSTRAINT; Schema: storage; Owner: -
--

ALTER TABLE ONLY "storage"."buckets"
    ADD CONSTRAINT "buckets_pkey" PRIMARY KEY ("id");


--
-- Name: iceberg_namespaces iceberg_namespaces_pkey; Type: CONSTRAINT; Schema: storage; Owner: -
--

ALTER TABLE ONLY "storage"."iceberg_namespaces"
    ADD CONSTRAINT "iceberg_namespaces_pkey" PRIMARY KEY ("id");


--
-- Name: iceberg_tables iceberg_tables_pkey; Type: CONSTRAINT; Schema: storage; Owner: -
--

ALTER TABLE ONLY "storage"."iceberg_tables"
    ADD CONSTRAINT "iceberg_tables_pkey" PRIMARY KEY ("id");


--
-- Name: migrations migrations_name_key; Type: CONSTRAINT; Schema: storage; Owner: -
--

ALTER TABLE ONLY "storage"."migrations"
    ADD CONSTRAINT "migrations_name_key" UNIQUE ("name");


--
-- Name: migrations migrations_pkey; Type: CONSTRAINT; Schema: storage; Owner: -
--

ALTER TABLE ONLY "storage"."migrations"
    ADD CONSTRAINT "migrations_pkey" PRIMARY KEY ("id");


--
-- Name: objects objects_pkey; Type: CONSTRAINT; Schema: storage; Owner: -
--

ALTER TABLE ONLY "storage"."objects"
    ADD CONSTRAINT "objects_pkey" PRIMARY KEY ("id");


--
-- Name: prefixes prefixes_pkey; Type: CONSTRAINT; Schema: storage; Owner: -
--

ALTER TABLE ONLY "storage"."prefixes"
    ADD CONSTRAINT "prefixes_pkey" PRIMARY KEY ("bucket_id", "level", "name");


--
-- Name: s3_multipart_uploads_parts s3_multipart_uploads_parts_pkey; Type: CONSTRAINT; Schema: storage; Owner: -
--

ALTER TABLE ONLY "storage"."s3_multipart_uploads_parts"
    ADD CONSTRAINT "s3_multipart_uploads_parts_pkey" PRIMARY KEY ("id");


--
-- Name: s3_multipart_uploads s3_multipart_uploads_pkey; Type: CONSTRAINT; Schema: storage; Owner: -
--

ALTER TABLE ONLY "storage"."s3_multipart_uploads"
    ADD CONSTRAINT "s3_multipart_uploads_pkey" PRIMARY KEY ("id");


--
-- Name: hooks hooks_pkey; Type: CONSTRAINT; Schema: supabase_functions; Owner: -
--

ALTER TABLE ONLY "supabase_functions"."hooks"
    ADD CONSTRAINT "hooks_pkey" PRIMARY KEY ("id");


--
-- Name: migrations migrations_pkey; Type: CONSTRAINT; Schema: supabase_functions; Owner: -
--

ALTER TABLE ONLY "supabase_functions"."migrations"
    ADD CONSTRAINT "migrations_pkey" PRIMARY KEY ("version");


--
-- Name: extensions_tenant_external_id_index; Type: INDEX; Schema: _realtime; Owner: -
--

CREATE INDEX "extensions_tenant_external_id_index" ON "_realtime"."extensions" USING "btree" ("tenant_external_id");


--
-- Name: extensions_tenant_external_id_type_index; Type: INDEX; Schema: _realtime; Owner: -
--

CREATE UNIQUE INDEX "extensions_tenant_external_id_type_index" ON "_realtime"."extensions" USING "btree" ("tenant_external_id", "type");


--
-- Name: tenants_external_id_index; Type: INDEX; Schema: _realtime; Owner: -
--

CREATE UNIQUE INDEX "tenants_external_id_index" ON "_realtime"."tenants" USING "btree" ("external_id");


--
-- Name: audit_logs_instance_id_idx; Type: INDEX; Schema: auth; Owner: -
--

CREATE INDEX "audit_logs_instance_id_idx" ON "auth"."audit_log_entries" USING "btree" ("instance_id");


--
-- Name: confirmation_token_idx; Type: INDEX; Schema: auth; Owner: -
--

CREATE UNIQUE INDEX "confirmation_token_idx" ON "auth"."users" USING "btree" ("confirmation_token") WHERE (("confirmation_token")::"text" !~ '^[0-9 ]*$'::"text");


--
-- Name: email_change_token_current_idx; Type: INDEX; Schema: auth; Owner: -
--

CREATE UNIQUE INDEX "email_change_token_current_idx" ON "auth"."users" USING "btree" ("email_change_token_current") WHERE (("email_change_token_current")::"text" !~ '^[0-9 ]*$'::"text");


--
-- Name: email_change_token_new_idx; Type: INDEX; Schema: auth; Owner: -
--

CREATE UNIQUE INDEX "email_change_token_new_idx" ON "auth"."users" USING "btree" ("email_change_token_new") WHERE (("email_change_token_new")::"text" !~ '^[0-9 ]*$'::"text");


--
-- Name: factor_id_created_at_idx; Type: INDEX; Schema: auth; Owner: -
--

CREATE INDEX "factor_id_created_at_idx" ON "auth"."mfa_factors" USING "btree" ("user_id", "created_at");


--
-- Name: flow_state_created_at_idx; Type: INDEX; Schema: auth; Owner: -
--

CREATE INDEX "flow_state_created_at_idx" ON "auth"."flow_state" USING "btree" ("created_at" DESC);


--
-- Name: identities_email_idx; Type: INDEX; Schema: auth; Owner: -
--

CREATE INDEX "identities_email_idx" ON "auth"."identities" USING "btree" ("email" "text_pattern_ops");


--
-- Name: INDEX "identities_email_idx"; Type: COMMENT; Schema: auth; Owner: -
--

COMMENT ON INDEX "auth"."identities_email_idx" IS 'Auth: Ensures indexed queries on the email column';


--
-- Name: identities_user_id_idx; Type: INDEX; Schema: auth; Owner: -
--

CREATE INDEX "identities_user_id_idx" ON "auth"."identities" USING "btree" ("user_id");


--
-- Name: idx_auth_code; Type: INDEX; Schema: auth; Owner: -
--

CREATE INDEX "idx_auth_code" ON "auth"."flow_state" USING "btree" ("auth_code");


--
-- Name: idx_user_id_auth_method; Type: INDEX; Schema: auth; Owner: -
--

CREATE INDEX "idx_user_id_auth_method" ON "auth"."flow_state" USING "btree" ("user_id", "authentication_method");


--
-- Name: mfa_challenge_created_at_idx; Type: INDEX; Schema: auth; Owner: -
--

CREATE INDEX "mfa_challenge_created_at_idx" ON "auth"."mfa_challenges" USING "btree" ("created_at" DESC);


--
-- Name: mfa_factors_user_friendly_name_unique; Type: INDEX; Schema: auth; Owner: -
--

CREATE UNIQUE INDEX "mfa_factors_user_friendly_name_unique" ON "auth"."mfa_factors" USING "btree" ("friendly_name", "user_id") WHERE (TRIM(BOTH FROM "friendly_name") <> ''::"text");


--
-- Name: mfa_factors_user_id_idx; Type: INDEX; Schema: auth; Owner: -
--

CREATE INDEX "mfa_factors_user_id_idx" ON "auth"."mfa_factors" USING "btree" ("user_id");


--
-- Name: oauth_auth_pending_exp_idx; Type: INDEX; Schema: auth; Owner: -
--

CREATE INDEX "oauth_auth_pending_exp_idx" ON "auth"."oauth_authorizations" USING "btree" ("expires_at") WHERE ("status" = 'pending'::"auth"."oauth_authorization_status");


--
-- Name: oauth_clients_deleted_at_idx; Type: INDEX; Schema: auth; Owner: -
--

CREATE INDEX "oauth_clients_deleted_at_idx" ON "auth"."oauth_clients" USING "btree" ("deleted_at");


--
-- Name: oauth_consents_active_client_idx; Type: INDEX; Schema: auth; Owner: -
--

CREATE INDEX "oauth_consents_active_client_idx" ON "auth"."oauth_consents" USING "btree" ("client_id") WHERE ("revoked_at" IS NULL);


--
-- Name: oauth_consents_active_user_client_idx; Type: INDEX; Schema: auth; Owner: -
--

CREATE INDEX "oauth_consents_active_user_client_idx" ON "auth"."oauth_consents" USING "btree" ("user_id", "client_id") WHERE ("revoked_at" IS NULL);


--
-- Name: oauth_consents_user_order_idx; Type: INDEX; Schema: auth; Owner: -
--

CREATE INDEX "oauth_consents_user_order_idx" ON "auth"."oauth_consents" USING "btree" ("user_id", "granted_at" DESC);


--
-- Name: one_time_tokens_relates_to_hash_idx; Type: INDEX; Schema: auth; Owner: -
--

CREATE INDEX "one_time_tokens_relates_to_hash_idx" ON "auth"."one_time_tokens" USING "hash" ("relates_to");


--
-- Name: one_time_tokens_token_hash_hash_idx; Type: INDEX; Schema: auth; Owner: -
--

CREATE INDEX "one_time_tokens_token_hash_hash_idx" ON "auth"."one_time_tokens" USING "hash" ("token_hash");


--
-- Name: one_time_tokens_user_id_token_type_key; Type: INDEX; Schema: auth; Owner: -
--

CREATE UNIQUE INDEX "one_time_tokens_user_id_token_type_key" ON "auth"."one_time_tokens" USING "btree" ("user_id", "token_type");


--
-- Name: reauthentication_token_idx; Type: INDEX; Schema: auth; Owner: -
--

CREATE UNIQUE INDEX "reauthentication_token_idx" ON "auth"."users" USING "btree" ("reauthentication_token") WHERE (("reauthentication_token")::"text" !~ '^[0-9 ]*$'::"text");


--
-- Name: recovery_token_idx; Type: INDEX; Schema: auth; Owner: -
--

CREATE UNIQUE INDEX "recovery_token_idx" ON "auth"."users" USING "btree" ("recovery_token") WHERE (("recovery_token")::"text" !~ '^[0-9 ]*$'::"text");


--
-- Name: refresh_tokens_instance_id_idx; Type: INDEX; Schema: auth; Owner: -
--

CREATE INDEX "refresh_tokens_instance_id_idx" ON "auth"."refresh_tokens" USING "btree" ("instance_id");


--
-- Name: refresh_tokens_instance_id_user_id_idx; Type: INDEX; Schema: auth; Owner: -
--

CREATE INDEX "refresh_tokens_instance_id_user_id_idx" ON "auth"."refresh_tokens" USING "btree" ("instance_id", "user_id");


--
-- Name: refresh_tokens_parent_idx; Type: INDEX; Schema: auth; Owner: -
--

CREATE INDEX "refresh_tokens_parent_idx" ON "auth"."refresh_tokens" USING "btree" ("parent");


--
-- Name: refresh_tokens_session_id_revoked_idx; Type: INDEX; Schema: auth; Owner: -
--

CREATE INDEX "refresh_tokens_session_id_revoked_idx" ON "auth"."refresh_tokens" USING "btree" ("session_id", "revoked");


--
-- Name: refresh_tokens_updated_at_idx; Type: INDEX; Schema: auth; Owner: -
--

CREATE INDEX "refresh_tokens_updated_at_idx" ON "auth"."refresh_tokens" USING "btree" ("updated_at" DESC);


--
-- Name: saml_providers_sso_provider_id_idx; Type: INDEX; Schema: auth; Owner: -
--

CREATE INDEX "saml_providers_sso_provider_id_idx" ON "auth"."saml_providers" USING "btree" ("sso_provider_id");


--
-- Name: saml_relay_states_created_at_idx; Type: INDEX; Schema: auth; Owner: -
--

CREATE INDEX "saml_relay_states_created_at_idx" ON "auth"."saml_relay_states" USING "btree" ("created_at" DESC);


--
-- Name: saml_relay_states_for_email_idx; Type: INDEX; Schema: auth; Owner: -
--

CREATE INDEX "saml_relay_states_for_email_idx" ON "auth"."saml_relay_states" USING "btree" ("for_email");


--
-- Name: saml_relay_states_sso_provider_id_idx; Type: INDEX; Schema: auth; Owner: -
--

CREATE INDEX "saml_relay_states_sso_provider_id_idx" ON "auth"."saml_relay_states" USING "btree" ("sso_provider_id");


--
-- Name: sessions_not_after_idx; Type: INDEX; Schema: auth; Owner: -
--

CREATE INDEX "sessions_not_after_idx" ON "auth"."sessions" USING "btree" ("not_after" DESC);


--
-- Name: sessions_oauth_client_id_idx; Type: INDEX; Schema: auth; Owner: -
--

CREATE INDEX "sessions_oauth_client_id_idx" ON "auth"."sessions" USING "btree" ("oauth_client_id");


--
-- Name: sessions_user_id_idx; Type: INDEX; Schema: auth; Owner: -
--

CREATE INDEX "sessions_user_id_idx" ON "auth"."sessions" USING "btree" ("user_id");


--
-- Name: sso_domains_domain_idx; Type: INDEX; Schema: auth; Owner: -
--

CREATE UNIQUE INDEX "sso_domains_domain_idx" ON "auth"."sso_domains" USING "btree" ("lower"("domain"));


--
-- Name: sso_domains_sso_provider_id_idx; Type: INDEX; Schema: auth; Owner: -
--

CREATE INDEX "sso_domains_sso_provider_id_idx" ON "auth"."sso_domains" USING "btree" ("sso_provider_id");


--
-- Name: sso_providers_resource_id_idx; Type: INDEX; Schema: auth; Owner: -
--

CREATE UNIQUE INDEX "sso_providers_resource_id_idx" ON "auth"."sso_providers" USING "btree" ("lower"("resource_id"));


--
-- Name: sso_providers_resource_id_pattern_idx; Type: INDEX; Schema: auth; Owner: -
--

CREATE INDEX "sso_providers_resource_id_pattern_idx" ON "auth"."sso_providers" USING "btree" ("resource_id" "text_pattern_ops");


--
-- Name: unique_phone_factor_per_user; Type: INDEX; Schema: auth; Owner: -
--

CREATE UNIQUE INDEX "unique_phone_factor_per_user" ON "auth"."mfa_factors" USING "btree" ("user_id", "phone");


--
-- Name: user_id_created_at_idx; Type: INDEX; Schema: auth; Owner: -
--

CREATE INDEX "user_id_created_at_idx" ON "auth"."sessions" USING "btree" ("user_id", "created_at");


--
-- Name: users_email_partial_key; Type: INDEX; Schema: auth; Owner: -
--

CREATE UNIQUE INDEX "users_email_partial_key" ON "auth"."users" USING "btree" ("email") WHERE ("is_sso_user" = false);


--
-- Name: INDEX "users_email_partial_key"; Type: COMMENT; Schema: auth; Owner: -
--

COMMENT ON INDEX "auth"."users_email_partial_key" IS 'Auth: A partial unique index that applies only when is_sso_user is false';


--
-- Name: users_instance_id_email_idx; Type: INDEX; Schema: auth; Owner: -
--

CREATE INDEX "users_instance_id_email_idx" ON "auth"."users" USING "btree" ("instance_id", "lower"(("email")::"text"));


--
-- Name: users_instance_id_idx; Type: INDEX; Schema: auth; Owner: -
--

CREATE INDEX "users_instance_id_idx" ON "auth"."users" USING "btree" ("instance_id");


--
-- Name: users_is_anonymous_idx; Type: INDEX; Schema: auth; Owner: -
--

CREATE INDEX "users_is_anonymous_idx" ON "auth"."users" USING "btree" ("is_anonymous");


--
-- Name: idx_adaptive_rules_active; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_adaptive_rules_active" ON "public"."adaptive_streaming_rules" USING "btree" ("active");


--
-- Name: idx_adaptive_rules_priority; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_adaptive_rules_priority" ON "public"."adaptive_streaming_rules" USING "btree" ("priority" DESC);


--
-- Name: idx_audit_logs_action; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_audit_logs_action" ON "public"."audit_logs" USING "btree" ("action");


--
-- Name: idx_audit_logs_course_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_audit_logs_course_id" ON "public"."audit_logs" USING "btree" ("course_id");


--
-- Name: idx_audit_logs_created_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_audit_logs_created_at" ON "public"."audit_logs" USING "btree" ("created_at");


--
-- Name: idx_audit_logs_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_audit_logs_user_id" ON "public"."audit_logs" USING "btree" ("user_id");


--
-- Name: idx_bandwidth_tests_session; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_bandwidth_tests_session" ON "public"."bandwidth_tests" USING "btree" ("session_id");


--
-- Name: idx_bandwidth_tests_timestamp; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_bandwidth_tests_timestamp" ON "public"."bandwidth_tests" USING "btree" ("timestamp");


--
-- Name: idx_bandwidth_tests_type; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_bandwidth_tests_type" ON "public"."bandwidth_tests" USING "btree" ("test_type");


--
-- Name: idx_bandwidth_tests_user; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_bandwidth_tests_user" ON "public"."bandwidth_tests" USING "btree" ("user_id");


--
-- Name: idx_chat_history_created_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_chat_history_created_at" ON "public"."chat_history" USING "btree" ("created_at");


--
-- Name: idx_chat_history_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_chat_history_user_id" ON "public"."chat_history" USING "btree" ("user_id");


--
-- Name: idx_course_access_analytics_course; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_course_access_analytics_course" ON "public"."course_access_logs" USING "btree" ("course_id", "access_granted", "access_type");


--
-- Name: idx_course_access_cache_access_level; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_course_access_cache_access_level" ON "public"."course_access_cache" USING "btree" ("access_level");


--
-- Name: idx_course_access_cache_expires_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_course_access_cache_expires_at" ON "public"."course_access_cache" USING "btree" ("expires_at");


--
-- Name: idx_course_access_cache_user_course; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_course_access_cache_user_course" ON "public"."course_access_cache" USING "btree" ("user_id", "course_id");


--
-- Name: idx_course_access_logs_access_type; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_course_access_logs_access_type" ON "public"."course_access_logs" USING "btree" ("access_type");


--
-- Name: idx_course_access_logs_created_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_course_access_logs_created_at" ON "public"."course_access_logs" USING "btree" ("created_at");


--
-- Name: idx_course_access_logs_user_course; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_course_access_logs_user_course" ON "public"."course_access_logs" USING "btree" ("user_id", "course_id");


--
-- Name: idx_course_access_logs_user_lecture; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_course_access_logs_user_lecture" ON "public"."course_access_logs" USING "btree" ("user_id", "lecture_id");


--
-- Name: idx_course_access_user_course_time; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_course_access_user_course_time" ON "public"."course_access_logs" USING "btree" ("user_id", "course_id", "created_at");


--
-- Name: idx_course_resources_course_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_course_resources_course_id" ON "public"."course_resources" USING "btree" ("course_id");


--
-- Name: idx_course_resources_file_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_course_resources_file_id" ON "public"."course_resources" USING "btree" ("file_id");


--
-- Name: idx_course_resources_order; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_course_resources_order" ON "public"."course_resources" USING "btree" ("course_id", "display_order");


--
-- Name: idx_course_resources_type; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_course_resources_type" ON "public"."course_resources" USING "btree" ("resource_type");


--
-- Name: idx_courses_category; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_courses_category" ON "public"."courses" USING "btree" ("category");


--
-- Name: idx_courses_created_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_courses_created_at" ON "public"."courses" USING "btree" ("created_at");


--
-- Name: idx_courses_deleted_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_courses_deleted_at" ON "public"."courses" USING "btree" ("deleted_at") WHERE ("deleted_at" IS NULL);


--
-- Name: idx_courses_instructor_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_courses_instructor_id" ON "public"."courses" USING "btree" ("instructor_id");


--
-- Name: idx_courses_is_paid; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_courses_is_paid" ON "public"."courses" USING "btree" ("is_paid");


--
-- Name: idx_courses_lemon_squeezy_product_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_courses_lemon_squeezy_product_id" ON "public"."courses" USING "btree" ("lemon_squeezy_product_id");


--
-- Name: idx_courses_lemon_squeezy_variant_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_courses_lemon_squeezy_variant_id" ON "public"."courses" USING "btree" ("lemon_squeezy_variant_id");


--
-- Name: idx_courses_level; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_courses_level" ON "public"."courses" USING "btree" ("level");


--
-- Name: idx_courses_preview_enabled; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_courses_preview_enabled" ON "public"."courses" USING "btree" ("preview_enabled");


--
-- Name: idx_courses_price; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_courses_price" ON "public"."courses" USING "btree" ("price");


--
-- Name: idx_courses_rating; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_courses_rating" ON "public"."courses" USING "btree" ("rating");


--
-- Name: idx_courses_search; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_courses_search" ON "public"."courses" USING "gin" ("to_tsvector"('"english"'::"regconfig", ((((("title")::"text" || ' '::"text") || "description") || ' '::"text") || ("category")::"text")));


--
-- Name: idx_courses_status; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_courses_status" ON "public"."courses" USING "btree" ("status");


--
-- Name: idx_courses_tags; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_courses_tags" ON "public"."courses" USING "gin" ("tags");


--
-- Name: idx_enrollment_user_payment_status; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_enrollment_user_payment_status" ON "public"."enrollments" USING "btree" ("user_id", "payment_status");


--
-- Name: idx_enrollments_access_expires_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_enrollments_access_expires_at" ON "public"."enrollments" USING "btree" ("access_expires_at");


--
-- Name: idx_enrollments_completed_lectures; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_enrollments_completed_lectures" ON "public"."enrollments" USING "btree" ("completed_lectures");


--
-- Name: idx_enrollments_course_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_enrollments_course_id" ON "public"."enrollments" USING "btree" ("course_id");


--
-- Name: idx_enrollments_created_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_enrollments_created_at" ON "public"."enrollments" USING "btree" ("created_at");


--
-- Name: idx_enrollments_deleted_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_enrollments_deleted_at" ON "public"."enrollments" USING "btree" ("deleted_at") WHERE ("deleted_at" IS NULL);


--
-- Name: idx_enrollments_enrolled_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_enrollments_enrolled_at" ON "public"."enrollments" USING "btree" ("enrolled_at");


--
-- Name: idx_enrollments_last_accessed; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_enrollments_last_accessed" ON "public"."enrollments" USING "btree" ("last_accessed");


--
-- Name: idx_enrollments_payment_status; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_enrollments_payment_status" ON "public"."enrollments" USING "btree" ("payment_status");


--
-- Name: idx_enrollments_payment_verified_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_enrollments_payment_verified_at" ON "public"."enrollments" USING "btree" ("payment_verified_at");


--
-- Name: idx_enrollments_progress; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_enrollments_progress" ON "public"."enrollments" USING "btree" ("progress_percentage");


--
-- Name: idx_enrollments_status; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_enrollments_status" ON "public"."enrollments" USING "btree" ("status");


--
-- Name: idx_enrollments_total_lectures; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_enrollments_total_lectures" ON "public"."enrollments" USING "btree" ("total_lectures");


--
-- Name: idx_enrollments_transaction_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_enrollments_transaction_id" ON "public"."enrollments" USING "btree" ("transaction_id");


--
-- Name: idx_enrollments_updated_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_enrollments_updated_at" ON "public"."enrollments" USING "btree" ("updated_at");


--
-- Name: idx_enrollments_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_enrollments_user_id" ON "public"."enrollments" USING "btree" ("user_id");


--
-- Name: idx_file_permissions_file_user; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_file_permissions_file_user" ON "public"."file_permissions" USING "btree" ("file_id", "user_id");


--
-- Name: idx_file_permissions_type; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_file_permissions_type" ON "public"."file_permissions" USING "btree" ("permission_type");


--
-- Name: idx_file_permissions_unique; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX "idx_file_permissions_unique" ON "public"."file_permissions" USING "btree" ("file_id", "user_id", "permission_type");


--
-- Name: idx_file_permissions_user; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_file_permissions_user" ON "public"."file_permissions" USING "btree" ("user_id");


--
-- Name: idx_files_bucket_key; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_files_bucket_key" ON "public"."files" USING "btree" ("bucket_name", "object_key");


--
-- Name: idx_files_bucket_object_active; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX "idx_files_bucket_object_active" ON "public"."files" USING "btree" ("bucket_name", "object_key") WHERE ("deleted_at" IS NULL);


--
-- Name: idx_files_content_type; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_files_content_type" ON "public"."files" USING "btree" ("content_type");


--
-- Name: idx_files_created_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_files_created_at" ON "public"."files" USING "btree" ("created_at");


--
-- Name: idx_files_deleted_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_files_deleted_at" ON "public"."files" USING "btree" ("deleted_at");


--
-- Name: idx_files_is_public; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_files_is_public" ON "public"."files" USING "btree" ("is_public");


--
-- Name: idx_files_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_files_user_id" ON "public"."files" USING "btree" ("upload_user_id");


--
-- Name: idx_forum_mentions_mentioned_user; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_forum_mentions_mentioned_user" ON "public"."forum_mentions" USING "btree" ("mentioned_user_id");


--
-- Name: idx_forum_mentions_mentioner_user; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_forum_mentions_mentioner_user" ON "public"."forum_mentions" USING "btree" ("mentioner_user_id");


--
-- Name: idx_forum_mentions_post_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_forum_mentions_post_id" ON "public"."forum_mentions" USING "btree" ("post_id");


--
-- Name: idx_forum_mentions_unread; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_forum_mentions_unread" ON "public"."forum_mentions" USING "btree" ("mentioned_user_id", "is_read") WHERE ("is_read" = false);


--
-- Name: idx_forum_notifications_reference; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_forum_notifications_reference" ON "public"."forum_notifications" USING "btree" ("reference_type", "reference_id");


--
-- Name: idx_forum_notifications_type; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_forum_notifications_type" ON "public"."forum_notifications" USING "btree" ("type");


--
-- Name: idx_forum_notifications_unread; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_forum_notifications_unread" ON "public"."forum_notifications" USING "btree" ("user_id", "is_read") WHERE ("is_read" = false);


--
-- Name: idx_forum_notifications_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_forum_notifications_user_id" ON "public"."forum_notifications" USING "btree" ("user_id");


--
-- Name: idx_forum_posts_author_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_forum_posts_author_id" ON "public"."forum_posts" USING "btree" ("author_id");


--
-- Name: idx_forum_posts_created_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_forum_posts_created_at" ON "public"."forum_posts" USING "btree" ("created_at" DESC);


--
-- Name: idx_forum_posts_is_answer; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_forum_posts_is_answer" ON "public"."forum_posts" USING "btree" ("is_answer") WHERE ("is_answer" = true);


--
-- Name: idx_forum_posts_is_pinned; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_forum_posts_is_pinned" ON "public"."forum_posts" USING "btree" ("is_pinned") WHERE ("is_pinned" = true);


--
-- Name: idx_forum_posts_parent_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_forum_posts_parent_id" ON "public"."forum_posts" USING "btree" ("parent_id") WHERE ("parent_id" IS NOT NULL);


--
-- Name: idx_forum_posts_pin_order; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_forum_posts_pin_order" ON "public"."forum_posts" USING "btree" ("pin_order") WHERE ("pin_order" IS NOT NULL);


--
-- Name: idx_forum_posts_status; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_forum_posts_status" ON "public"."forum_posts" USING "btree" ("status");


--
-- Name: idx_forum_posts_topic_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_forum_posts_topic_id" ON "public"."forum_posts" USING "btree" ("topic_id");


--
-- Name: idx_forum_posts_topic_status; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_forum_posts_topic_status" ON "public"."forum_posts" USING "btree" ("topic_id", "status");


--
-- Name: idx_forum_subscriptions_topic_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_forum_subscriptions_topic_id" ON "public"."forum_topic_subscriptions" USING "btree" ("topic_id");


--
-- Name: idx_forum_subscriptions_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_forum_subscriptions_user_id" ON "public"."forum_topic_subscriptions" USING "btree" ("user_id");


--
-- Name: idx_forum_topics_course_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_forum_topics_course_id" ON "public"."forum_topics" USING "btree" ("course_id");


--
-- Name: idx_forum_topics_course_status; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_forum_topics_course_status" ON "public"."forum_topics" USING "btree" ("course_id", "status") WHERE ("course_id" IS NOT NULL);


--
-- Name: idx_forum_topics_creator_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_forum_topics_creator_id" ON "public"."forum_topics" USING "btree" ("creator_id");


--
-- Name: idx_forum_topics_is_pinned; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_forum_topics_is_pinned" ON "public"."forum_topics" USING "btree" ("is_pinned");


--
-- Name: idx_forum_topics_pin_order; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_forum_topics_pin_order" ON "public"."forum_topics" USING "btree" ("pin_order") WHERE ("pin_order" IS NOT NULL);


--
-- Name: idx_forum_topics_status; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_forum_topics_status" ON "public"."forum_topics" USING "btree" ("status");


--
-- Name: idx_forum_votes_post_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_forum_votes_post_id" ON "public"."forum_votes" USING "btree" ("post_id");


--
-- Name: idx_forum_votes_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_forum_votes_user_id" ON "public"."forum_votes" USING "btree" ("user_id");


--
-- Name: idx_forum_votes_vote_type; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_forum_votes_vote_type" ON "public"."forum_votes" USING "btree" ("vote_type");


--
-- Name: idx_lecture_preview_sessions_exhausted; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_lecture_preview_sessions_exhausted" ON "public"."lecture_preview_sessions" USING "btree" ("preview_exhausted");


--
-- Name: idx_lecture_preview_sessions_last_accessed; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_lecture_preview_sessions_last_accessed" ON "public"."lecture_preview_sessions" USING "btree" ("last_accessed_at");


--
-- Name: idx_lecture_preview_sessions_user_lecture; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_lecture_preview_sessions_user_lecture" ON "public"."lecture_preview_sessions" USING "btree" ("user_id", "lecture_id");


--
-- Name: idx_lecture_resources_file_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_lecture_resources_file_id" ON "public"."lecture_resources" USING "btree" ("file_id");


--
-- Name: idx_lecture_resources_lecture_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_lecture_resources_lecture_id" ON "public"."lecture_resources" USING "btree" ("lecture_id");


--
-- Name: idx_lecture_resources_order; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_lecture_resources_order" ON "public"."lecture_resources" USING "btree" ("lecture_id", "display_order");


--
-- Name: idx_lecture_resources_type; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_lecture_resources_type" ON "public"."lecture_resources" USING "btree" ("resource_type");


--
-- Name: idx_lectures_access_level; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_lectures_access_level" ON "public"."lectures" USING "btree" ("access_level");


--
-- Name: idx_lectures_course_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_lectures_course_id" ON "public"."lectures" USING "btree" ("course_id");


--
-- Name: idx_lectures_deleted_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_lectures_deleted_at" ON "public"."lectures" USING "btree" ("deleted_at") WHERE ("deleted_at" IS NULL);


--
-- Name: idx_lectures_order_number; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_lectures_order_number" ON "public"."lectures" USING "btree" ("course_id", "order_number");


--
-- Name: idx_lectures_preview_available; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_lectures_preview_available" ON "public"."lectures" USING "btree" ("preview_available");


--
-- Name: idx_lectures_status; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_lectures_status" ON "public"."lectures" USING "btree" ("status");


--
-- Name: idx_lemon_squeezy_products_product_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_lemon_squeezy_products_product_id" ON "public"."lemon_squeezy_products" USING "btree" ("lemon_squeezy_product_id");


--
-- Name: idx_lemon_squeezy_variants_product_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_lemon_squeezy_variants_product_id" ON "public"."lemon_squeezy_variants" USING "btree" ("lemon_squeezy_product_id");


--
-- Name: idx_lemon_squeezy_variants_variant_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_lemon_squeezy_variants_variant_id" ON "public"."lemon_squeezy_variants" USING "btree" ("lemon_squeezy_variant_id");


--
-- Name: idx_network_analytics_date; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_network_analytics_date" ON "public"."network_analytics_daily" USING "btree" ("date");


--
-- Name: idx_network_analytics_user; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_network_analytics_user" ON "public"."network_analytics_daily" USING "btree" ("user_id");


--
-- Name: idx_network_analytics_video; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_network_analytics_video" ON "public"."network_analytics_daily" USING "btree" ("video_id");


--
-- Name: idx_network_dashboard_session_status; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_network_dashboard_session_status" ON "public"."network_monitoring_dashboard" USING "btree" ("session_status");


--
-- Name: idx_network_dashboard_video_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_network_dashboard_video_id" ON "public"."network_monitoring_dashboard" USING "btree" ("video_id");


--
-- Name: idx_network_events_session; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_network_events_session" ON "public"."network_events" USING "btree" ("session_id");


--
-- Name: idx_network_events_severity; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_network_events_severity" ON "public"."network_events" USING "btree" ("severity");


--
-- Name: idx_network_events_timestamp; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_network_events_timestamp" ON "public"."network_events" USING "btree" ("timestamp");


--
-- Name: idx_network_events_type; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_network_events_type" ON "public"."network_events" USING "btree" ("event_type");


--
-- Name: idx_network_events_unresolved; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_network_events_unresolved" ON "public"."network_events" USING "btree" ("resolved") WHERE ("resolved" = false);


--
-- Name: idx_network_events_user; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_network_events_user" ON "public"."network_events" USING "btree" ("user_id");


--
-- Name: idx_network_metrics_connection_type; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_network_metrics_connection_type" ON "public"."network_metrics" USING "btree" ("connection_type");


--
-- Name: idx_network_metrics_device_type; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_network_metrics_device_type" ON "public"."network_metrics" USING "btree" ("device_type");


--
-- Name: idx_network_metrics_quality_score; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_network_metrics_quality_score" ON "public"."network_metrics" USING "btree" ("quality_score");


--
-- Name: idx_network_metrics_session; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_network_metrics_session" ON "public"."network_metrics" USING "btree" ("session_id");


--
-- Name: idx_network_metrics_timestamp; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_network_metrics_timestamp" ON "public"."network_metrics" USING "btree" ("timestamp");


--
-- Name: idx_network_metrics_user_timestamp; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_network_metrics_user_timestamp" ON "public"."network_metrics" USING "btree" ("user_id", "timestamp");


--
-- Name: idx_notes_course_lecture; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_notes_course_lecture" ON "public"."notes" USING "btree" ("course_id", "lecture_id");


--
-- Name: idx_notes_created_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_notes_created_at" ON "public"."notes" USING "btree" ("created_at");


--
-- Name: idx_notes_timestamp; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_notes_timestamp" ON "public"."notes" USING "btree" ("timestamp_seconds");


--
-- Name: idx_notes_user_course_lecture; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_notes_user_course_lecture" ON "public"."notes" USING "btree" ("user_id", "course_id", "lecture_id");


--
-- Name: idx_oauth_accounts_provider; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_oauth_accounts_provider" ON "public"."oauth_accounts" USING "btree" ("provider");


--
-- Name: idx_oauth_accounts_provider_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_oauth_accounts_provider_id" ON "public"."oauth_accounts" USING "btree" ("provider_id");


--
-- Name: idx_oauth_accounts_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_oauth_accounts_user_id" ON "public"."oauth_accounts" USING "btree" ("user_id");


--
-- Name: idx_payment_analytics_date; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_payment_analytics_date" ON "public"."transactions" USING "btree" ("date_trunc"('day'::"text", "created_at"));


--
-- Name: idx_payment_events_provider_event_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_payment_events_provider_event_id" ON "public"."payment_events" USING "btree" ("provider", "provider_event_id");


--
-- Name: idx_payment_events_transaction_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_payment_events_transaction_id" ON "public"."payment_events" USING "btree" ("transaction_id");


--
-- Name: idx_payment_events_type_processed; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_payment_events_type_processed" ON "public"."payment_events" USING "btree" ("event_type", "processed");


--
-- Name: idx_payment_events_user_course; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_payment_events_user_course" ON "public"."payment_events" USING "btree" ("user_id", "course_id");


--
-- Name: idx_payment_methods_stripe_customer_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_payment_methods_stripe_customer_id" ON "public"."payment_methods" USING "btree" ("stripe_customer_id");


--
-- Name: idx_payment_methods_stripe_payment_method_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_payment_methods_stripe_payment_method_id" ON "public"."payment_methods" USING "btree" ("stripe_payment_method_id");


--
-- Name: idx_payment_methods_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_payment_methods_user_id" ON "public"."payment_methods" USING "btree" ("user_id");


--
-- Name: idx_progress_completed; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_progress_completed" ON "public"."progress" USING "btree" ("completed");


--
-- Name: idx_progress_completed_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_progress_completed_at" ON "public"."progress" USING "btree" ("completed_at");


--
-- Name: idx_progress_course_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_progress_course_id" ON "public"."progress" USING "btree" ("course_id");


--
-- Name: idx_progress_is_completed; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_progress_is_completed" ON "public"."progress" USING "btree" ("is_completed");


--
-- Name: idx_progress_last_accessed; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_progress_last_accessed" ON "public"."progress" USING "btree" ("last_accessed");


--
-- Name: idx_progress_lecture_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_progress_lecture_id" ON "public"."progress" USING "btree" ("lecture_id");


--
-- Name: idx_progress_progress_percentage; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_progress_progress_percentage" ON "public"."progress" USING "btree" ("progress_percentage");


--
-- Name: idx_progress_user_course; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_progress_user_course" ON "public"."progress" USING "btree" ("user_id", "course_id");


--
-- Name: idx_progress_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_progress_user_id" ON "public"."progress" USING "btree" ("user_id");


--
-- Name: idx_progress_watch_time; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_progress_watch_time" ON "public"."progress" USING "btree" ("watch_time_seconds");


--
-- Name: idx_stripe_customers_stripe_customer_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_stripe_customers_stripe_customer_id" ON "public"."stripe_customers" USING "btree" ("stripe_customer_id");


--
-- Name: idx_stripe_customers_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_stripe_customers_user_id" ON "public"."stripe_customers" USING "btree" ("user_id");


--
-- Name: idx_stripe_products_course_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_stripe_products_course_id" ON "public"."stripe_products" USING "btree" ("course_id");


--
-- Name: idx_stripe_products_stripe_price_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_stripe_products_stripe_price_id" ON "public"."stripe_products" USING "btree" ("stripe_price_id");


--
-- Name: idx_stripe_products_stripe_product_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_stripe_products_stripe_product_id" ON "public"."stripe_products" USING "btree" ("stripe_product_id");


--
-- Name: idx_stripe_webhook_events_created_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_stripe_webhook_events_created_at" ON "public"."stripe_webhook_events" USING "btree" ("created_at");


--
-- Name: idx_stripe_webhook_events_event_type; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_stripe_webhook_events_event_type" ON "public"."stripe_webhook_events" USING "btree" ("event_type");


--
-- Name: idx_stripe_webhook_events_processed; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_stripe_webhook_events_processed" ON "public"."stripe_webhook_events" USING "btree" ("processed");


--
-- Name: idx_stripe_webhook_events_processing_attempts; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_stripe_webhook_events_processing_attempts" ON "public"."stripe_webhook_events" USING "btree" ("processing_attempts");


--
-- Name: idx_stripe_webhook_events_stripe_event_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_stripe_webhook_events_stripe_event_id" ON "public"."stripe_webhook_events" USING "btree" ("stripe_event_id");


--
-- Name: idx_stripe_webhook_events_type; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_stripe_webhook_events_type" ON "public"."stripe_webhook_events" USING "btree" ("event_type");


--
-- Name: idx_stripe_webhook_events_unprocessed; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_stripe_webhook_events_unprocessed" ON "public"."stripe_webhook_events" USING "btree" ("processed", "processing_attempts", "created_at") WHERE ("processed" = false);


--
-- Name: idx_subscriptions_status; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_subscriptions_status" ON "public"."subscriptions" USING "btree" ("status");


--
-- Name: idx_subscriptions_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_subscriptions_user_id" ON "public"."subscriptions" USING "btree" ("user_id");


--
-- Name: idx_transaction_user_status_course; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_transaction_user_status_course" ON "public"."transactions" USING "btree" ("user_id", "status", "course_id");


--
-- Name: idx_transactions_completed; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_transactions_completed" ON "public"."transactions" USING "btree" ("user_id", "course_id", "payment_verified_at") WHERE (("status")::"text" = 'completed'::"text");


--
-- Name: idx_transactions_course_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_transactions_course_id" ON "public"."transactions" USING "btree" ("course_id");


--
-- Name: idx_transactions_lemon_squeezy_checkout_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_transactions_lemon_squeezy_checkout_id" ON "public"."transactions" USING "btree" ("lemon_squeezy_checkout_id");


--
-- Name: idx_transactions_lemon_squeezy_order_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_transactions_lemon_squeezy_order_id" ON "public"."transactions" USING "btree" ("lemon_squeezy_order_id");


--
-- Name: idx_transactions_payment_provider; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_transactions_payment_provider" ON "public"."transactions" USING "btree" ("payment_provider");


--
-- Name: idx_transactions_payment_verified_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_transactions_payment_verified_at" ON "public"."transactions" USING "btree" ("payment_verified_at");


--
-- Name: idx_transactions_pending; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_transactions_pending" ON "public"."transactions" USING "btree" ("user_id", "course_id", "created_at") WHERE (("status")::"text" = 'pending'::"text");


--
-- Name: idx_transactions_status; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_transactions_status" ON "public"."transactions" USING "btree" ("status");


--
-- Name: idx_transactions_status_updated; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_transactions_status_updated" ON "public"."transactions" USING "btree" ("status", "updated_at");


--
-- Name: idx_transactions_stripe_charge_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_transactions_stripe_charge_id" ON "public"."transactions" USING "btree" ("stripe_charge_id");


--
-- Name: idx_transactions_stripe_customer_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_transactions_stripe_customer_id" ON "public"."transactions" USING "btree" ("stripe_customer_id");


--
-- Name: idx_transactions_stripe_payment_intent_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_transactions_stripe_payment_intent_id" ON "public"."transactions" USING "btree" ("stripe_payment_intent_id");


--
-- Name: idx_transactions_stripe_session_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_transactions_stripe_session_id" ON "public"."transactions" USING "btree" ("stripe_session_id");


--
-- Name: idx_transactions_user_course; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_transactions_user_course" ON "public"."transactions" USING "btree" ("user_id", "course_id");


--
-- Name: idx_transactions_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_transactions_user_id" ON "public"."transactions" USING "btree" ("user_id");


--
-- Name: idx_transactions_webhook_event_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_transactions_webhook_event_id" ON "public"."transactions" USING "btree" ("webhook_event_id");


--
-- Name: idx_upload_sessions_bucket_key; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_upload_sessions_bucket_key" ON "public"."upload_sessions" USING "btree" ("bucket_name", "object_key");


--
-- Name: idx_upload_sessions_expires_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_upload_sessions_expires_at" ON "public"."upload_sessions" USING "btree" ("expires_at");


--
-- Name: idx_upload_sessions_upload_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_upload_sessions_upload_id" ON "public"."upload_sessions" USING "btree" ("upload_id");


--
-- Name: idx_upload_sessions_upload_id_unique; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX "idx_upload_sessions_upload_id_unique" ON "public"."upload_sessions" USING "btree" ("upload_id");


--
-- Name: idx_upload_sessions_user_status; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_upload_sessions_user_status" ON "public"."upload_sessions" USING "btree" ("user_id", "status");


--
-- Name: idx_user_payment_methods_active; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_user_payment_methods_active" ON "public"."user_payment_methods" USING "btree" ("user_id", "is_active") WHERE ("is_active" = true);


--
-- Name: idx_user_payment_methods_default; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_user_payment_methods_default" ON "public"."user_payment_methods" USING "btree" ("user_id", "is_default") WHERE ("is_default" = true);


--
-- Name: idx_user_payment_methods_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_user_payment_methods_user_id" ON "public"."user_payment_methods" USING "btree" ("user_id");


--
-- Name: idx_users_email; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_users_email" ON "public"."users" USING "btree" ("email");


--
-- Name: idx_users_provider_id; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX "idx_users_provider_id" ON "public"."users" USING "btree" ("provider", "provider_id") WHERE (("provider")::"text" <> 'local'::"text");


--
-- Name: idx_users_role; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_users_role" ON "public"."users" USING "btree" ("role");


--
-- Name: idx_users_username; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_users_username" ON "public"."users" USING "btree" ("username");


--
-- Name: idx_video_analytics_video_date; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_video_analytics_video_date" ON "public"."video_analytics" USING "btree" ("video_id", "date");


--
-- Name: idx_video_permissions_video_user; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_video_permissions_video_user" ON "public"."video_permissions" USING "btree" ("video_id", "user_id");


--
-- Name: idx_video_qualities_video_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_video_qualities_video_id" ON "public"."video_qualities" USING "btree" ("video_id");


--
-- Name: idx_videos_cloudflare_uid; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_videos_cloudflare_uid" ON "public"."videos" USING "btree" ("cloudflare_uid");


--
-- Name: idx_videos_course_lecture; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_videos_course_lecture" ON "public"."videos" USING "btree" ("course_id", "lecture_id");


--
-- Name: idx_videos_deleted_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_videos_deleted_at" ON "public"."videos" USING "btree" ("deleted_at");


--
-- Name: idx_videos_search; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_videos_search" ON "public"."videos" USING "gin" ("to_tsvector"('"english"'::"regconfig", ((("title")::"text" || ' '::"text") || "description"))) WHERE ("deleted_at" IS NULL);


--
-- Name: idx_videos_status; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_videos_status" ON "public"."videos" USING "btree" ("status");


--
-- Name: idx_videos_status_visibility; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_videos_status_visibility" ON "public"."videos" USING "btree" ("status", "visibility") WHERE ("deleted_at" IS NULL);


--
-- Name: idx_videos_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_videos_user_id" ON "public"."videos" USING "btree" ("upload_user_id");


--
-- Name: idx_viewing_sessions_session; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_viewing_sessions_session" ON "public"."viewing_sessions" USING "btree" ("session_id");


--
-- Name: idx_viewing_sessions_user_video; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_viewing_sessions_user_video" ON "public"."viewing_sessions" USING "btree" ("user_id", "video_id");


--
-- Name: idx_webhook_events_event_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_webhook_events_event_id" ON "public"."lemon_squeezy_webhook_events" USING "btree" ("event_id");


--
-- Name: idx_webhook_events_event_name; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_webhook_events_event_name" ON "public"."lemon_squeezy_webhook_events" USING "btree" ("event_name");


--
-- Name: idx_webhook_events_processed_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "idx_webhook_events_processed_at" ON "public"."lemon_squeezy_webhook_events" USING "btree" ("processed_at");


--
-- Name: ix_realtime_subscription_entity; Type: INDEX; Schema: realtime; Owner: -
--

CREATE INDEX "ix_realtime_subscription_entity" ON "realtime"."subscription" USING "btree" ("entity");


--
-- Name: messages_inserted_at_topic_index; Type: INDEX; Schema: realtime; Owner: -
--

CREATE INDEX "messages_inserted_at_topic_index" ON ONLY "realtime"."messages" USING "btree" ("inserted_at" DESC, "topic") WHERE (("extension" = 'broadcast'::"text") AND ("private" IS TRUE));


--
-- Name: messages_2025_10_04_inserted_at_topic_idx; Type: INDEX; Schema: realtime; Owner: -
--

CREATE INDEX "messages_2025_10_04_inserted_at_topic_idx" ON "realtime"."messages_2025_10_04" USING "btree" ("inserted_at" DESC, "topic") WHERE (("extension" = 'broadcast'::"text") AND ("private" IS TRUE));


--
-- Name: messages_2025_10_05_inserted_at_topic_idx; Type: INDEX; Schema: realtime; Owner: -
--

CREATE INDEX "messages_2025_10_05_inserted_at_topic_idx" ON "realtime"."messages_2025_10_05" USING "btree" ("inserted_at" DESC, "topic") WHERE (("extension" = 'broadcast'::"text") AND ("private" IS TRUE));


--
-- Name: messages_2025_10_06_inserted_at_topic_idx; Type: INDEX; Schema: realtime; Owner: -
--

CREATE INDEX "messages_2025_10_06_inserted_at_topic_idx" ON "realtime"."messages_2025_10_06" USING "btree" ("inserted_at" DESC, "topic") WHERE (("extension" = 'broadcast'::"text") AND ("private" IS TRUE));


--
-- Name: messages_2025_10_07_inserted_at_topic_idx; Type: INDEX; Schema: realtime; Owner: -
--

CREATE INDEX "messages_2025_10_07_inserted_at_topic_idx" ON "realtime"."messages_2025_10_07" USING "btree" ("inserted_at" DESC, "topic") WHERE (("extension" = 'broadcast'::"text") AND ("private" IS TRUE));


--
-- Name: messages_2025_10_08_inserted_at_topic_idx; Type: INDEX; Schema: realtime; Owner: -
--

CREATE INDEX "messages_2025_10_08_inserted_at_topic_idx" ON "realtime"."messages_2025_10_08" USING "btree" ("inserted_at" DESC, "topic") WHERE (("extension" = 'broadcast'::"text") AND ("private" IS TRUE));


--
-- Name: messages_2025_10_09_inserted_at_topic_idx; Type: INDEX; Schema: realtime; Owner: -
--

CREATE INDEX "messages_2025_10_09_inserted_at_topic_idx" ON "realtime"."messages_2025_10_09" USING "btree" ("inserted_at" DESC, "topic") WHERE (("extension" = 'broadcast'::"text") AND ("private" IS TRUE));


--
-- Name: messages_2025_10_10_inserted_at_topic_idx; Type: INDEX; Schema: realtime; Owner: -
--

CREATE INDEX "messages_2025_10_10_inserted_at_topic_idx" ON "realtime"."messages_2025_10_10" USING "btree" ("inserted_at" DESC, "topic") WHERE (("extension" = 'broadcast'::"text") AND ("private" IS TRUE));


--
-- Name: subscription_subscription_id_entity_filters_key; Type: INDEX; Schema: realtime; Owner: -
--

CREATE UNIQUE INDEX "subscription_subscription_id_entity_filters_key" ON "realtime"."subscription" USING "btree" ("subscription_id", "entity", "filters");


--
-- Name: bname; Type: INDEX; Schema: storage; Owner: -
--

CREATE UNIQUE INDEX "bname" ON "storage"."buckets" USING "btree" ("name");


--
-- Name: bucketid_objname; Type: INDEX; Schema: storage; Owner: -
--

CREATE UNIQUE INDEX "bucketid_objname" ON "storage"."objects" USING "btree" ("bucket_id", "name");


--
-- Name: idx_iceberg_namespaces_bucket_id; Type: INDEX; Schema: storage; Owner: -
--

CREATE UNIQUE INDEX "idx_iceberg_namespaces_bucket_id" ON "storage"."iceberg_namespaces" USING "btree" ("bucket_id", "name");


--
-- Name: idx_iceberg_tables_namespace_id; Type: INDEX; Schema: storage; Owner: -
--

CREATE UNIQUE INDEX "idx_iceberg_tables_namespace_id" ON "storage"."iceberg_tables" USING "btree" ("namespace_id", "name");


--
-- Name: idx_multipart_uploads_list; Type: INDEX; Schema: storage; Owner: -
--

CREATE INDEX "idx_multipart_uploads_list" ON "storage"."s3_multipart_uploads" USING "btree" ("bucket_id", "key", "created_at");


--
-- Name: idx_name_bucket_level_unique; Type: INDEX; Schema: storage; Owner: -
--

CREATE UNIQUE INDEX "idx_name_bucket_level_unique" ON "storage"."objects" USING "btree" ("name" COLLATE "C", "bucket_id", "level");


--
-- Name: idx_objects_bucket_id_name; Type: INDEX; Schema: storage; Owner: -
--

CREATE INDEX "idx_objects_bucket_id_name" ON "storage"."objects" USING "btree" ("bucket_id", "name" COLLATE "C");


--
-- Name: idx_objects_lower_name; Type: INDEX; Schema: storage; Owner: -
--

CREATE INDEX "idx_objects_lower_name" ON "storage"."objects" USING "btree" (("path_tokens"["level"]), "lower"("name") "text_pattern_ops", "bucket_id", "level");


--
-- Name: idx_prefixes_lower_name; Type: INDEX; Schema: storage; Owner: -
--

CREATE INDEX "idx_prefixes_lower_name" ON "storage"."prefixes" USING "btree" ("bucket_id", "level", (("string_to_array"("name", '/'::"text"))["level"]), "lower"("name") "text_pattern_ops");


--
-- Name: name_prefix_search; Type: INDEX; Schema: storage; Owner: -
--

CREATE INDEX "name_prefix_search" ON "storage"."objects" USING "btree" ("name" "text_pattern_ops");


--
-- Name: objects_bucket_id_level_idx; Type: INDEX; Schema: storage; Owner: -
--

CREATE UNIQUE INDEX "objects_bucket_id_level_idx" ON "storage"."objects" USING "btree" ("bucket_id", "level", "name" COLLATE "C");


--
-- Name: supabase_functions_hooks_h_table_id_h_name_idx; Type: INDEX; Schema: supabase_functions; Owner: -
--

CREATE INDEX "supabase_functions_hooks_h_table_id_h_name_idx" ON "supabase_functions"."hooks" USING "btree" ("hook_table_id", "hook_name");


--
-- Name: supabase_functions_hooks_request_id_idx; Type: INDEX; Schema: supabase_functions; Owner: -
--

CREATE INDEX "supabase_functions_hooks_request_id_idx" ON "supabase_functions"."hooks" USING "btree" ("request_id");


--
-- Name: messages_2025_10_04_inserted_at_topic_idx; Type: INDEX ATTACH; Schema: realtime; Owner: -
--

ALTER INDEX "realtime"."messages_inserted_at_topic_index" ATTACH PARTITION "realtime"."messages_2025_10_04_inserted_at_topic_idx";


--
-- Name: messages_2025_10_04_pkey; Type: INDEX ATTACH; Schema: realtime; Owner: -
--

ALTER INDEX "realtime"."messages_pkey" ATTACH PARTITION "realtime"."messages_2025_10_04_pkey";


--
-- Name: messages_2025_10_05_inserted_at_topic_idx; Type: INDEX ATTACH; Schema: realtime; Owner: -
--

ALTER INDEX "realtime"."messages_inserted_at_topic_index" ATTACH PARTITION "realtime"."messages_2025_10_05_inserted_at_topic_idx";


--
-- Name: messages_2025_10_05_pkey; Type: INDEX ATTACH; Schema: realtime; Owner: -
--

ALTER INDEX "realtime"."messages_pkey" ATTACH PARTITION "realtime"."messages_2025_10_05_pkey";


--
-- Name: messages_2025_10_06_inserted_at_topic_idx; Type: INDEX ATTACH; Schema: realtime; Owner: -
--

ALTER INDEX "realtime"."messages_inserted_at_topic_index" ATTACH PARTITION "realtime"."messages_2025_10_06_inserted_at_topic_idx";


--
-- Name: messages_2025_10_06_pkey; Type: INDEX ATTACH; Schema: realtime; Owner: -
--

ALTER INDEX "realtime"."messages_pkey" ATTACH PARTITION "realtime"."messages_2025_10_06_pkey";


--
-- Name: messages_2025_10_07_inserted_at_topic_idx; Type: INDEX ATTACH; Schema: realtime; Owner: -
--

ALTER INDEX "realtime"."messages_inserted_at_topic_index" ATTACH PARTITION "realtime"."messages_2025_10_07_inserted_at_topic_idx";


--
-- Name: messages_2025_10_07_pkey; Type: INDEX ATTACH; Schema: realtime; Owner: -
--

ALTER INDEX "realtime"."messages_pkey" ATTACH PARTITION "realtime"."messages_2025_10_07_pkey";


--
-- Name: messages_2025_10_08_inserted_at_topic_idx; Type: INDEX ATTACH; Schema: realtime; Owner: -
--

ALTER INDEX "realtime"."messages_inserted_at_topic_index" ATTACH PARTITION "realtime"."messages_2025_10_08_inserted_at_topic_idx";


--
-- Name: messages_2025_10_08_pkey; Type: INDEX ATTACH; Schema: realtime; Owner: -
--

ALTER INDEX "realtime"."messages_pkey" ATTACH PARTITION "realtime"."messages_2025_10_08_pkey";


--
-- Name: messages_2025_10_09_inserted_at_topic_idx; Type: INDEX ATTACH; Schema: realtime; Owner: -
--

ALTER INDEX "realtime"."messages_inserted_at_topic_index" ATTACH PARTITION "realtime"."messages_2025_10_09_inserted_at_topic_idx";


--
-- Name: messages_2025_10_09_pkey; Type: INDEX ATTACH; Schema: realtime; Owner: -
--

ALTER INDEX "realtime"."messages_pkey" ATTACH PARTITION "realtime"."messages_2025_10_09_pkey";


--
-- Name: messages_2025_10_10_inserted_at_topic_idx; Type: INDEX ATTACH; Schema: realtime; Owner: -
--

ALTER INDEX "realtime"."messages_inserted_at_topic_index" ATTACH PARTITION "realtime"."messages_2025_10_10_inserted_at_topic_idx";


--
-- Name: messages_2025_10_10_pkey; Type: INDEX ATTACH; Schema: realtime; Owner: -
--

ALTER INDEX "realtime"."messages_pkey" ATTACH PARTITION "realtime"."messages_2025_10_10_pkey";


--
-- Name: notes trigger_notes_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER "trigger_notes_updated_at" BEFORE UPDATE ON "public"."notes" FOR EACH ROW EXECUTE FUNCTION "public"."update_notes_updated_at"();


--
-- Name: course_resources trigger_prevent_course_resources_insert; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER "trigger_prevent_course_resources_insert" BEFORE INSERT ON "public"."course_resources" FOR EACH ROW EXECUTE FUNCTION "public"."prevent_course_resources_insert"();


--
-- Name: transactions trigger_update_enrollment_payment_status; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER "trigger_update_enrollment_payment_status" AFTER UPDATE ON "public"."transactions" FOR EACH ROW EXECUTE FUNCTION "public"."update_enrollment_payment_status"();


--
-- Name: network_metrics trigger_update_network_analytics; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER "trigger_update_network_analytics" AFTER INSERT ON "public"."network_metrics" FOR EACH ROW EXECUTE FUNCTION "public"."update_network_analytics_daily"();


--
-- Name: enrollments update_enrollments_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER "update_enrollments_updated_at" BEFORE UPDATE ON "public"."enrollments" FOR EACH ROW EXECUTE FUNCTION "public"."update_updated_at_column"();


--
-- Name: progress update_progress_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER "update_progress_updated_at" BEFORE UPDATE ON "public"."progress" FOR EACH ROW EXECUTE FUNCTION "public"."update_updated_at_column"();


--
-- Name: stripe_customers update_stripe_customers_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER "update_stripe_customers_updated_at" BEFORE UPDATE ON "public"."stripe_customers" FOR EACH ROW EXECUTE FUNCTION "public"."update_stripe_customers_updated_at"();


--
-- Name: stripe_products update_stripe_products_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER "update_stripe_products_updated_at" BEFORE UPDATE ON "public"."stripe_products" FOR EACH ROW EXECUTE FUNCTION "public"."update_stripe_products_updated_at"();


--
-- Name: subscription tr_check_filters; Type: TRIGGER; Schema: realtime; Owner: -
--

CREATE TRIGGER "tr_check_filters" BEFORE INSERT OR UPDATE ON "realtime"."subscription" FOR EACH ROW EXECUTE FUNCTION "realtime"."subscription_check_filters"();


--
-- Name: buckets enforce_bucket_name_length_trigger; Type: TRIGGER; Schema: storage; Owner: -
--

CREATE TRIGGER "enforce_bucket_name_length_trigger" BEFORE INSERT OR UPDATE OF "name" ON "storage"."buckets" FOR EACH ROW EXECUTE FUNCTION "storage"."enforce_bucket_name_length"();


--
-- Name: objects objects_delete_delete_prefix; Type: TRIGGER; Schema: storage; Owner: -
--

CREATE TRIGGER "objects_delete_delete_prefix" AFTER DELETE ON "storage"."objects" FOR EACH ROW EXECUTE FUNCTION "storage"."delete_prefix_hierarchy_trigger"();


--
-- Name: objects objects_insert_create_prefix; Type: TRIGGER; Schema: storage; Owner: -
--

CREATE TRIGGER "objects_insert_create_prefix" BEFORE INSERT ON "storage"."objects" FOR EACH ROW EXECUTE FUNCTION "storage"."objects_insert_prefix_trigger"();


--
-- Name: objects objects_update_create_prefix; Type: TRIGGER; Schema: storage; Owner: -
--

CREATE TRIGGER "objects_update_create_prefix" BEFORE UPDATE ON "storage"."objects" FOR EACH ROW WHEN ((("new"."name" <> "old"."name") OR ("new"."bucket_id" <> "old"."bucket_id"))) EXECUTE FUNCTION "storage"."objects_update_prefix_trigger"();


--
-- Name: prefixes prefixes_create_hierarchy; Type: TRIGGER; Schema: storage; Owner: -
--

CREATE TRIGGER "prefixes_create_hierarchy" BEFORE INSERT ON "storage"."prefixes" FOR EACH ROW WHEN (("pg_trigger_depth"() < 1)) EXECUTE FUNCTION "storage"."prefixes_insert_trigger"();


--
-- Name: prefixes prefixes_delete_hierarchy; Type: TRIGGER; Schema: storage; Owner: -
--

CREATE TRIGGER "prefixes_delete_hierarchy" AFTER DELETE ON "storage"."prefixes" FOR EACH ROW EXECUTE FUNCTION "storage"."delete_prefix_hierarchy_trigger"();


--
-- Name: objects update_objects_updated_at; Type: TRIGGER; Schema: storage; Owner: -
--

CREATE TRIGGER "update_objects_updated_at" BEFORE UPDATE ON "storage"."objects" FOR EACH ROW EXECUTE FUNCTION "storage"."update_updated_at_column"();


--
-- Name: extensions extensions_tenant_external_id_fkey; Type: FK CONSTRAINT; Schema: _realtime; Owner: -
--

ALTER TABLE ONLY "_realtime"."extensions"
    ADD CONSTRAINT "extensions_tenant_external_id_fkey" FOREIGN KEY ("tenant_external_id") REFERENCES "_realtime"."tenants"("external_id") ON DELETE CASCADE;


--
-- Name: identities identities_user_id_fkey; Type: FK CONSTRAINT; Schema: auth; Owner: -
--

ALTER TABLE ONLY "auth"."identities"
    ADD CONSTRAINT "identities_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "auth"."users"("id") ON DELETE CASCADE;


--
-- Name: mfa_amr_claims mfa_amr_claims_session_id_fkey; Type: FK CONSTRAINT; Schema: auth; Owner: -
--

ALTER TABLE ONLY "auth"."mfa_amr_claims"
    ADD CONSTRAINT "mfa_amr_claims_session_id_fkey" FOREIGN KEY ("session_id") REFERENCES "auth"."sessions"("id") ON DELETE CASCADE;


--
-- Name: mfa_challenges mfa_challenges_auth_factor_id_fkey; Type: FK CONSTRAINT; Schema: auth; Owner: -
--

ALTER TABLE ONLY "auth"."mfa_challenges"
    ADD CONSTRAINT "mfa_challenges_auth_factor_id_fkey" FOREIGN KEY ("factor_id") REFERENCES "auth"."mfa_factors"("id") ON DELETE CASCADE;


--
-- Name: mfa_factors mfa_factors_user_id_fkey; Type: FK CONSTRAINT; Schema: auth; Owner: -
--

ALTER TABLE ONLY "auth"."mfa_factors"
    ADD CONSTRAINT "mfa_factors_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "auth"."users"("id") ON DELETE CASCADE;


--
-- Name: oauth_authorizations oauth_authorizations_client_id_fkey; Type: FK CONSTRAINT; Schema: auth; Owner: -
--

ALTER TABLE ONLY "auth"."oauth_authorizations"
    ADD CONSTRAINT "oauth_authorizations_client_id_fkey" FOREIGN KEY ("client_id") REFERENCES "auth"."oauth_clients"("id") ON DELETE CASCADE;


--
-- Name: oauth_authorizations oauth_authorizations_user_id_fkey; Type: FK CONSTRAINT; Schema: auth; Owner: -
--

ALTER TABLE ONLY "auth"."oauth_authorizations"
    ADD CONSTRAINT "oauth_authorizations_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "auth"."users"("id") ON DELETE CASCADE;


--
-- Name: oauth_consents oauth_consents_client_id_fkey; Type: FK CONSTRAINT; Schema: auth; Owner: -
--

ALTER TABLE ONLY "auth"."oauth_consents"
    ADD CONSTRAINT "oauth_consents_client_id_fkey" FOREIGN KEY ("client_id") REFERENCES "auth"."oauth_clients"("id") ON DELETE CASCADE;


--
-- Name: oauth_consents oauth_consents_user_id_fkey; Type: FK CONSTRAINT; Schema: auth; Owner: -
--

ALTER TABLE ONLY "auth"."oauth_consents"
    ADD CONSTRAINT "oauth_consents_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "auth"."users"("id") ON DELETE CASCADE;


--
-- Name: one_time_tokens one_time_tokens_user_id_fkey; Type: FK CONSTRAINT; Schema: auth; Owner: -
--

ALTER TABLE ONLY "auth"."one_time_tokens"
    ADD CONSTRAINT "one_time_tokens_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "auth"."users"("id") ON DELETE CASCADE;


--
-- Name: refresh_tokens refresh_tokens_session_id_fkey; Type: FK CONSTRAINT; Schema: auth; Owner: -
--

ALTER TABLE ONLY "auth"."refresh_tokens"
    ADD CONSTRAINT "refresh_tokens_session_id_fkey" FOREIGN KEY ("session_id") REFERENCES "auth"."sessions"("id") ON DELETE CASCADE;


--
-- Name: saml_providers saml_providers_sso_provider_id_fkey; Type: FK CONSTRAINT; Schema: auth; Owner: -
--

ALTER TABLE ONLY "auth"."saml_providers"
    ADD CONSTRAINT "saml_providers_sso_provider_id_fkey" FOREIGN KEY ("sso_provider_id") REFERENCES "auth"."sso_providers"("id") ON DELETE CASCADE;


--
-- Name: saml_relay_states saml_relay_states_flow_state_id_fkey; Type: FK CONSTRAINT; Schema: auth; Owner: -
--

ALTER TABLE ONLY "auth"."saml_relay_states"
    ADD CONSTRAINT "saml_relay_states_flow_state_id_fkey" FOREIGN KEY ("flow_state_id") REFERENCES "auth"."flow_state"("id") ON DELETE CASCADE;


--
-- Name: saml_relay_states saml_relay_states_sso_provider_id_fkey; Type: FK CONSTRAINT; Schema: auth; Owner: -
--

ALTER TABLE ONLY "auth"."saml_relay_states"
    ADD CONSTRAINT "saml_relay_states_sso_provider_id_fkey" FOREIGN KEY ("sso_provider_id") REFERENCES "auth"."sso_providers"("id") ON DELETE CASCADE;


--
-- Name: sessions sessions_oauth_client_id_fkey; Type: FK CONSTRAINT; Schema: auth; Owner: -
--

ALTER TABLE ONLY "auth"."sessions"
    ADD CONSTRAINT "sessions_oauth_client_id_fkey" FOREIGN KEY ("oauth_client_id") REFERENCES "auth"."oauth_clients"("id") ON DELETE CASCADE;


--
-- Name: sessions sessions_user_id_fkey; Type: FK CONSTRAINT; Schema: auth; Owner: -
--

ALTER TABLE ONLY "auth"."sessions"
    ADD CONSTRAINT "sessions_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "auth"."users"("id") ON DELETE CASCADE;


--
-- Name: sso_domains sso_domains_sso_provider_id_fkey; Type: FK CONSTRAINT; Schema: auth; Owner: -
--

ALTER TABLE ONLY "auth"."sso_domains"
    ADD CONSTRAINT "sso_domains_sso_provider_id_fkey" FOREIGN KEY ("sso_provider_id") REFERENCES "auth"."sso_providers"("id") ON DELETE CASCADE;


--
-- Name: audit_logs audit_logs_course_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."audit_logs"
    ADD CONSTRAINT "audit_logs_course_id_fkey" FOREIGN KEY ("course_id") REFERENCES "public"."courses"("id");


--
-- Name: audit_logs audit_logs_lecture_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."audit_logs"
    ADD CONSTRAINT "audit_logs_lecture_id_fkey" FOREIGN KEY ("lecture_id") REFERENCES "public"."lectures"("id");


--
-- Name: audit_logs audit_logs_transaction_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."audit_logs"
    ADD CONSTRAINT "audit_logs_transaction_id_fkey" FOREIGN KEY ("transaction_id") REFERENCES "public"."transactions"("id");


--
-- Name: audit_logs audit_logs_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."audit_logs"
    ADD CONSTRAINT "audit_logs_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."users"("id");


--
-- Name: chat_history chat_history_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."chat_history"
    ADD CONSTRAINT "chat_history_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."users"("id") ON DELETE CASCADE;


--
-- Name: course_access_logs course_access_logs_course_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."course_access_logs"
    ADD CONSTRAINT "course_access_logs_course_id_fkey" FOREIGN KEY ("course_id") REFERENCES "public"."courses"("id");


--
-- Name: course_access_logs course_access_logs_lecture_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."course_access_logs"
    ADD CONSTRAINT "course_access_logs_lecture_id_fkey" FOREIGN KEY ("lecture_id") REFERENCES "public"."lectures"("id");


--
-- Name: course_access_logs course_access_logs_transaction_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."course_access_logs"
    ADD CONSTRAINT "course_access_logs_transaction_id_fkey" FOREIGN KEY ("transaction_id") REFERENCES "public"."transactions"("id");


--
-- Name: course_access_logs course_access_logs_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."course_access_logs"
    ADD CONSTRAINT "course_access_logs_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."users"("id");


--
-- Name: course_resources course_resources_course_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."course_resources"
    ADD CONSTRAINT "course_resources_course_id_fkey" FOREIGN KEY ("course_id") REFERENCES "public"."courses"("id") ON DELETE CASCADE;


--
-- Name: course_resources course_resources_file_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."course_resources"
    ADD CONSTRAINT "course_resources_file_id_fkey" FOREIGN KEY ("file_id") REFERENCES "public"."files"("id") ON DELETE CASCADE;


--
-- Name: enrollments enrollments_course_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."enrollments"
    ADD CONSTRAINT "enrollments_course_id_fkey" FOREIGN KEY ("course_id") REFERENCES "public"."courses"("id") ON DELETE CASCADE;


--
-- Name: enrollments enrollments_transaction_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."enrollments"
    ADD CONSTRAINT "enrollments_transaction_id_fkey" FOREIGN KEY ("transaction_id") REFERENCES "public"."transactions"("id");


--
-- Name: file_permissions file_permissions_file_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."file_permissions"
    ADD CONSTRAINT "file_permissions_file_id_fkey" FOREIGN KEY ("file_id") REFERENCES "public"."files"("id") ON DELETE CASCADE;


--
-- Name: courses fk_courses_instructor_id; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."courses"
    ADD CONSTRAINT "fk_courses_instructor_id" FOREIGN KEY ("instructor_id") REFERENCES "public"."users"("id");


--
-- Name: enrollments fk_enrollments_user_id; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."enrollments"
    ADD CONSTRAINT "fk_enrollments_user_id" FOREIGN KEY ("user_id") REFERENCES "public"."users"("id");


--
-- Name: lecture_resources fk_lecture_resources_file_id; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."lecture_resources"
    ADD CONSTRAINT "fk_lecture_resources_file_id" FOREIGN KEY ("file_id") REFERENCES "public"."files"("id") ON DELETE CASCADE;


--
-- Name: lecture_resources fk_lecture_resources_lecture_id; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."lecture_resources"
    ADD CONSTRAINT "fk_lecture_resources_lecture_id" FOREIGN KEY ("lecture_id") REFERENCES "public"."lectures"("id") ON DELETE CASCADE;


--
-- Name: lemon_squeezy_variants fk_lemon_squeezy_variants_product; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."lemon_squeezy_variants"
    ADD CONSTRAINT "fk_lemon_squeezy_variants_product" FOREIGN KEY ("lemon_squeezy_product_id") REFERENCES "public"."lemon_squeezy_products"("lemon_squeezy_product_id") ON DELETE CASCADE;


--
-- Name: notes fk_notes_course; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."notes"
    ADD CONSTRAINT "fk_notes_course" FOREIGN KEY ("course_id") REFERENCES "public"."courses"("id") ON DELETE CASCADE;


--
-- Name: notes fk_notes_lecture; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."notes"
    ADD CONSTRAINT "fk_notes_lecture" FOREIGN KEY ("lecture_id") REFERENCES "public"."lectures"("id") ON DELETE CASCADE;


--
-- Name: notes fk_notes_user; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."notes"
    ADD CONSTRAINT "fk_notes_user" FOREIGN KEY ("user_id") REFERENCES "public"."users"("id") ON DELETE CASCADE;


--
-- Name: progress fk_progress_course_id; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."progress"
    ADD CONSTRAINT "fk_progress_course_id" FOREIGN KEY ("course_id") REFERENCES "public"."courses"("id") ON DELETE CASCADE;


--
-- Name: forum_mentions forum_mentions_mentioned_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."forum_mentions"
    ADD CONSTRAINT "forum_mentions_mentioned_user_id_fkey" FOREIGN KEY ("mentioned_user_id") REFERENCES "public"."users"("id") ON DELETE CASCADE;


--
-- Name: forum_mentions forum_mentions_mentioner_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."forum_mentions"
    ADD CONSTRAINT "forum_mentions_mentioner_user_id_fkey" FOREIGN KEY ("mentioner_user_id") REFERENCES "public"."users"("id") ON DELETE CASCADE;


--
-- Name: forum_mentions forum_mentions_post_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."forum_mentions"
    ADD CONSTRAINT "forum_mentions_post_id_fkey" FOREIGN KEY ("post_id") REFERENCES "public"."forum_posts"("id") ON DELETE CASCADE;


--
-- Name: forum_notifications forum_notifications_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."forum_notifications"
    ADD CONSTRAINT "forum_notifications_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."users"("id") ON DELETE CASCADE;


--
-- Name: forum_posts forum_posts_author_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."forum_posts"
    ADD CONSTRAINT "forum_posts_author_id_fkey" FOREIGN KEY ("author_id") REFERENCES "public"."users"("id") ON DELETE CASCADE;


--
-- Name: forum_posts forum_posts_parent_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."forum_posts"
    ADD CONSTRAINT "forum_posts_parent_id_fkey" FOREIGN KEY ("parent_id") REFERENCES "public"."forum_posts"("id") ON DELETE CASCADE;


--
-- Name: forum_posts forum_posts_topic_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."forum_posts"
    ADD CONSTRAINT "forum_posts_topic_id_fkey" FOREIGN KEY ("topic_id") REFERENCES "public"."forum_topics"("id") ON DELETE CASCADE;


--
-- Name: forum_topic_subscriptions forum_topic_subscriptions_topic_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."forum_topic_subscriptions"
    ADD CONSTRAINT "forum_topic_subscriptions_topic_id_fkey" FOREIGN KEY ("topic_id") REFERENCES "public"."forum_topics"("id") ON DELETE CASCADE;


--
-- Name: forum_topic_subscriptions forum_topic_subscriptions_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."forum_topic_subscriptions"
    ADD CONSTRAINT "forum_topic_subscriptions_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."users"("id") ON DELETE CASCADE;


--
-- Name: forum_topics forum_topics_creator_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."forum_topics"
    ADD CONSTRAINT "forum_topics_creator_id_fkey" FOREIGN KEY ("creator_id") REFERENCES "public"."users"("id") ON DELETE CASCADE;


--
-- Name: forum_votes forum_votes_post_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."forum_votes"
    ADD CONSTRAINT "forum_votes_post_id_fkey" FOREIGN KEY ("post_id") REFERENCES "public"."forum_posts"("id") ON DELETE CASCADE;


--
-- Name: forum_votes forum_votes_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."forum_votes"
    ADD CONSTRAINT "forum_votes_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."users"("id") ON DELETE CASCADE;


--
-- Name: lecture_preview_sessions lecture_preview_sessions_lecture_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."lecture_preview_sessions"
    ADD CONSTRAINT "lecture_preview_sessions_lecture_id_fkey" FOREIGN KEY ("lecture_id") REFERENCES "public"."lectures"("id");


--
-- Name: lectures lectures_course_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."lectures"
    ADD CONSTRAINT "lectures_course_id_fkey" FOREIGN KEY ("course_id") REFERENCES "public"."courses"("id") ON DELETE CASCADE;


--
-- Name: oauth_accounts oauth_accounts_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."oauth_accounts"
    ADD CONSTRAINT "oauth_accounts_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."users"("id") ON DELETE CASCADE;


--
-- Name: payment_events payment_events_course_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."payment_events"
    ADD CONSTRAINT "payment_events_course_id_fkey" FOREIGN KEY ("course_id") REFERENCES "public"."courses"("id");


--
-- Name: payment_events payment_events_transaction_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."payment_events"
    ADD CONSTRAINT "payment_events_transaction_id_fkey" FOREIGN KEY ("transaction_id") REFERENCES "public"."transactions"("id");


--
-- Name: payment_events payment_events_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."payment_events"
    ADD CONSTRAINT "payment_events_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."users"("id");


--
-- Name: payment_methods payment_methods_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."payment_methods"
    ADD CONSTRAINT "payment_methods_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."users"("id") ON DELETE CASCADE;


--
-- Name: progress progress_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."progress"
    ADD CONSTRAINT "progress_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."users"("id") ON DELETE CASCADE;


--
-- Name: stripe_customers stripe_customers_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."stripe_customers"
    ADD CONSTRAINT "stripe_customers_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."users"("id") ON DELETE CASCADE;


--
-- Name: stripe_products stripe_products_course_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."stripe_products"
    ADD CONSTRAINT "stripe_products_course_id_fkey" FOREIGN KEY ("course_id") REFERENCES "public"."courses"("id") ON DELETE CASCADE;


--
-- Name: subscriptions subscriptions_payment_method_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."subscriptions"
    ADD CONSTRAINT "subscriptions_payment_method_id_fkey" FOREIGN KEY ("payment_method_id") REFERENCES "public"."payment_methods"("id");


--
-- Name: subscriptions subscriptions_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."subscriptions"
    ADD CONSTRAINT "subscriptions_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."users"("id") ON DELETE CASCADE;


--
-- Name: transactions transactions_payment_method_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."transactions"
    ADD CONSTRAINT "transactions_payment_method_id_fkey" FOREIGN KEY ("payment_method_id") REFERENCES "public"."payment_methods"("id");


--
-- Name: transactions transactions_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."transactions"
    ADD CONSTRAINT "transactions_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."users"("id") ON DELETE CASCADE;


--
-- Name: user_payment_methods user_payment_methods_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."user_payment_methods"
    ADD CONSTRAINT "user_payment_methods_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."users"("id");


--
-- Name: video_analytics video_analytics_video_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."video_analytics"
    ADD CONSTRAINT "video_analytics_video_id_fkey" FOREIGN KEY ("video_id") REFERENCES "public"."videos"("id") ON DELETE CASCADE;


--
-- Name: video_permissions video_permissions_video_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."video_permissions"
    ADD CONSTRAINT "video_permissions_video_id_fkey" FOREIGN KEY ("video_id") REFERENCES "public"."videos"("id") ON DELETE CASCADE;


--
-- Name: video_qualities video_qualities_video_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."video_qualities"
    ADD CONSTRAINT "video_qualities_video_id_fkey" FOREIGN KEY ("video_id") REFERENCES "public"."videos"("id") ON DELETE CASCADE;


--
-- Name: viewing_sessions viewing_sessions_video_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY "public"."viewing_sessions"
    ADD CONSTRAINT "viewing_sessions_video_id_fkey" FOREIGN KEY ("video_id") REFERENCES "public"."videos"("id") ON DELETE CASCADE;


--
-- Name: iceberg_namespaces iceberg_namespaces_bucket_id_fkey; Type: FK CONSTRAINT; Schema: storage; Owner: -
--

ALTER TABLE ONLY "storage"."iceberg_namespaces"
    ADD CONSTRAINT "iceberg_namespaces_bucket_id_fkey" FOREIGN KEY ("bucket_id") REFERENCES "storage"."buckets_analytics"("id") ON DELETE CASCADE;


--
-- Name: iceberg_tables iceberg_tables_bucket_id_fkey; Type: FK CONSTRAINT; Schema: storage; Owner: -
--

ALTER TABLE ONLY "storage"."iceberg_tables"
    ADD CONSTRAINT "iceberg_tables_bucket_id_fkey" FOREIGN KEY ("bucket_id") REFERENCES "storage"."buckets_analytics"("id") ON DELETE CASCADE;


--
-- Name: iceberg_tables iceberg_tables_namespace_id_fkey; Type: FK CONSTRAINT; Schema: storage; Owner: -
--

ALTER TABLE ONLY "storage"."iceberg_tables"
    ADD CONSTRAINT "iceberg_tables_namespace_id_fkey" FOREIGN KEY ("namespace_id") REFERENCES "storage"."iceberg_namespaces"("id") ON DELETE CASCADE;


--
-- Name: objects objects_bucketId_fkey; Type: FK CONSTRAINT; Schema: storage; Owner: -
--

ALTER TABLE ONLY "storage"."objects"
    ADD CONSTRAINT "objects_bucketId_fkey" FOREIGN KEY ("bucket_id") REFERENCES "storage"."buckets"("id");


--
-- Name: prefixes prefixes_bucketId_fkey; Type: FK CONSTRAINT; Schema: storage; Owner: -
--

ALTER TABLE ONLY "storage"."prefixes"
    ADD CONSTRAINT "prefixes_bucketId_fkey" FOREIGN KEY ("bucket_id") REFERENCES "storage"."buckets"("id");


--
-- Name: s3_multipart_uploads s3_multipart_uploads_bucket_id_fkey; Type: FK CONSTRAINT; Schema: storage; Owner: -
--

ALTER TABLE ONLY "storage"."s3_multipart_uploads"
    ADD CONSTRAINT "s3_multipart_uploads_bucket_id_fkey" FOREIGN KEY ("bucket_id") REFERENCES "storage"."buckets"("id");


--
-- Name: s3_multipart_uploads_parts s3_multipart_uploads_parts_bucket_id_fkey; Type: FK CONSTRAINT; Schema: storage; Owner: -
--

ALTER TABLE ONLY "storage"."s3_multipart_uploads_parts"
    ADD CONSTRAINT "s3_multipart_uploads_parts_bucket_id_fkey" FOREIGN KEY ("bucket_id") REFERENCES "storage"."buckets"("id");


--
-- Name: s3_multipart_uploads_parts s3_multipart_uploads_parts_upload_id_fkey; Type: FK CONSTRAINT; Schema: storage; Owner: -
--

ALTER TABLE ONLY "storage"."s3_multipart_uploads_parts"
    ADD CONSTRAINT "s3_multipart_uploads_parts_upload_id_fkey" FOREIGN KEY ("upload_id") REFERENCES "storage"."s3_multipart_uploads"("id") ON DELETE CASCADE;


--
-- Name: audit_log_entries; Type: ROW SECURITY; Schema: auth; Owner: -
--

ALTER TABLE "auth"."audit_log_entries" ENABLE ROW LEVEL SECURITY;

--
-- Name: flow_state; Type: ROW SECURITY; Schema: auth; Owner: -
--

ALTER TABLE "auth"."flow_state" ENABLE ROW LEVEL SECURITY;

--
-- Name: identities; Type: ROW SECURITY; Schema: auth; Owner: -
--

ALTER TABLE "auth"."identities" ENABLE ROW LEVEL SECURITY;

--
-- Name: instances; Type: ROW SECURITY; Schema: auth; Owner: -
--

ALTER TABLE "auth"."instances" ENABLE ROW LEVEL SECURITY;

--
-- Name: mfa_amr_claims; Type: ROW SECURITY; Schema: auth; Owner: -
--

ALTER TABLE "auth"."mfa_amr_claims" ENABLE ROW LEVEL SECURITY;

--
-- Name: mfa_challenges; Type: ROW SECURITY; Schema: auth; Owner: -
--

ALTER TABLE "auth"."mfa_challenges" ENABLE ROW LEVEL SECURITY;

--
-- Name: mfa_factors; Type: ROW SECURITY; Schema: auth; Owner: -
--

ALTER TABLE "auth"."mfa_factors" ENABLE ROW LEVEL SECURITY;

--
-- Name: one_time_tokens; Type: ROW SECURITY; Schema: auth; Owner: -
--

ALTER TABLE "auth"."one_time_tokens" ENABLE ROW LEVEL SECURITY;

--
-- Name: refresh_tokens; Type: ROW SECURITY; Schema: auth; Owner: -
--

ALTER TABLE "auth"."refresh_tokens" ENABLE ROW LEVEL SECURITY;

--
-- Name: saml_providers; Type: ROW SECURITY; Schema: auth; Owner: -
--

ALTER TABLE "auth"."saml_providers" ENABLE ROW LEVEL SECURITY;

--
-- Name: saml_relay_states; Type: ROW SECURITY; Schema: auth; Owner: -
--

ALTER TABLE "auth"."saml_relay_states" ENABLE ROW LEVEL SECURITY;

--
-- Name: schema_migrations; Type: ROW SECURITY; Schema: auth; Owner: -
--

ALTER TABLE "auth"."schema_migrations" ENABLE ROW LEVEL SECURITY;

--
-- Name: sessions; Type: ROW SECURITY; Schema: auth; Owner: -
--

ALTER TABLE "auth"."sessions" ENABLE ROW LEVEL SECURITY;

--
-- Name: sso_domains; Type: ROW SECURITY; Schema: auth; Owner: -
--

ALTER TABLE "auth"."sso_domains" ENABLE ROW LEVEL SECURITY;

--
-- Name: sso_providers; Type: ROW SECURITY; Schema: auth; Owner: -
--

ALTER TABLE "auth"."sso_providers" ENABLE ROW LEVEL SECURITY;

--
-- Name: users; Type: ROW SECURITY; Schema: auth; Owner: -
--

ALTER TABLE "auth"."users" ENABLE ROW LEVEL SECURITY;

--
-- Name: messages; Type: ROW SECURITY; Schema: realtime; Owner: -
--

ALTER TABLE "realtime"."messages" ENABLE ROW LEVEL SECURITY;

--
-- Name: buckets; Type: ROW SECURITY; Schema: storage; Owner: -
--

ALTER TABLE "storage"."buckets" ENABLE ROW LEVEL SECURITY;

--
-- Name: buckets_analytics; Type: ROW SECURITY; Schema: storage; Owner: -
--

ALTER TABLE "storage"."buckets_analytics" ENABLE ROW LEVEL SECURITY;

--
-- Name: iceberg_namespaces; Type: ROW SECURITY; Schema: storage; Owner: -
--

ALTER TABLE "storage"."iceberg_namespaces" ENABLE ROW LEVEL SECURITY;

--
-- Name: iceberg_tables; Type: ROW SECURITY; Schema: storage; Owner: -
--

ALTER TABLE "storage"."iceberg_tables" ENABLE ROW LEVEL SECURITY;

--
-- Name: migrations; Type: ROW SECURITY; Schema: storage; Owner: -
--

ALTER TABLE "storage"."migrations" ENABLE ROW LEVEL SECURITY;

--
-- Name: objects; Type: ROW SECURITY; Schema: storage; Owner: -
--

ALTER TABLE "storage"."objects" ENABLE ROW LEVEL SECURITY;

--
-- Name: prefixes; Type: ROW SECURITY; Schema: storage; Owner: -
--

ALTER TABLE "storage"."prefixes" ENABLE ROW LEVEL SECURITY;

--
-- Name: s3_multipart_uploads; Type: ROW SECURITY; Schema: storage; Owner: -
--

ALTER TABLE "storage"."s3_multipart_uploads" ENABLE ROW LEVEL SECURITY;

--
-- Name: s3_multipart_uploads_parts; Type: ROW SECURITY; Schema: storage; Owner: -
--

ALTER TABLE "storage"."s3_multipart_uploads_parts" ENABLE ROW LEVEL SECURITY;

--
-- Name: supabase_realtime; Type: PUBLICATION; Schema: -; Owner: -
--

CREATE PUBLICATION "supabase_realtime" WITH (publish = 'insert, update, delete, truncate');


--
-- Name: issue_graphql_placeholder; Type: EVENT TRIGGER; Schema: -; Owner: -
--

CREATE EVENT TRIGGER "issue_graphql_placeholder" ON "sql_drop"
         WHEN TAG IN ('DROP EXTENSION')
   EXECUTE FUNCTION "extensions"."set_graphql_placeholder"();


--
-- Name: issue_pg_cron_access; Type: EVENT TRIGGER; Schema: -; Owner: -
--

CREATE EVENT TRIGGER "issue_pg_cron_access" ON "ddl_command_end"
         WHEN TAG IN ('CREATE EXTENSION')
   EXECUTE FUNCTION "extensions"."grant_pg_cron_access"();


--
-- Name: issue_pg_graphql_access; Type: EVENT TRIGGER; Schema: -; Owner: -
--

CREATE EVENT TRIGGER "issue_pg_graphql_access" ON "ddl_command_end"
         WHEN TAG IN ('CREATE FUNCTION')
   EXECUTE FUNCTION "extensions"."grant_pg_graphql_access"();


--
-- Name: issue_pg_net_access; Type: EVENT TRIGGER; Schema: -; Owner: -
--

CREATE EVENT TRIGGER "issue_pg_net_access" ON "ddl_command_end"
         WHEN TAG IN ('CREATE EXTENSION')
   EXECUTE FUNCTION "extensions"."grant_pg_net_access"();


--
-- Name: pgrst_ddl_watch; Type: EVENT TRIGGER; Schema: -; Owner: -
--

CREATE EVENT TRIGGER "pgrst_ddl_watch" ON "ddl_command_end"
   EXECUTE FUNCTION "extensions"."pgrst_ddl_watch"();


--
-- Name: pgrst_drop_watch; Type: EVENT TRIGGER; Schema: -; Owner: -
--

CREATE EVENT TRIGGER "pgrst_drop_watch" ON "sql_drop"
   EXECUTE FUNCTION "extensions"."pgrst_drop_watch"();


--
-- PostgreSQL database dump complete
--

