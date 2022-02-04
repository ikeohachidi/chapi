export interface PermOrigin {
    id: number;
    url: string;
    routeId: number; 
}

export interface RoutePermOrigin {[routeId: number]: PermOrigin[]}

export class MergeOptions {
    routeId = 0;
    mergeHeader = false;
    mergeQuery = false;
    mergeBody = false;

    constructor(init?: Partial<MergeOptions>) {
        if (init) {
            init.mergeQuery ? this.mergeQuery = init.mergeQuery : null;
            init.mergeHeader ? this.mergeHeader = init.mergeHeader : null;
            init.mergeBody ? this.mergeBody = init.mergeBody : null;
        }
        this.routeId = 0;
    }
}

export interface RouteMergeOption {[routeId: number]: MergeOptions}