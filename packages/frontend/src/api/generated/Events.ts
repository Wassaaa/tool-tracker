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

export class Events<SecurityDataType = unknown> extends HttpClient<SecurityDataType>  {

            /**
 * @description Get a list of events with pagination and optional filtering
 *
 * @tags events
 * @name EventsList
 * @summary List all events
 * @request GET:/events
 */
eventsList: (query?: {
  /**
   * Limit
   * @default 50
   */
    limit?: number,
  /**
   * Offset
   * @default 0
   */
    offset?: number,
  /** Filter by event type */
    type?: string,
  /** Filter by tool ID */
    tool_id?: string,
  /** Filter by user ID */
    user_id?: string,

}, params: RequestParams = {}) =>
    this.request<Record<string,(DomainEvent)[]>, Record<string,string>>({
        path: `/events`,
        method: 'GET',
        query: query,                        type: ContentType.Json,        format: "json",        ...params,
    }),            /**
 * @description Get a specific event by its ID
 *
 * @tags events
 * @name EventsDetail
 * @summary Get an event by ID
 * @request GET:/events/{id}
 */
eventsDetail: (id: string, params: RequestParams = {}) =>
    this.request<DomainEvent, Record<string,string>>({
        path: `/events/${id}`,
        method: 'GET',
                                type: ContentType.Json,        format: "json",        ...params,
    }),    }
