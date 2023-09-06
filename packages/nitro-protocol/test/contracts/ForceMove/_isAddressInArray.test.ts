import {Contract} from 'ethers';

import ForceMoveArtifact from '../../../artifacts/contracts/test/TESTForceMove.sol/TESTForceMove.json';
import {generateParticipants, getTestProvider, setupContract} from '../../test-helpers';

const provider = getTestProvider();
let ForceMove: Contract;

const nParticipants = 4;
const {participants} = generateParticipants(nParticipants);

beforeAll(async () => {
  ForceMove = setupContract(provider, ForceMoveArtifact, process.env.TEST_FORCE_MOVE_ADDRESS);
});

describe('_isAddressInArray', () => {
  const suspect = participants[0];
  const addresses = participants.slice(1);

  it('verifies absence of suspect', async () => {
    expect(await ForceMove.isAddressInArray(suspect, addresses)).toBe(false);
  });
  it('finds an address hiding in an array', async () => {
    addresses[1] = suspect;
    expect(await ForceMove.isAddressInArray(suspect, addresses)).toBe(true);
  });
});
