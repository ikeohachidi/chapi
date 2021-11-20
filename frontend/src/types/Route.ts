import { HTTPMethod } from "./HTTP";

export default class Route {
    id?: number;
    name: string;
    description: string;
    destination: string;
    url: string;
    method: HTTPMethod;
    path: string;
    queries: Query[];
    body: string;
    authorizedURLs: AuthorizedURL[];
    projectId: number;

    constructor() {
        this.id = 0;
        this.name = '';
        this.description = '';
        this.destination = '';
        this.url = '';
        this.method = HTTPMethod.GET;
        this.path = '';
        this.queries = [];
        this.body = '';
        this.authorizedURLs = [];
        this.projectId = 0;
    }
}

export type CreateRouteRequest = {
    projectId: number;
    path: string;
    method: string;
    description: string;
} 

export type AuthorizedURL = {
    id: string;
    url: string;
    routeId: string;
}

export type Query = {
    id?: number;
    routeId?: number;
    name: string;
    value: string;
}

export type ProjectRoute = {
    projectId: string;
    routes: Route[]
}

export type ProjectRouteQuery = {
    projectId: number;
    routeId: number;
    query: Query;
}