import { useNavigate } from "react-router-dom";
import Header from "@/components/Header";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { Heart, Coffee, ExternalLink, Zap, Shield, Users } from "lucide-react";
import { useLanguage } from "@/contexts/LanguageContext";
import "./Donate.css";

const Donate = () => {
  const navigate = useNavigate();
  const { t } = useLanguage();

  const handleSupport = () => {
    // Redirect to Buy Me a Coffee
    window.open("https://buymeacoffee.com/muhammad_deve", "_blank");
  };

  return (
    <div className="min-h-screen bg-background">
      <Header onSearch={(q) => navigate(`/browse?search=${encodeURIComponent(q)}`)} />

      <div className="container mx-auto py-20 px-4">
        <div className="max-w-4xl mx-auto">
          {/* Header */}
          <div className="text-center mb-16">
            <div className="inline-flex items-center justify-center w-16 h-16 rounded-full bg-primary/20 mb-6">
              <Heart className="h-8 w-8 text-primary" />
            </div>
            <h1 className="text-5xl font-bold mb-4 text-foreground">Support Streamly</h1>
            <p className="text-xl text-muted-foreground max-w-2xl mx-auto">
              Help us keep providing free access to thousands of TV channels worldwide
            </p>
          </div>

          {/* Why Support Cards */}
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-16">
            {/* Card 1 */}
            <Card className="p-6 bg-card border-border hover:border-primary/50 transition-all hover:shadow-lg">
              <Zap className="h-8 w-8 text-primary mb-3" />
              <h3 className="text-lg font-bold mb-2 text-foreground">Keep It Running</h3>
              <p className="text-sm text-muted-foreground">
                Server costs, maintenance, and updates require ongoing support
              </p>
            </Card>

            {/* Card 2 */}
            <Card className="p-6 bg-card border-border hover:border-primary/50 transition-all hover:shadow-lg">
              <Shield className="h-8 w-8 text-primary mb-3" />
              <h3 className="text-lg font-bold mb-2 text-foreground">Stay Independent</h3>
              <p className="text-sm text-muted-foreground">
                Your support helps us remain completely free and independent
              </p>
            </Card>

            {/* Card 3 */}
            <Card className="p-6 bg-card border-border hover:border-primary/50 transition-all hover:shadow-lg">
              <Users className="h-8 w-8 text-primary mb-3" />
              <h3 className="text-lg font-bold mb-2 text-foreground">Community Driven</h3>
              <p className="text-sm text-muted-foreground">
                Every contribution helps improve the service for everyone
              </p>
            </Card>
          </div>

          {/* Main Support Section */}
          <Card className="p-12 bg-gradient-to-br from-primary/10 to-primary/5 border-2 border-primary/30 group hover:border-primary transition-all duration-300 hover:shadow-xl hover:shadow-primary/20">
            <div className="text-center mb-8">
              <Coffee className="h-12 w-12 text-primary mx-auto mb-4" />
              <h2 className="text-3xl font-bold mb-3 text-foreground">Buy Me a Coffee</h2>
              <p className="text-lg text-muted-foreground mb-8">
                Support us on Buy Me a Coffee and help keep Streamly running
              </p>
              
              <Button
                onClick={handleSupport}
                size="lg"
                className="bg-primary hover:bg-primary/90 text-primary-foreground font-bold text-lg px-8 py-6 rounded-lg transition-all hover:shadow-xl"
              >
                <Coffee className="h-6 w-6 mr-2" />
                Support Us Now
                <ExternalLink className="h-6 w-6 ml-2" />
              </Button>

              <p className="text-sm text-muted-foreground mt-6">
                You'll be redirected to Buy Me a Coffee where you can choose your support amount
              </p>
            </div>
          </Card>

          {/* Benefits Section */}
          <div className="mt-16 grid grid-cols-1 md:grid-cols-2 gap-8">
            {/* Left Card */}
            <Card className="p-8 bg-card border-border">
              <h3 className="text-2xl font-bold mb-6 text-foreground">What Your Support Does</h3>
              <ul className="space-y-4">
                <li className="flex items-start gap-3">
                  <span className="text-primary font-bold text-lg">✓</span>
                  <div>
                    <p className="font-semibold text-foreground">Keeps the Platform Running</p>
                    <p className="text-sm text-muted-foreground">24/7 server uptime and maintenance</p>
                  </div>
                </li>
                <li className="flex items-start gap-3">
                  <span className="text-primary font-bold text-lg">✓</span>
                  <div>
                    <p className="font-semibold text-foreground">New Features & Channels</p>
                    <p className="text-sm text-muted-foreground">Continuous improvements and expansions</p>
                  </div>
                </li>
                <li className="flex items-start gap-3">
                  <span className="text-primary font-bold text-lg">✓</span>
                  <div>
                    <p className="font-semibold text-foreground">Better Quality</p>
                    <p className="text-sm text-muted-foreground">Improved streaming quality and reliability</p>
                  </div>
                </li>
              </ul>
            </Card>

            {/* Right Card */}
            <Card className="p-8 bg-gradient-to-br from-primary/20 to-primary/10 border-primary/30">
              <h3 className="text-2xl font-bold mb-4 text-foreground">Thank You!</h3>
              <p className="text-muted-foreground mb-6">
                Every contribution, no matter the size, makes a huge difference. Your support directly helps us maintain and improve Streamly for everyone.
              </p>
              <div className="bg-background/50 rounded-lg p-4 border border-primary/20">
                <p className="text-sm text-muted-foreground">
                  <span className="font-semibold text-foreground">Questions?</span> Feel free to reach out through Buy Me a Coffee. We'd love to hear from you!
                </p>
              </div>
            </Card>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Donate;
