import { useState, useEffect, useRef } from "react";
import { Link, useNavigate } from "react-router-dom";
import { Tv, Search, Menu, X, Sun, Moon, Globe } from "lucide-react";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { useTheme } from "@/contexts/ThemeContext";
import { useLanguage, Language } from "@/contexts/LanguageContext";
import { searchChannels, Channel } from "@/lib/channels";

interface HeaderProps {
  onSearch: (query: string) => void;
}

const Header = ({ onSearch }: HeaderProps) => {
  const [searchQuery, setSearchQuery] = useState("");
  const [searchResults, setSearchResults] = useState<Channel[]>([]);
  const [showResults, setShowResults] = useState(false);
  const [isSearching, setIsSearching] = useState(false);
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);
  const { theme, toggleTheme } = useTheme();
  const { language, setLanguage, t } = useLanguage();
  const navigate = useNavigate();
  const searchRef = useRef<HTMLDivElement>(null);

  // Close search results when clicking outside
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (searchRef.current && !searchRef.current.contains(event.target as Node)) {
        setShowResults(false);
      }
    };

    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, []);

  // Debounced search - real-time as user types
  useEffect(() => {
    const timer = setTimeout(async () => {
      if (searchQuery.trim().length >= 1) {
        setIsSearching(true);
        const results = await searchChannels(searchQuery);
        setSearchResults(results.channels);
        setShowResults(true);
        setIsSearching(false);
      } else {
        setSearchResults([]);
        setShowResults(false);
      }
    }, 300);

    return () => clearTimeout(timer);
  }, [searchQuery]);

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    if (searchResults.length > 0) {
      // Navigate to first result
      handleChannelClick(searchResults[0]);
    }
  };

  const handleChannelClick = (channel: Channel) => {
    setShowResults(false);
    setSearchQuery("");
    navigate(`/channel/${encodeURIComponent(channel.id)}`, {
      state: { channelData: channel }
    });
  };

  const languageNames: Record<Language, string> = {
    en: "English",
    ru: "Русский",
    uz: "O'zbek",
  };

  return (
    <header className="sticky top-0 z-50 w-full border-b border-border bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
      <div className="container flex h-16 items-center justify-between px-4">
        <div className="flex items-center gap-8">
          <Link to="/" className="flex items-center gap-2 text-foreground hover:text-primary transition-colors">
            <Tv className="h-6 w-6 text-primary" />
            <span className="font-bold text-xl">Streamly</span>
          </Link>

          <nav className="hidden md:flex items-center gap-6">
            <Link to="/" className="text-sm font-medium text-muted-foreground hover:text-foreground transition-colors">
              {t("header.home")}
            </Link>
            <Link to="/browse" className="text-sm font-medium text-muted-foreground hover:text-foreground transition-colors">
              {t("header.browse")}
            </Link>
            <Link to="/donate" className="text-sm font-medium text-muted-foreground hover:text-foreground transition-colors">
              {t("header.donate")}
            </Link>
          </nav>
        </div>

        <div className="flex items-center gap-2">
          <form onSubmit={handleSearch} className="hidden sm:flex items-center gap-2">
            <div className="relative" ref={searchRef}>
              <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
              <Input
                type="search"
                placeholder={t("header.search")}
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="pl-9 w-64 bg-secondary border-border"
              />
              
              {/* Search Results Dropdown */}
              {showResults && searchQuery.trim().length >= 1 && (
                <div className="absolute top-full mt-2 w-96 max-h-96 overflow-y-auto bg-popover border border-border rounded-lg shadow-lg z-50">
                  {isSearching ? (
                    <div className="p-4 text-center text-muted-foreground">
                      Searching...
                    </div>
                  ) : searchResults.length > 0 ? (
                    <div className="py-2">
                      {searchResults.map((channel) => (
                        <button
                          key={channel.id}
                          onClick={() => handleChannelClick(channel)}
                          className="w-full px-4 py-3 flex items-center gap-3 hover:bg-accent transition-colors text-left"
                        >
                          <img
                            src={channel.logo}
                            alt={channel.name}
                            className="w-12 h-12 rounded object-cover"
                            onError={(e) => {
                              e.currentTarget.src = 'https://images.unsplash.com/photo-1504711434969-e33886168f5c?w=200&h=200&fit=crop';
                            }}
                          />
                          <div className="flex-1 min-w-0">
                            <p className="font-medium text-foreground truncate">{channel.name}</p>
                            <p className="text-xs text-muted-foreground truncate">
                              {channel.category} • {channel.country}
                            </p>
                          </div>
                        </button>
                      ))}
                    </div>
                  ) : (
                    <div className="p-4 text-center text-muted-foreground">
                      No channels found
                    </div>
                  )}
                </div>
              )}
            </div>
          </form>

          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" size="icon">
                <Globe className="h-5 w-5" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end" className="bg-popover z-50">
              {(Object.keys(languageNames) as Language[]).map((lang) => (
                <DropdownMenuItem
                  key={lang}
                  onClick={() => setLanguage(lang)}
                  className={language === lang ? "bg-accent" : ""}
                >
                  {languageNames[lang]}
                </DropdownMenuItem>
              ))}
            </DropdownMenuContent>
          </DropdownMenu>

          <Button variant="ghost" size="icon" onClick={toggleTheme}>
            {theme === "light" ? <Moon className="h-5 w-5" /> : <Sun className="h-5 w-5" />}
          </Button>

          <Button
            variant="ghost"
            size="icon"
            className="md:hidden"
            onClick={() => setIsMobileMenuOpen(!isMobileMenuOpen)}
          >
            {isMobileMenuOpen ? <X className="h-5 w-5" /> : <Menu className="h-5 w-5" />}
          </Button>
        </div>
      </div>

      {/* Mobile menu */}
      {isMobileMenuOpen && (
        <div className="md:hidden border-t border-border bg-background">
          <nav className="container flex flex-col gap-4 p-4">
            <Link
              to="/"
              className="text-sm font-medium text-muted-foreground hover:text-foreground transition-colors"
              onClick={() => setIsMobileMenuOpen(false)}
            >
              {t("header.home")}
            </Link>
            <Link
              to="/browse"
              className="text-sm font-medium text-muted-foreground hover:text-foreground transition-colors"
              onClick={() => setIsMobileMenuOpen(false)}
            >
              {t("header.browse")}
            </Link>
            <Link
              to="/donate"
              className="text-sm font-medium text-muted-foreground hover:text-foreground transition-colors"
              onClick={() => setIsMobileMenuOpen(false)}
            >
              {t("header.donate")}
            </Link>
            
            <form onSubmit={handleSearch} className="sm:hidden">
              <div className="relative">
                <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                <Input
                  type="search"
                  placeholder={t("header.search")}
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  className="pl-9 bg-secondary border-border"
                />
              </div>
            </form>
          </nav>
        </div>
      )}
    </header>
  );
};

export default Header;
