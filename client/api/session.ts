import { ISession, ISessionLog } from '@/state/models/session-model';
import { SessionStatus, SessionKillReason } from '@/state/models/session-model-enums';
import Axios from 'axios';
import { buildRequest } from './common';

export interface IAPISession extends Omit<ISession, 'logs'> {
    logs: ISessionLog[];
}

export interface IAPISessionLogsAndStatus {
    logs: ISessionLog[];
    status: SessionStatus;
}

export interface IApiSessionStatus {
    status: SessionStatus;
    age: number;
    killReason: SessionKillReason;
    replacedBy: string;
}

export function killSessionAPI(uuid: string) {
    return buildRequest<void>(() => Axios.delete(`/_polo_/api/session/${uuid}`));
}

export function retrieveSessionAPI(uuid: string) {
    return buildRequest<IAPISession>(() => Axios.get(`/_polo_/api/session/${uuid}`));
}

export function trackSessionAPI(uuid: string) {
    return buildRequest<void>(() => Axios.post(`/_polo_/api/session/${uuid}/track`));
}

export function untrackSessionAPI() {
    return buildRequest<void>(() => Axios.delete(`/_polo_/api/session/<none>/track`))
}

export function retrieveSessionStatusAPI(uuid: string) {
    return buildRequest<IApiSessionStatus>(() => Axios.get(`/_polo_/api/session/${uuid}/status`));
}

export function retrieveLogsAndStatusAPI(uuid: string, lastUUID: string = "<none>") {
    return buildRequest<IAPISessionLogsAndStatus>(() => Axios.get(`/_polo_/api/session/${uuid}/logs/${lastUUID}`));
}