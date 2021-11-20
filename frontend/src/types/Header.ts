export default interface Header {
    routeId?: number;
    id?: number;
    name: string;
    value: string;
}

export const HeaderDefault = {
    routeId: 0,
    id: 0,
    name: '',
    value: ''
}