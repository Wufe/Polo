import { IAPISession } from '@/api/session';
import React, {  } from 'react';
import './index.scss';
import { HelperSession } from './session/helper-session';
import { render } from 'react-dom';
import { HelperStatusContext } from './contexts';
import { HelperOverlay } from './overlay/helper-overlay';
import { HelperStatusProvider } from './status/helper-status-provider';

type TProps = {
    session: IAPISession;
}
export const HelperApp = (props: TProps) => {

    if (!props.session)
        return null;

    const { age, status, killReason, replacedBy } = props.session;

    return <HelperStatusProvider
        uuid={props.session.uuid}
        initial={{
            age,
            status,
            killReason,
            replacedBy
        }}>
        <HelperStatusContext.Consumer>
            {({ helperStatus }) => <HelperOverlay helperStatus={helperStatus} />}
        </HelperStatusContext.Consumer>
        <div className="session-helper__component">
            <HelperSession session={props.session} />
        </div>
    </HelperStatusProvider>
};

render(<HelperApp session={window.currentSession} />, document.getElementById('polo-session-helper'));

declare global {
    interface Window {
        currentSession: any;
    }
}