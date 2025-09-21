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

export class Users<SecurityDataType = unknown> extends HttpClient<SecurityDataType>  {

            /**
 * @description Get a list of users with pagination and optional role filtering
 *
 * @tags users
 * @name UsersList
 * @summary List all users
 * @request GET:/users
 */
usersList: (query?: {
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
  /** Filter by role (EMPLOYEE, ADMIN, MANAGER) */
    role?: string,

}, params: RequestParams = {}) =>
    this.request<Record<string,(DomainUser)[]>, Record<string,string>>({
        path: `/users`,
        method: 'GET',
        query: query,                        type: ContentType.Json,        format: "json",        ...params,
    }),            /**
 * @description Create a new user with name, email, and role
 *
 * @tags users
 * @name UsersCreate
 * @summary Create a new user
 * @request POST:/users
 */
usersCreate: (user: ServerCreateUserRequest, params: RequestParams = {}) =>
    this.request<DomainUser, Record<string,string>>({
        path: `/users`,
        method: 'POST',
                body: user,                type: ContentType.Json,        format: "json",        ...params,
    }),            /**
 * @description Get a specific user by their ID
 *
 * @tags users
 * @name UsersDetail
 * @summary Get a user by ID
 * @request GET:/users/{id}
 */
usersDetail: (id: string, params: RequestParams = {}) =>
    this.request<DomainUser, Record<string,string>>({
        path: `/users/${id}`,
        method: 'GET',
                                type: ContentType.Json,        format: "json",        ...params,
    }),            /**
 * @description Update a user's name, email, and role
 *
 * @tags users
 * @name UsersUpdate
 * @summary Update a user
 * @request PUT:/users/{id}
 */
usersUpdate: (id: string, user: ServerUpdateUserRequest, params: RequestParams = {}) =>
    this.request<DomainUser, Record<string,string>>({
        path: `/users/${id}`,
        method: 'PUT',
                body: user,                type: ContentType.Json,        format: "json",        ...params,
    }),            /**
 * @description Delete a user from the system
 *
 * @tags users
 * @name UsersDelete
 * @summary Delete a user
 * @request DELETE:/users/{id}
 */
usersDelete: (id: string, params: RequestParams = {}) =>
    this.request<void, Record<string,string>>({
        path: `/users/${id}`,
        method: 'DELETE',
                                type: ContentType.Json,                ...params,
    }),            /**
 * @description Get the complete activity history for a specific user
 *
 * @tags users
 * @name ActivityList
 * @summary Get user activity
 * @request GET:/users/{id}/activity
 */
activityList: (id: string, params: RequestParams = {}) =>
    this.request<Record<string,(DomainEvent)[]>, Record<string,string>>({
        path: `/users/${id}/activity`,
        method: 'GET',
                                type: ContentType.Json,        format: "json",        ...params,
    }),            /**
 * @description Get list of tools currently checked out by a specific user
 *
 * @tags users
 * @name ToolsList
 * @summary Get tools assigned to user
 * @request GET:/users/{id}/tools
 */
toolsList: (id: string, params: RequestParams = {}) =>
    this.request<Record<string,(string)[]>, Record<string,string>>({
        path: `/users/${id}/tools`,
        method: 'GET',
                                type: ContentType.Json,        format: "json",        ...params,
    }),    }
