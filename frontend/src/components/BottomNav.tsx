"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { List, Map as MapIcon, PlusCircle, User } from "lucide-react";

export function BottomNav() {
    const pathname = usePathname();

    // Hide nav on login/signup to keep focus simple
    if (pathname === "/login" || pathname === "/signup") {
        return null;
    }

    const navItems = [
        { href: "/feed", icon: List, label: "Feed" },
        { href: "/map", icon: MapIcon, label: "Map" },
        { href: "/create", icon: PlusCircle, label: "Create" },
        { href: "/profile", icon: User, label: "Profile" },
    ];

    return (
        <nav className="sm:hidden fixed bottom-0 left-0 right-0 z-50 bg-background/80 backdrop-blur-lg border-t pb-safe">
            <ul className="flex justify-around items-center h-16 px-2">
                {navItems.map((item) => {
                    const isActive = pathname === item.href || (item.href !== "/" && pathname.startsWith(item.href));
                    return (
                        <li key={item.href} className="w-full">
                            <Link
                                href={item.href}
                                className="flex flex-col items-center justify-center w-full h-full space-y-1 group"
                                aria-label={item.label}
                            >
                                <div
                                    className={`p-1.5 rounded-full transition-colors ${isActive
                                            ? "bg-primary/10 text-primary"
                                            : "text-muted-foreground group-hover:bg-muted group-hover:text-foreground"
                                        }`}
                                >
                                    <item.icon className={`h-5 w-5 ${isActive ? "stroke-[2.5px]" : "stroke-2"}`} />
                                </div>
                                <span className={`text-[10px] font-medium transition-colors ${isActive ? "text-primary" : "text-muted-foreground"}`}>
                                    {item.label}
                                </span>
                            </Link>
                        </li>
                    );
                })}
            </ul>
        </nav>
    );
}
