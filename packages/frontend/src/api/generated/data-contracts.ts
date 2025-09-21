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

export enum DomainUserRole {
  UserRoleEmployee = "EMPLOYEE",
  UserRoleAdmin = "ADMIN",
  UserRoleManager = "MANAGER",
}

export enum DomainToolStatus {
  ToolStatusInOffice = "IN_OFFICE",
  ToolStatusCheckedOut = "CHECKED_OUT",
  ToolStatusMaintenance = "MAINTENANCE",
  ToolStatusLost = "LOST",
}

export enum DomainEventType {
  EventTypeToolCreated = "TOOL_CREATED",
  EventTypeToolUpdated = "TOOL_UPDATED",
  EventTypeToolDeleted = "TOOL_DELETED",
  EventTypeToolCheckedOut = "TOOL_CHECKED_OUT",
  EventTypeToolCheckedIn = "TOOL_CHECKED_IN",
  EventTypeToolMaintenance = "TOOL_MAINTENANCE",
  EventTypeToolLost = "TOOL_LOST",
  EventTypeUserCreated = "USER_CREATED",
  EventTypeUserUpdated = "USER_UPDATED",
  EventTypeUserDeleted = "USER_DELETED",
}

export interface DomainEvent {
  actor_id?: string;
  created_at?: string;
  id?: string;
  metadata?: string;
  notes?: string;
  tool_id?: string;
  type?: DomainEventType;
  user_id?: string;
}

export interface DomainTool {
  created_at?: string;
  current_user_id?: string;
  id?: string;
  last_checked_out_at?: string;
  name?: string;
  status?: DomainToolStatus;
  updated_at?: string;
}

export interface DomainUser {
  created_at?: string;
  email?: string;
  id?: string;
  name?: string;
  role?: DomainUserRole;
  updated_at?: string;
}

export interface ServerCheckinToolRequest {
  notes?: string;
  user_id: string;
}

export interface ServerCheckoutToolRequest {
  notes?: string;
  user_id: string;
}

export interface ServerCreateToolRequest {
  name: string;
  status?: DomainToolStatus;
}

export interface ServerCreateUserRequest {
  email: string;
  name: string;
  role?: DomainUserRole;
}

export interface ServerMaintenanceRequest {
  notes?: string;
  user_id: string;
}

export interface ServerMarkLostRequest {
  notes?: string;
  user_id: string;
}

export interface ServerStatsResponse {
  tools_by_status?: {
    checked_out?: number;
    in_office?: number;
    lost?: number;
    maintenance?: number;
  };
  total_events?: number;
  total_tools?: number;
  total_users?: number;
  users_by_role?: {
    admins?: number;
    employees?: number;
    managers?: number;
  };
}

export interface ServerUpdateToolRequest {
  name: string;
  status?: DomainToolStatus;
}

export interface ServerUpdateUserRequest {
  email: string;
  name: string;
  role: DomainUserRole;
}
