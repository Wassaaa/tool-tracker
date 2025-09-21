/* eslint-disable */
/* tslint:disable */
// @ts-nocheck
/*
 * ---------------------------------------------------------------
 * ## THIS FILE WAS GENERATED VIA SWAGGER-TYPESCRIPT-API        ##
 * ##                                                           ##
 * ## AUTHOR: acacode                                           ##
 * ## SOURCE: https://github.com/acacode/swagger-typescript-api ##
 * ---------------------------------------------------------------
 */



import { HttpClient, RequestParams, ContentType, HttpResponse } from "./http-client";
import { DomainUserRole, DomainToolStatus, DomainEventType, DomainEvent, DomainTool, DomainUser, ServerCheckinToolRequest, ServerCheckoutToolRequest, ServerCreateToolRequest, ServerCreateUserRequest, ServerMaintenanceRequest, ServerMarkLostRequest, ServerStatsResponse, ServerUpdateToolRequest, ServerUpdateUserRequest } from "./data-contracts"

export class Admin<SecurityDataType = unknown> extends HttpClient<SecurityDataType>  {

            /**
 * @description Get recent audit events for administrative review
 *
 * @tags admin
 * @name AuditList
 * @summary Get audit log
 * @request GET:/admin/audit
 */
auditList: (params: RequestParams = {}) =>
    this.request<Record<string,any>, Record<string,string>>({
        path: `/admin/audit`,
        method: 'GET',
                                type: ContentType.Json,        format: "json",        ...params,
    }),            /**
 * @description Get comprehensive statistics about tools, users, and events
 *
 * @tags admin
 * @name StatsList
 * @summary Get system statistics
 * @request GET:/admin/stats
 */
statsList: (params: RequestParams = {}) =>
    this.request<ServerStatsResponse, Record<string,string>>({
        path: `/admin/stats`,
        method: 'GET',
                                type: ContentType.Json,        format: "json",        ...params,
    }),    }
