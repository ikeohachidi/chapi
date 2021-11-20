export interface PermOrigin {
    id: number;
    url: string;
    routeId: number; 
}

export interface RoutePermOrigin {[routeId: number]: PermOrigin[]};