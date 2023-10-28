/**********************************************************************
 *
 *   Component hook generated by Quest
 *
 *   Code Logic for the component goes in this hook
 *   To setup bindings that use the data here or call the functions here, use the Quest editor
 *   Do not change the name of the hook
 *
 *   For help and further details refer to: https://www.quest.ai/docs
 *
 *
 **********************************************************************/

import { useMemo } from 'react';
import { CardRewardProps } from 'src/types';
import { useNavigate } from 'react-router-dom';
const useCardReward = (props: CardRewardProps) => {
  const navigate = useNavigate();
  
  const { reward } = props;

  const imageUrl = useMemo(() => {
    switch (reward.name) {
      case 'Target':
        return '/images/target.png';
      case 'Starbucks Coffee':
        return '/images/starbucks.png';
      case 'Jumba Juice':
        return '/images/jamba.png';
      case 'Grubhub':
        return '/images/grubhub.png';
    }

  }, [reward.name]);
  const rewardDetailsPage = (rewardID: string) => {
    navigate(`/reward-details/${rewardID}`);
  };

  const fns: any = { rewardDetailsPage };
  const data: any = { imageUrl };

  return { data, fns };
};

export default useCardReward;
