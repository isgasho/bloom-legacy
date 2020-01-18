import axios from 'axios';

const { log } = require('@bloom42/astro');

const CALL_URL = 'http://localhost:8042/electronCall';

const empty = {};

async function call(method: string, params: any): Promise<any> {
  const message = JSON.stringify({
    method,
    params,
  });
  log.with({ msg: message }).debug('jsonMessage');

  const res = await axios.post(CALL_URL, message);
  log.with({ res: res.data }).debug('resMessage');

  const { data } = res;
  if (data.error !== null) {
    throw data.error;
  }

  return data.data;
}

function toIsoDate(date: string | null): Date | null {
  if (date === null) {
    return null;
  }
  return new Date(date).toISOString() as unknown as Date;
}

async function init(): Promise<void> {
  return call('core.init', empty);
}

export default {
  call,
  toIsoDate,
  empty,
  init,
};