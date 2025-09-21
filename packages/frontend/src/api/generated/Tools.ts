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

export class Tools<SecurityDataType = unknown> extends HttpClient<SecurityDataType>  {

            /**
 * @description Get a list of tools with pagination and optional status filtering
 *
 * @tags tools
 * @name ToolsList
 * @summary List all tools
 * @request GET:/tools
 */
toolsList: (query?: {
  /**
   * Limit
   * @default 10
   */
    limit?: number,
  /**
   * Offset
   * @default 0
   */
    offset?: number,
  /** Filter by status */
    status?: string,

}, params: RequestParams = {}) =>
    this.request<Record<string,(DomainTool)[]>, Record<string,string>>({
        path: `/tools`,
        method: 'GET',
        query: query,                        type: ContentType.Json,        format: "json",        ...params,
    }),            /**
 * @description Create a new tool with name and status
 *
 * @tags tools
 * @name ToolsCreate
 * @summary Create a new tool
 * @request POST:/tools
 */
toolsCreate: (tool: ServerCreateToolRequest, params: RequestParams = {}) =>
    this.request<DomainTool, Record<string,string>>({
        path: `/tools`,
        method: 'POST',
                body: tool,                type: ContentType.Json,        format: "json",        ...params,
    }),            /**
 * @description Get a specific tool by its ID
 *
 * @tags tools
 * @name ToolsDetail
 * @summary Get a tool by ID
 * @request GET:/tools/{id}
 */
toolsDetail: (id: string, params: RequestParams = {}) =>
    this.request<DomainTool, Record<string,string>>({
        path: `/tools/${id}`,
        method: 'GET',
                                type: ContentType.Json,        format: "json",        ...params,
    }),            /**
 * @description Update a tool's name and status
 *
 * @tags tools
 * @name ToolsUpdate
 * @summary Update a tool
 * @request PUT:/tools/{id}
 */
toolsUpdate: (id: string, tool: ServerUpdateToolRequest, params: RequestParams = {}) =>
    this.request<DomainTool, Record<string,string>>({
        path: `/tools/${id}`,
        method: 'PUT',
                body: tool,                type: ContentType.Json,        format: "json",        ...params,
    }),            /**
 * @description Delete a tool from the system
 *
 * @tags tools
 * @name ToolsDelete
 * @summary Delete a tool
 * @request DELETE:/tools/{id}
 */
toolsDelete: (id: string, params: RequestParams = {}) =>
    this.request<void, Record<string,string>>({
        path: `/tools/${id}`,
        method: 'DELETE',
                                type: ContentType.Json,                ...params,
    }),            /**
 * @description Check in a tool that was previously checked out
 *
 * @tags tools
 * @name CheckinCreate
 * @summary Check in a tool from a user
 * @request POST:/tools/{id}/checkin
 */
checkinCreate: (id: string, checkin: ServerCheckinToolRequest, params: RequestParams = {}) =>
    this.request<Record<string,any>, Record<string,string>>({
        path: `/tools/${id}/checkin`,
        method: 'POST',
                body: checkin,                type: ContentType.Json,        format: "json",        ...params,
    }),            /**
 * @description Check out a tool to a specific user with optional notes
 *
 * @tags tools
 * @name CheckoutCreate
 * @summary Check out a tool to a user
 * @request POST:/tools/{id}/checkout
 */
checkoutCreate: (id: string, checkout: ServerCheckoutToolRequest, params: RequestParams = {}) =>
    this.request<Record<string,any>, Record<string,string>>({
        path: `/tools/${id}/checkout`,
        method: 'POST',
                body: checkout,                type: ContentType.Json,        format: "json",        ...params,
    }),            /**
 * @description Get the complete event history for a specific tool
 *
 * @tags tools
 * @name HistoryList
 * @summary Get tool history
 * @request GET:/tools/{id}/history
 */
historyList: (id: string, params: RequestParams = {}) =>
    this.request<Record<string,(DomainEvent)[]>, Record<string,string>>({
        path: `/tools/${id}/history`,
        method: 'GET',
                                type: ContentType.Json,        format: "json",        ...params,
    }),            /**
 * @description Mark a tool as lost or missing
 *
 * @tags tools
 * @name LostCreate
 * @summary Mark a tool as lost
 * @request POST:/tools/{id}/lost
 */
lostCreate: (id: string, lost: ServerMarkLostRequest, params: RequestParams = {}) =>
    this.request<Record<string,any>, Record<string,string>>({
        path: `/tools/${id}/lost`,
        method: 'POST',
                body: lost,                type: ContentType.Json,        format: "json",        ...params,
    }),            /**
 * @description Mark a tool as being in maintenance
 *
 * @tags tools
 * @name MaintenanceCreate
 * @summary Send a tool to maintenance
 * @request POST:/tools/{id}/maintenance
 */
maintenanceCreate: (id: string, maintenance: ServerMaintenanceRequest, params: RequestParams = {}) =>
    this.request<Record<string,any>, Record<string,string>>({
        path: `/tools/${id}/maintenance`,
        method: 'POST',
                body: maintenance,                type: ContentType.Json,        format: "json",        ...params,
    }),    }
