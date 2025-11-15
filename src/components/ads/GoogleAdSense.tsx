import { useEffect } from 'react';

interface AdSenseObject {
  [key: string]: unknown;
}

declare global {
  interface Window {
    adsbygoogle: (AdSenseObject | undefined)[];
  }
}

interface GoogleAdSenseProps {
  adSlot: string;
  adFormat?: string;
  className?: string;
  style?: React.CSSProperties;
}

const GoogleAdSense = ({ 
  adSlot, 
  adFormat = 'auto', 
  className = '', 
  style = { display: 'block' } 
}: GoogleAdSenseProps) => {
  useEffect(() => {
    try {
      if (typeof window !== 'undefined' && window.adsbygoogle) {
        window.adsbygoogle.push({});
      }
    } catch (error) {
      console.error('AdSense error:', error);
    }
  }, []);

  return (
    <ins
      className={`adsbygoogle ${className}`}
      style={style}
      data-ad-client="ca-pub-3638974529306119"
      data-ad-slot={adSlot}
      data-ad-format={adFormat}
      data-full-width-responsive="true"
    />
  );
};

export default GoogleAdSense;
