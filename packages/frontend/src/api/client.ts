// API Client wrapper using the generated code
import { Tools } from './generated/Tools';
import { Users } from './generated/Users';
import { Events } from './generated/Events';
import { Admin } from './generated/Admin';
import type { 
  DomainTool, 
  DomainUser, 
  DomainEvent,
  DomainToolStatus,
  DomainUserRole,
  ServerCreateToolRequest,
  ServerCreateUserRequest,
  ServerUpdateToolRequest,
  ServerUpdateUserRequest,
  ServerCheckoutToolRequest,
  ServerCheckinToolRequest
} from './generated/data-contracts';

// Base API configuration
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api';

// Create API client instances
export const toolsApi = new Tools({
  baseUrl: API_BASE_URL,
});

export const usersApi = new Users({
  baseUrl: API_BASE_URL,
});

export const eventsApi = new Events({
  baseUrl: API_BASE_URL,
});

export const adminApi = new Admin({
  baseUrl: API_BASE_URL,
});

// Re-export types for convenience
export type {
  DomainTool,
  DomainUser,
  DomainEvent,
  DomainToolStatus,
  DomainUserRole,
  ServerCreateToolRequest,
  ServerCreateUserRequest,
  ServerUpdateToolRequest,
  ServerUpdateUserRequest,
  ServerCheckoutToolRequest,
  ServerCheckinToolRequest,
};

// Higher-level API functions with better error handling
export class ToolTrackerAPI {
  
  // Tools API
  static async getTools(params?: { limit?: number; offset?: number; status?: string }) {
    try {
      const response = await toolsApi.toolsList(params);
      return response.data;
    } catch (error) {
      console.error('Failed to fetch tools:', error);
      throw new Error('Failed to fetch tools');
    }
  }

  static async getTool(id: string) {
    try {
      const response = await toolsApi.toolsDetail(id);
      return response.data;
    } catch (error) {
      console.error(`Failed to fetch tool ${id}:`, error);
      throw new Error(`Failed to fetch tool`);
    }
  }

  static async createTool(data: ServerCreateToolRequest) {
    try {
      const response = await toolsApi.toolsCreate(data);
      return response.data;
    } catch (error) {
      console.error('Failed to create tool:', error);
      throw new Error('Failed to create tool');
    }
  }

  static async updateTool(id: string, data: ServerUpdateToolRequest) {
    try {
      const response = await toolsApi.toolsUpdate(id, data);
      return response.data;
    } catch (error) {
      console.error(`Failed to update tool ${id}:`, error);
      throw new Error('Failed to update tool');
    }
  }

  static async deleteTool(id: string) {
    try {
      await toolsApi.toolsDelete(id);
    } catch (error) {
      console.error(`Failed to delete tool ${id}:`, error);
      throw new Error('Failed to delete tool');
    }
  }

  static async checkoutTool(id: string, data: ServerCheckoutToolRequest) {
    try {
      const response = await toolsApi.toolsCheckoutCreate(id, data);
      return response.data;
    } catch (error) {
      console.error(`Failed to checkout tool ${id}:`, error);
      throw new Error('Failed to checkout tool');
    }
  }

  static async checkinTool(id: string, data: ServerCheckinToolRequest) {
    try {
      const response = await toolsApi.toolsCheckinCreate(id, data);
      return response.data;
    } catch (error) {
      console.error(`Failed to checkin tool ${id}:`, error);
      throw new Error('Failed to checkin tool');
    }
  }

  // Users API
  static async getUsers(params?: { limit?: number; offset?: number; role?: string }) {
    try {
      const response = await usersApi.usersList(params);
      return response.data;
    } catch (error) {
      console.error('Failed to fetch users:', error);
      throw new Error('Failed to fetch users');
    }
  }

  static async getUser(id: string) {
    try {
      const response = await usersApi.usersDetail(id);
      return response.data;
    } catch (error) {
      console.error(`Failed to fetch user ${id}:`, error);
      throw new Error(`Failed to fetch user`);
    }
  }

  static async createUser(data: ServerCreateUserRequest) {
    try {
      const response = await usersApi.usersCreate(data);
      return response.data;
    } catch (error) {
      console.error('Failed to create user:', error);
      throw new Error('Failed to create user');
    }
  }

  // Events API
  static async getEvents(params?: { 
    limit?: number; 
    offset?: number; 
    type?: string;
    tool_id?: string;
    user_id?: string;
  }) {
    try {
      const response = await eventsApi.eventsList(params);
      return response.data;
    } catch (error) {
      console.error('Failed to fetch events:', error);
      throw new Error('Failed to fetch events');
    }
  }

  // Admin API
  static async getStats() {
    try {
      const response = await adminApi.adminStatsList();
      return response.data;
    } catch (error) {
      console.error('Failed to fetch stats:', error);
      throw new Error('Failed to fetch stats');
    }
  }
}