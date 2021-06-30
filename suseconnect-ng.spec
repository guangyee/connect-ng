#
# spec file for package suseconnect-ng
#
# Copyright (c) 2021 SUSE LLC
#
# All modifications and additions to the file contributed by third parties
# remain the property of their copyright owners, unless otherwise agreed
# upon. The license for this file, and modifications and additions to the
# file, is the same license as for the pristine package itself (unless the
# license for the pristine package is not an Open Source License, in which
# case the license is the MIT License). An "Open Source License" is a
# license that conforms to the Open Source Definition (Version 1.9)
# published by the Open Source Initiative.

# Please submit bugfixes or comments via https://bugs.opensuse.org/
#

%global provider_prefix github.com/SUSE/connect-ng
%global import_path     %{provider_prefix}

Name:           suseconnect-ng
Version:        0.0.0~git.d73f68f
Release:        0
URL:            https://github.com/SUSE/connect-ng
License:        LGPL-2.1-or-later
Summary:        Utility to register a system with the SUSE Customer Center
Group:          System/Management
Source:         connect-ng-%{version}.tar.xz
Source1:        %name-rpmlintrc
BuildRequires:  golang-packaging

%description
This package provides a command line tool for connecting a
client system to the SUSE Customer Center. It will connect the system to your
product subscriptions and enable the product repositories/services locally.


%{go_nostrip}
%{go_provides}

%prep
%setup -q -n connect-ng-%{version}

%build
find %_builddir/..
%goprep %{import_path}
find %_builddir/..
go list -m all
%gobuild cmd
go build -buildmode=c-shared -o %_builddir/go/src/github.com/SUSE/connect-ng/ext/libsuseconnect.so ext/main.go
find %_builddir/..

%install
%goinstall
ln -s cmd %buildroot/%_bindir/SUSEConnect
mkdir %buildroot/%_sbindir
ln -s ../bin/cmd %buildroot/%_sbindir/SUSEConnect
#TODO package ruby module
#cp /home/abuild/rpmbuild/BUILD/go/src/github.com/SUSE/connect-ng/ext/libsuseconnect.so %_libdir/libsuseconnect.so
#TODO man pages not yet available in source, these are the names frome the ruby version
#/usr/share/man/man5/SUSEConnect.5.gz
#/usr/share/man/man8/SUSEConnect.8.gz

%gofilelist
find %_builddir/..

%check
%gotest github.com/SUSE/connect-ng/connect

%files -f file.lst
%license LICENSE LICENSE.LGPL
%doc README.md
%_bindir/cmd
%_bindir/SUSEConnect
%_sbindir/SUSEConnect
